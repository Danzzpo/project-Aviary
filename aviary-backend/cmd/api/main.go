package main

import (
	"aviary-backend/config"
	"aviary-backend/controllers"
	"aviary-backend/middleware"
	"aviary-backend/models"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Connect DB & Migrate
	config.ConnectDatabase()

	// --- UPDATE MIGRASI DI SINI ---
	// Menambahkan model Production dan Egg
	// Kita gunakan Set table options InnoDB untuk keamanan relasi
	err = config.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.User{}, 
		&models.Bird{}, 
		&models.Pair{}, 
		&models.Production{}, // <--- BARU
		&models.Egg{},        // <--- BARU
	)
	
	if err != nil {
		log.Fatal("Gagal migrasi:", err)
	}

	// 3. Setup Router
	r := gin.Default()

	// --- SETUP CORS ---
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"}, // Izinkan 5173 & 5174 jaga-jaga
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Pong!", "status": "online"})
	})

	// Public Routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Protected Routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Birds
		protected.GET("/birds", controllers.GetBirds)
		protected.POST("/birds", controllers.CreateBird)
		protected.PUT("/birds/:id", controllers.UpdateBird)
		
		// Dashboard
		protected.GET("/dashboard/stats", controllers.GetDashboardStats)
		
		// Pairing (Penjodohan)
		protected.GET("/pairs", controllers.GetActivePairs)
		protected.POST("/pairs", controllers.CreatePair)
		protected.PUT("/pairs/:id/disband", controllers.DisbandPair)

		// --- PRODUCTION & EGGS (BARU) ---
		protected.GET("/pairs/:pair_id/production", controllers.GetActiveProduction) // Cek status telur
		protected.POST("/pairs/:pair_id/eggs", controllers.AddEgg)                   // Tambah telur
		protected.PUT("/eggs/:egg_id/status", controllers.UpdateEggStatus)           // Update (Menetas/Zonk)
		protected.DELETE("/eggs/:egg_id", controllers.DeleteEgg)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}