package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Palma99/feature-flag-service/config"
	"github.com/Palma99/feature-flag-service/internals/application/services"
	usecase "github.com/Palma99/feature-flag-service/internals/application/usecase"
	app_middleware "github.com/Palma99/feature-flag-service/internals/infrastructure/middleware"
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

	publicServiceController := interfaces.NewPublicServiceController(flagInteractor)

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

	r.Route("/project", func(r chi.Router) {
		r.Use(app_middleware.CheckAuthMiddleware(jwtService))

		r.Post("/", projectController.CreateProject)
		r.Get("/{id}", projectController.GetProject)
		r.Get("/list", projectController.GetProjects)
	})

	r.Route("/environment", func(r chi.Router) {
		r.Use(app_middleware.CheckAuthMiddleware(jwtService))

		r.Get("/{id}", environmentController.GetEnvironment)
		r.Post("/", environmentController.CreateEnvironment)
	})

	r.Route("/flag", func(r chi.Router) {
		r.Use(app_middleware.CheckAuthMiddleware(jwtService))

		r.Post("/", flagController.CreateFlag)
		r.Delete("/{id}", flagController.DeleteFlag)
		r.Put("/environment/{id}", flagController.UpdateFlagEnvironment)
		r.Get("/project/{id}", flagController.GetProjectFlags)
	})

	r.Route("/public/v1", func(r chi.Router) {
		r.Use(app_middleware.CheckPublicKeyAuthMiddleware(keyService))

		r.Get("/flags", publicServiceController.GetFlagsByPublicKey)
	})

	fmt.Println("Server started on http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
