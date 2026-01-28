package controllers

import (
	"aviary-backend/config"
	"aviary-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardStats struct {
	TotalBirds     int64 `json:"total_birds"`
	TotalAvailable int64 `json:"total_available"`
	TotalSold      int64 `json:"total_sold"`
	TotalDeceased  int64 `json:"total_deceased"`

	ActivePairs    int64 `json:"active_pairs"`
	IncubatingEggs int64 `json:"incubating_eggs"`
}

func GetDashboardStats(c *gin.Context) {
	// 1. Ambil User ID
	userID, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uid := userID.(uint)

	var stats DashboardStats
	db := config.DB // Hapus .Debug() jika ingin terminal bersih

	// ==========================
	// 1. STATISTIK BURUNG
	// ==========================
	// Hitung hanya burung milik user ini
	db.Model(&models.Bird{}).Where("user_id = ?", uid).Count(&stats.TotalBirds)
	db.Model(&models.Bird{}).Where("user_id = ? AND status = ?", uid, "AVAILABLE").Count(&stats.TotalAvailable)
	db.Model(&models.Bird{}).Where("user_id = ? AND status = ?", uid, "SOLD").Count(&stats.TotalSold)
	db.Model(&models.Bird{}).Where("user_id = ? AND status = ?", uid, "DEAD").Count(&stats.TotalDeceased)

	// ==========================
	// 2. STATISTIK PASANGAN
	// ==========================
	// Hanya hitung pasangan yang masih ACTIVE
	db.Model(&models.Pair{}).Where("user_id = ? AND status = ?", uid, "ACTIVE").Count(&stats.ActivePairs)

	// ==========================
	// 3. STATISTIK TELUR (PERBAIKAN TOTAL)
	// ==========================
	// Logic Baru:
	// 1. Join ke Production & Pairs.
	// 2. Filter User ID (Wajib).
	// 3. Filter Status Telur (PENDING/FERTILE).
	// 4. FILTER PENTING: Pastikan Status Pasangan masih 'ACTIVE'.
	//    (Jika pasangan sudah dibubarkan/dihapus, telurnya jangan dihitung lagi di dashboard)
	
	db.Table("eggs").
		Joins("JOIN productions ON productions.id = eggs.production_id").
		Joins("JOIN pairs ON pairs.id = productions.pair_id").
		Where("pairs.user_id = ?", uid).                     // Kunci User ID
		Where("pairs.status = ?", "ACTIVE").                 // Kunci Pasangan Wajib Hidup
		Where("eggs.status IN ?", []string{"PENDING", "FERTILE"}).
		Count(&stats.IncubatingEggs)

	c.JSON(http.StatusOK, stats)
}