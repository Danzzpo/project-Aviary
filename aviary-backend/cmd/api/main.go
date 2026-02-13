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

	// Pastikan semua tabel terdaftar di sini
	err = config.DB.AutoMigrate(
		&models.User{},
		&models.Bird{},
		&models.Pair{},
		&models.Production{},
		&models.Egg{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi:", err)
	}

	// 3. Setup Router
	r := gin.Default()

	// --- SETUP CORS (SANGAT PENTING UNTUK COOKIE) ---
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- TAMBAHAN 1: AKSES FOLDER UPLOADS ---
	// Ini wajib agar Frontend bisa menampilkan foto profil (misal: http://localhost:8080/uploads/avatars/...)
	r.Static("/uploads", "./uploads") 

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Pong!", "status": "online"})
	})

	// --- AUTH ROUTES (PUBLIC - Tidak Butuh Login) ---
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/logout", controllers.Logout)
	}

	// --- PROTECTED ROUTES (BUTUH LOGIN) ---
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// --- TAMBAHAN 2: ROUTE UPDATE PROFILE ---
		// Route ini ditaruh di sini karena butuh login (ambil ID user dari token)
		protected.PUT("/auth/profile", controllers.UpdateProfile)

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

		// Production & Eggs
		protected.GET("/pairs/:pair_id/production", controllers.GetActiveProduction)
		protected.POST("/pairs/:pair_id/eggs", controllers.AddEgg)
		protected.PUT("/eggs/:egg_id/status", controllers.UpdateEggStatus)
		protected.DELETE("/eggs/:egg_id", controllers.DeleteEgg)

		// Finance
		protected.GET("/finance", controllers.GetTransactions)
		protected.GET("/finance/summary", controllers.GetFinanceSummary)
		protected.POST("/finance", controllers.CreateTransaction)
		protected.DELETE("/finance/:id", controllers.DeleteTransaction)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}