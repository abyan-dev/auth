package server

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/abyan-dev/auth/pkg/handler"
	"github.com/abyan-dev/auth/pkg/middleware"
	"github.com/abyan-dev/auth/pkg/model"
	"github.com/abyan-dev/auth/pkg/utils"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	slog.Info("Loading environment variables...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	slog.Info("Connecting to database...")
	db, err := utils.InitDB(config)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	slog.Info("Applying database migrations...")
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Error auto-migrating database: %v", err)
	}

	slog.Info("Setting up the app...")

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	slog.Info("Loading routes...")

	initRouter(app)

	return app
}

func initRouter(app fiber.Router) {
	api := app.Group("/api")

	api.Get("/health", handler.Health)
	api.Get("/health/protected", middleware.RequireAuthenticated(), handler.HealthProtected)
	api.Post("/auth/register/request", handler.RequestRegistration)
	api.Post("/auth/register/verify", middleware.RequireAuthenticated(), handler.VerifyRegistration)
	api.Post("/auth/login", handler.Login)
	api.Post("/auth/logout", middleware.RequireAuthenticated(), handler.Logout)
}

func Run(app *fiber.App) {
	slog.Info("Server is now listening on port 8080...")
	if err := app.Listen(":8080"); err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}
