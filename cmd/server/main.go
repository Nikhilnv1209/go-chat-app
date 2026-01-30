package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"chat-app/internal/config"
	"chat-app/internal/database"
	"chat-app/internal/handlers"
	"chat-app/internal/middleware"
	"chat-app/internal/models"
	"chat-app/internal/repository"
	"chat-app/internal/service"
	"chat-app/internal/websocket"
	"chat-app/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	cfg := config.Load()

	// 2. Initialize Database
	database.InitDB(cfg)
	db := database.GetDB()

	// 3. Auto Migrate Models
	log.Println("Running AutoMigrate...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.GroupMember{},
		&models.Message{},
		&models.MessageReceipt{},
		&models.Conversation{},
		&models.RefreshToken{},
	)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	log.Println("Database Migration Clean!")

	// 4. Initialization
	// Repositories
	userRepo := repository.NewUserRepository(db)
	msgRepo := repository.NewMessageRepository(db)
	convRepo := repository.NewConversationRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	receiptRepo := repository.NewMessageReceiptRepository(db)    // [F06]
	refreshTokenRepo := repository.NewRefreshTokenRepository(db) // [F09]

	// WebSocket Hub
	// We create this early because MessageService needs it
	hub := websocket.NewHub(userRepo, convRepo)
	go hub.Run()

	// Services
	jwtService := jwt.NewService(jwt.Config{
		Secret:     cfg.JWT.Secret,
		Expiration: cfg.JWT.Expiration,
	})
	authService := service.NewAuthService(userRepo, refreshTokenRepo, jwtService)
	msgService := service.NewMessageService(msgRepo, convRepo, groupRepo, receiptRepo, userRepo, hub) // [F06][F07]

	groupService := service.NewGroupService(groupRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	wsHandler := handlers.NewWSHandler(hub, authService)
	groupHandler := handlers.NewGroupHandler(groupService)
	chatHandler := handlers.NewChatHandler(convRepo, msgRepo, userRepo, groupRepo, msgService)

	// INJECT MessageService into Hub/Client factory if needed?
	// Actually, the new handlers.WSHandler logic just passes the hub.
	// But the Client needs the msgService.
	// Clients are created in wsHandler.ServeWS. We need to pass msgService to wsHandler.
	wsHandler.MsgService = msgService

	// 5. Server Setup
	// Using gin.New() for explicit middleware control as per specs/03_Technical_Specification.md
	r := gin.New()
	r.Use(gin.Recovery())                // Panic recovery
	r.Use(middleware.LoggerMiddleware()) // Custom request logging [F00]
	r.Use(middleware.CORSMiddleware())   // CORS headers [F00]

	// Routes
	// Public routes (no auth required)
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
		authRoutes.POST("/logout", authHandler.Logout)
	}

	// Protected routes (require JWT auth)
	// Chat Routes (Inbox & History)
	chatRoutes := r.Group("/")
	chatRoutes.Use(middleware.AuthMiddleware(jwtService)) // [F00] Auth Middleware
	{
		chatRoutes.GET("/conversations", chatHandler.GetConversations)
		chatRoutes.GET("/messages", chatHandler.GetMessages)
		chatRoutes.POST("/messages/:id/read", chatHandler.MarkRead)
		chatRoutes.GET("/messages/:id/receipts", chatHandler.GetReceipts)
		chatRoutes.GET("/users", authHandler.SearchUsers)
		chatRoutes.GET("/users/:id", authHandler.GetUser)
	}

	// Group Routes (protected)
	groupRoutes := r.Group("/groups")
	groupRoutes.Use(middleware.AuthMiddleware(jwtService)) // [F00] Auth Middleware
	{
		groupRoutes.POST("", groupHandler.CreateGroup)
		groupRoutes.POST("/:id/members", groupHandler.AddMember)
	}

	// WebSocket Route
	r.GET("/ws", wsHandler.ServeWS)

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 6. Start Server with Timeouts and Graceful Shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server started on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed: ", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	log.Println("Shutting down WebSocket Hub...")
	if err := hub.Shutdown(shutdownCtx); err != nil {
		log.Printf("Hub shutdown error: %v", err)
	}

	log.Println("Shutting down HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Server shutdown complete")
}
