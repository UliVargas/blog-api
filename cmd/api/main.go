package main

import (
	"log"

	"github.com/UliVargas/blog-go/internal/config"
	"github.com/UliVargas/blog-go/internal/handlers"
	"github.com/UliVargas/blog-go/internal/models"
	"github.com/UliVargas/blog-go/internal/repository"
	"github.com/UliVargas/blog-go/internal/service"
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
	db.AutoMigrate(&models.User{})

	// Inicialización de servicios
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	// Inicialización de router
	router := gin.Default()

	// Rutas de usuarios
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.GET("/", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.POST("/", userHandler.Create)
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
