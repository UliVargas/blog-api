package main

import (
	"log"

	"github.com/UliVargas/blog-go/internal/application/service"
	"github.com/UliVargas/blog-go/internal/domain/model"
	"github.com/UliVargas/blog-go/internal/infrastructure/config"
	"github.com/UliVargas/blog-go/internal/infrastructure/repository"
	"github.com/UliVargas/blog-go/internal/presentation/handler"
	"github.com/UliVargas/blog-go/internal/presentation/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carga de varables de entorno de .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env", err)
	}

	// Inicialización de la base de datos
	db := config.DBConnect()
	db.AutoMigrate(&model.User{})

	// Inicialización de servicios
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

	// Inicialización de router
	router := gin.Default()

	// Rutas de usuarios
	api := router.Group("/api/v1")
	{
		// Rutas protegidas de usuarios
		protectedUsers := api.Group("/users")
		protectedUsers.Use(middleware.AuthMiddleware())
		{
			protectedUsers.GET("/", userHandler.GetAll)
			protectedUsers.GET("/:id", userHandler.GetByID)
		}

		// Rutas de autenticación
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}
	}

	// Rutas
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// Carga de configuración
	cfg := config.Load()

	// Ejecución del servidor
	router.Run(cfg.PORT)
}
