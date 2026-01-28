package controllers

import (
	"aviary-backend/config"
	"aviary-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Input Structs
type CreatePairInput struct {
	CageName string `json:"cage_name" binding:"required"`
	SireID   uint   `json:"sire_id" binding:"required"`
	DamID    uint   `json:"dam_id" binding:"required"`
	Date     string `json:"date"`
}

type AddEggInput struct {
	LaidDate string `json:"laid_date" binding:"required"`
}

type UpdateEggStatusInput struct {
	Status string `json:"status" binding:"required"`
}

// ==========================================
// 1. FITUR UTAMA: PASANGAN (PAIRING)
// ==========================================

// GET /api/pairs
func GetActivePairs(c *gin.Context) {
	userID, _ := c.Get("currentUser")
	var pairs []models.Pair

	// LANGKAH 1: Ambil Data Pasangan Saja (Tanpa Preload GORM)
	// Kita ambil manual agar lebih stabil
	if err := config.DB.Where("user_id = ? AND status = ?", userID, "ACTIVE").Find(&pairs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pasangan"})
		return
	}

	// LANGKAH 2: ISI DATA BURUNG SECARA MANUAL (MANUAL INJECTION)
	// Loop satu per satu untuk memastikan data terambil meskipun statusnya aneh di DB
	for i := range pairs {
		// Cari Sire (Bapak) - Pakai Unscoped agar tembus data terhapus
		var sire models.Bird
		if err := config.DB.Unscoped().Where("id = ?", pairs[i].SireID).First(&sire).Error; err == nil {
			pairs[i].Sire = &sire
		}

		// Cari Dam (Ibu) - Pakai Unscoped agar tembus data terhapus
		var dam models.Bird
		if err := config.DB.Unscoped().Where("id = ?", pairs[i].DamID).First(&dam).Error; err == nil {
			pairs[i].Dam = &dam
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": pairs})
}

// POST /api/pairs
func CreatePair(c *gin.Context) {
	userID, _ := c.Get("currentUser")
	var input CreatePairInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pairDate := time.Now()
	if input.Date != "" {
		parsed, err := time.Parse("2006-01-02", input.Date)
		if err == nil {
			pairDate = parsed
		}
	}

	tx := config.DB.Begin()

	// Cek Ketersediaan Burung (Manual Check)
	var count int64
	tx.Model(&models.Bird{}).Where("id IN ? AND status = ? AND user_id = ?", []uint{input.SireID, input.DamID}, "AVAILABLE", userID).Count(&count)

	if count != 2 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Salah satu burung tidak tersedia / bukan milik Anda"})
		return
	}

	// Simpan Pasangan
	pair := models.Pair{
		UserID:      userID.(uint),
		CageName:    input.CageName,
		SireID:      input.SireID,
		DamID:       input.DamID,
		PairingDate: pairDate,
		Status:      "ACTIVE",
	}

	if err := tx.Create(&pair).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan pasangan"})
		return
	}

	// Update Status Burung jadi PAIRED
	tx.Model(&models.Bird{}).Where("id IN ?", []uint{input.SireID, input.DamID}).Update("status", "PAIRED")

	// Buat Produksi Awal
	production := models.Production{PairID: pair.ID, StartDate: pairDate, Status: "ACTIVE"}
	tx.Create(&production)

	tx.Commit()

	// LOAD MANUAL UNTUK RESPON BALIK (Agar Frontend tidak error)
	var sire, dam models.Bird
	config.DB.Unscoped().First(&sire, pair.SireID)
	config.DB.Unscoped().First(&dam, pair.DamID)
	pair.Sire = &sire
	pair.Dam = &dam

	c.JSON(http.StatusCreated, gin.H{"data": pair, "message": "Berhasil menjodohkan!"})
}

// PUT /api/pairs/:id/disband
func DisbandPair(c *gin.Context) {
	pairID := c.Param("id")
	userID, _ := c.Get("currentUser")
	tx := config.DB.Begin()

	var pair models.Pair
	if err := tx.Where("id = ? AND user_id = ? AND status = ?", pairID, userID, "ACTIVE").First(&pair).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ditemukan"})
		return
	}

	now := time.Now()
	// Update Pasangan jadi HISTORY
	tx.Model(&pair).Updates(map[string]interface{}{"status": "HISTORY", "end_date": now})
	// Kembalikan Burung jadi AVAILABLE
	tx.Model(&models.Bird{}).Where("id IN ?", []uint{pair.SireID, pair.DamID}).Update("status", "AVAILABLE")
	// Tutup Produksi
	tx.Model(&models.Production{}).Where("pair_id = ? AND status = ?", pairID, "ACTIVE").Update("status", "COMPLETED")

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Dibubarkan"})
}

// ==========================================
// 2. FITUR TELUR & PRODUKSI
// ==========================================

func GetActiveProduction(c *gin.Context) {
	pairID := c.Param("pair_id")
	var production models.Production
	// Ambil produksi aktif + telur
	if err := config.DB.Preload("Eggs").Where("pair_id = ? AND status = ?", pairID, "ACTIVE").First(&production).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": production})
}

func AddEgg(c *gin.Context) {
	pairID := c.Param("pair_id")
	var input AddEggInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari Produksi Aktif, kalau tidak ada buat baru
	var production models.Production
	if err := config.DB.Where("pair_id = ? AND status = ?", pairID, "ACTIVE").First(&production).Error; err != nil {
		production = models.Production{PairID: stringToUint(pairID), StartDate: time.Now(), Status: "ACTIVE"}
		config.DB.Create(&production)
	}

	// Hitung urutan telur
	var count int64
	config.DB.Model(&models.Egg{}).Where("production_id = ?", production.ID).Count(&count)

	laidDate, _ := time.Parse("2006-01-02", input.LaidDate)
	egg := models.Egg{
		ProductionID: production.ID,
		EggOrder:     int(count) + 1,
		LaidDate:     laidDate,
		EstHatchDate: laidDate.AddDate(0, 0, 21), // Estimasi 21 hari
		Status:       "PENDING",
	}

	if err := config.DB.Create(&egg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan telur"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": egg})
}

func UpdateEggStatus(c *gin.Context) {
	eggID := c.Param("egg_id")
	var input UpdateEggStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var egg models.Egg
	if err := config.DB.First(&egg, eggID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Telur tidak ditemukan"})
		return
	}
	
	egg.Status = input.Status
	// Jika menetas, catat tanggalnya
	if input.Status == "HATCHED" {
		now := time.Now()
		egg.HatchDate = &now
	}
	
	config.DB.Save(&egg)
	c.JSON(http.StatusOK, gin.H{"data": egg})
}

// DELETE /api/eggs/:egg_id (FITUR BARU UNTUK HAPUS TELUR ZONK)
func DeleteEgg(c *gin.Context) {
	eggID := c.Param("egg_id")
	
	// Hapus permanen (Unscoped) agar database bersih
	if err := config.DB.Unscoped().Delete(&models.Egg{}, eggID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus telur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Telur berhasil dihapus"})
}

// Helper
func stringToUint(s string) uint {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(val)
}