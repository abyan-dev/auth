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
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	DB *gorm.DB
}

func (s *Server) New() *fiber.App {
	slog.Info("Loading environment variables...")

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	slog.Info("Connecting to database...")
	db, err := utils.InitDB(config)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	s.DB = db

	slog.Info("Applying database migrations...")
	if err := db.AutoMigrate(&model.User{}, &model.RevokedToken{}); err != nil {
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

	s.initRouter(app)

	return app
}

func (s *Server) initRouter(app fiber.Router) {
	api := app.Group("/api")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",                  // Allow specific origin
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS", // Allow all methods
		AllowHeaders:     "Content-Type, Authorization",            // Allow specific headers
		AllowCredentials: true,
	}))

	api.Get("/health", handler.Health)
	api.Get("/health/protected", middleware.RequireAuthenticated(), handler.HealthProtected)

	api.Post("/auth/register", handler.Register)
	api.Post("/auth/verify", handler.Verify)
	api.Post("/auth/login", handler.Login)
	api.Post("/auth/logout", middleware.RequireAuthenticated(), handler.Logout)
	api.Get("/auth/decode", middleware.RequireAuthenticated(), handler.Decode)
	api.Post("/auth/2fa/email", handler.OTPEmail)
}

func (s *Server) Run(app *fiber.App) {
	slog.Info("Running cleanup scheduler on a separate goroutine...")
	cleanupScheduler := CreateCleanupScheduler(s.DB)
	go cleanupScheduler.Start()

	slog.Info("Server is now listening on port 8080...")
	if err := app.Listen(":8080"); err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}
