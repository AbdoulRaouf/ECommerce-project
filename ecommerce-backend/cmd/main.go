package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thiomajid/ecommerce-backend/configs"
	"github.com/thiomajid/ecommerce-backend/internal/database"
	"github.com/thiomajid/ecommerce-backend/internal/handlers"
)

func main() {
	// Charger la configuration
	config := configs.LoadConfig()

	// Configurer le mode de Gin
	if config.ProductionMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialiser la base de données
	if err := database.InitDB(config.DatabasePath); err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données: %v", err)
	}

	// Créer un router Gin
	router := gin.Default()

	// Configurer CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.AllowedOrigins}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "User-ID", "Admin-Token"}
	router.Use(cors.New(corsConfig))

	// Routes pour les produits
	router.GET("/api/products", handlers.GetAllProducts)
	router.GET("/api/products/:id", handlers.GetProduct)
	router.GET("/api/categories", handlers.GetCategories)
	router.POST("/api/products", handlers.CreateProduct)
	router.PUT("/api/products/:id", handlers.UpdateProduct)
	router.DELETE("/api/products/:id", handlers.DeleteProduct)

	// Routes pour le panier
	router.GET("/api/cart", handlers.GetCartItems)
	router.POST("/api/cart", handlers.AddToCart)
	router.PUT("/api/cart/:id", handlers.UpdateCartItem)
	router.DELETE("/api/cart/:id", handlers.RemoveFromCart)
	router.DELETE("/api/cart", handlers.ClearCart)

	// Routes pour les commandes
	router.POST("/api/orders", handlers.CreateOrder)
	router.GET("/api/orders", handlers.GetOrders)
	router.GET("/api/orders/:id", handlers.GetOrderDetails)

	// Route d'authentification admin
	router.POST("/api/admin/login", handlers.LoginAdmin)

	// Routes admin (protégées par le middleware d'authentification)
	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(handlers.AdminAuthMiddleware())
	{
		// Tableau de bord admin
		adminRoutes.GET("/dashboard", handlers.AdminGetDashboardStats)

		// Gestion des produits
		adminRoutes.GET("/products", handlers.GetAllProducts)       // Réutilisation de la route existante
		adminRoutes.POST("/products", handlers.CreateProduct)       // Réutilisation de la route existante
		adminRoutes.PUT("/products/:id", handlers.UpdateProduct)    // Réutilisation de la route existante
		adminRoutes.DELETE("/products/:id", handlers.DeleteProduct) // Réutilisation de la route existante

		// Gestion des commandes
		adminRoutes.GET("/orders", handlers.AdminGetAllOrders)
		adminRoutes.GET("/orders/:id", handlers.AdminGetOrderDetails)
		adminRoutes.PUT("/orders/:id/status", handlers.AdminUpdateOrderStatus)
		adminRoutes.DELETE("/orders/:id", handlers.AdminDeleteOrder)
	}

	// Démarrer le serveur
	port := fmt.Sprintf(":%s", config.Port)
	log.Printf("Serveur démarré sur le port %s", config.Port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
