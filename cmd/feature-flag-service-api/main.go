package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/Palma99/feature-flag-service/config"
	"github.com/Palma99/feature-flag-service/internals/application/services"
	usecase "github.com/Palma99/feature-flag-service/internals/application/usecase"
	context_keys "github.com/Palma99/feature-flag-service/internals/infrastructure/context"
	infrastructure "github.com/Palma99/feature-flag-service/internals/infrastructure/repository"
	interfaces "github.com/Palma99/feature-flag-service/internals/interfaces/api"
	_ "github.com/lib/pq"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var db *sql.DB

var applicationConfig *config.Config

func init() {

	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	applicationConfig = config

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		applicationConfig.DB.User,
		applicationConfig.DB.Password,
		applicationConfig.DB.Host,
		applicationConfig.DB.Port,
		applicationConfig.DB.Database,
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	db = dbConn
}
func main() {
	keyService := services.NewKeyService()
	passwordService := services.NewPasswordService()
	jwtService := services.NewJWTService(applicationConfig.JwtSecret)

	userRepository := infrastructure.NewPgUserRepository(db)
	authInteractor := usecase.NewAuthInteractor(userRepository, passwordService, jwtService)
	authController := interfaces.NewApiUserController(authInteractor)

	projectRepository := infrastructure.NewPgProjectRepository(db)
	projectInteractor := usecase.NewProjectInteractor(projectRepository, keyService)
	projectController := interfaces.NewApiProjectController(projectInteractor)

	environmentRepository := infrastructure.NewPgEnvironmentRepository(db)
	environmentInteractor := usecase.NewEnvironmentInteractor(environmentRepository, projectRepository, keyService)
	environmentController := interfaces.NewApiEnvironmentController(environmentInteractor)

	flagRepository := infrastructure.NewPgFlagRepository(db)
	flagInteractor := usecase.NewFlagInteractor(flagRepository, environmentRepository, keyService, projectRepository)
	flagController := interfaces.NewApiFlagController(flagInteractor)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authController.GetToken)
	})

	// func checkAuth(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		// create new context from `r` request context, and assign key `"user"`
	// 		// to value of `"123"`
	// 		ctx := context.WithValue(r.Context(), "user", "123")

	// 		// call the next handler in the chain, passing the response writer and
	// 		// the updated request object with the new context value.
	// 		//
	// 		// note: context.Context values are nested, so any previously set
	// 		// values will be accessible as well, and the new `"user"` key
	// 		// will be accessible from this point forward.
	// 		next.ServeHTTP(w, r.WithContext(ctx))
	// 	})
	// }

	r.Route("/project", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// get token from auth header
				bearerToken := r.Header.Get("Authorization")
				parts := strings.Split(bearerToken, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				token := parts[1]

				// validate token
				if payload, err := jwtService.ValidateToken(token); err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else {

					ctx := context.WithValue(r.Context(), context_keys.UserIDKey, payload.UserID)
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			})
		})

		r.Post("/", projectController.CreateProject)
		r.Get("/{id}", projectController.GetProject)
		r.Get("/list", projectController.GetProjects)
	})

	r.Route("/environment", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// get token from auth header
				bearerToken := r.Header.Get("Authorization")
				parts := strings.Split(bearerToken, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				token := parts[1]

				// validate token
				if payload, err := jwtService.ValidateToken(token); err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else {

					ctx := context.WithValue(r.Context(), context_keys.UserIDKey, payload.UserID)
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			})
		})

		r.Get("/{id}", environmentController.GetEnvironment)
	})

	r.Route("/flag", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// get token from auth header
				bearerToken := r.Header.Get("Authorization")
				parts := strings.Split(bearerToken, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				token := parts[1]

				// validate token
				if payload, err := jwtService.ValidateToken(token); err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else {

					ctx := context.WithValue(r.Context(), context_keys.UserIDKey, payload.UserID)
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			})
		})

		r.Post("/", flagController.CreateFlag)
	})

	http.ListenAndServe(":3000", r)
}
