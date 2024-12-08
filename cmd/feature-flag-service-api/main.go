package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Palma99/feature-flag-service/config"
	"github.com/Palma99/feature-flag-service/internals/application/services"
	usecase "github.com/Palma99/feature-flag-service/internals/application/usecase"
	infrastructure "github.com/Palma99/feature-flag-service/internals/infrastructure/repository"
	"github.com/Palma99/feature-flag-service/internals/interfaces"
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
	// TODO: implement
	// keyService := services.NewKeyService()
	passwordService := services.NewPasswordService()
	jwtService := services.NewJWTService(applicationConfig.JwtSecret)

	userRepository := infrastructure.NewPgUserRepository(db)
	authInteractor := usecase.NewAuthInteractor(userRepository, passwordService, jwtService)
	authController := interfaces.NewUserController(authInteractor)

	// projectRepository := infrastructure.NewPgProjectRepository(db)
	// projectInteractor := usecase.NewProjectInteractor(projectRepository, keyService)
	// projectController := interfaces.NewProjectController(projectInteractor)

	// environmentRepository := infrastructure.NewPgEnvironmentRepository(db)
	// environmentInteractor := usecase.NewEnvironmentInteractor(environmentRepository, projectRepository, keyService)
	// environmentController := interfaces.NewEnvironmentInteractor(environmentInteractor)

	// flagRepository := infrastructure.NewPgFlagRepository(db)
	// flagInteractor := usecase.NewFlagInteractor(flagRepository, environmentRepository, keyService)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			var loginDTO usecase.UserLoginDTO
			err := json.NewDecoder(r.Body).Decode(&loginDTO)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token, err := authController.GetToken(loginDTO)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusOK)

			jsonResponse := map[string]string{
				"token": *token,
			}

			json.NewEncoder(w).Encode(jsonResponse)
		})
	})

	http.ListenAndServe(":3000", r)
}
