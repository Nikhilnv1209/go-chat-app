package main

import (
	"log"

	"chat-app/internal/config"
	"chat-app/internal/database"
	"chat-app/internal/handlers"
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
	database.InitDB()
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
	receiptRepo := repository.NewMessageReceiptRepository(db) // [F06]

	// WebSocket Hub
	// We create this early because MessageService needs it
	hub := websocket.NewHub(userRepo)
	go hub.Run()

	// Services
	jwtService := jwt.NewService(jwt.Config{
		Secret:     cfg.JWT.Secret,
		Expiration: cfg.JWT.Expiration,
	})
	authService := service.NewAuthService(userRepo, jwtService)
	msgService := service.NewMessageService(msgRepo, convRepo, groupRepo, receiptRepo, userRepo, hub) // [F06][F07]

	groupService := service.NewGroupService(groupRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	wsHandler := handlers.NewWSHandler(hub, authService)
	groupHandler := handlers.NewGroupHandler(groupService, authService)
	chatHandler := handlers.NewChatHandler(convRepo, msgRepo, userRepo, groupRepo, authService, msgService)

	// INJECT MessageService into Hub/Client factory if needed?
	// Actually, the new handlers.WSHandler logic just passes the hub.
	// But the Client needs the msgService.
	// Clients are created in wsHandler.ServeWS. We need to pass msgService to wsHandler.
	wsHandler.MsgService = msgService

	// 5. Server Setup
	r := gin.Default()

	// Routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Chat Routes (Inbox & History)
	// Grouping chat-related routes for better organization
	chatRoutes := r.Group("/") // Use root path for chat routes to maintain existing URLs
	{
		chatRoutes.GET("/conversations", chatHandler.GetConversations)
		chatRoutes.GET("/messages", chatHandler.GetMessages)
		chatRoutes.POST("/messages/:id/read", chatHandler.MarkRead)
		chatRoutes.GET("/messages/:id/receipts", chatHandler.GetReceipts)
	}

	// Group Routes
	groupRoutes := r.Group("/groups")
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

	log.Printf("Server started on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
