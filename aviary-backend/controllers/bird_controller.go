package controllers

import (
	"aviary-backend/config"
	"aviary-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- STRUCTS UNTUK VALIDASI INPUT ---

type CreateBirdInput struct {
	RingNumber string `json:"ring_number" binding:"required"`
	Species    string `json:"species" binding:"required"`
	Mutation   string `json:"mutation"`
	Gender     string `json:"gender"`
	SireID     *uint  `json:"sire_id"`
	DamID      *uint  `json:"dam_id"`
}

type UpdateBirdInput struct {
	RingNumber string `json:"ring_number"`
	Species    string `json:"species"`
	Mutation   string `json:"mutation"`
	Gender     string `json:"gender"`
	Status     string `json:"status"` // Penting untuk update status (SOLD/DEAD)
	SireID     *uint  `json:"sire_id"`
	DamID      *uint  `json:"dam_id"`
}

// --- HANDLER FUNCTIONS ---

// 1. GET /api/birds (Lihat Semua)
func GetBirds(c *gin.Context) {
	userID, _ := c.Get("currentUser")
	var birds []models.Bird

	// Query ke DB Slave (Read Replica)
	// Preload Sire & Dam agar data bapak/ibu ikut terbawa
	if err := config.DB.Where("user_id = ?", userID).Preload("Sire").Preload("Dam").Find(&birds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": birds})
}

// 2. POST /api/birds (Tambah Baru)
func CreateBird(c *gin.Context) {
	// Ambil UserID dari Middleware
	userID, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input CreateBirdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buat Objek Model
	bird := models.Bird{
		UserID:     userID.(uint),
		RingNumber: input.RingNumber,
		Species:    input.Species,
		Mutation:   input.Mutation,
		Gender:     input.Gender,
		SireID:     input.SireID, // ID Bapak
		DamID:      input.DamID,  // ID Ibu
		Status:     "AVAILABLE",  // Default saat buat baru
	}

	// Simpan ke DB Master
	if err := config.DB.Create(&bird).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data burung"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": bird})
}

// 3. PUT /api/birds/:id (Edit Data)
func UpdateBird(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("currentUser")

	var bird models.Bird
	// Cek apakah burung ini milik user yang login?
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&bird).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Burung tidak ditemukan"})
		return
	}

	var input UpdateBirdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update Field Manual (Agar aman dan terkontrol)
	bird.RingNumber = input.RingNumber
	bird.Species = input.Species
	bird.Mutation = input.Mutation
	bird.Gender = input.Gender
	bird.Status = input.Status // Update status (misal jadi SOLD)
	bird.SireID = input.SireID // Update Bapak
	bird.DamID = input.DamID   // Update Ibu

	// Simpan Perubahan ke DB Master
	config.DB.Save(&bird)

	c.JSON(http.StatusOK, gin.H{"data": bird, "message": "Berhasil diupdate"})
}