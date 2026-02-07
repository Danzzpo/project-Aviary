package controllers

import (
	"aviary-backend/config"
	"aviary-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionInput struct {
	Type        string  `json:"type" binding:"required"` // INCOME / EXPENSE
	Category    string  `json:"category" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Date        string  `json:"date" binding:"required"`
	Description string  `json:"description"`
}

// GET /api/finance
func GetTransactions(c *gin.Context) {
	userID, _ := c.Get("currentUser")
	var transactions []models.Transaction

	// Urutkan dari tanggal terbaru
	if err := config.DB.Where("user_id = ?", userID).Order("date desc, id desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

// GET /api/finance/summary (Hitung Saldo)
func GetFinanceSummary(c *gin.Context) {
	userID, _ := c.Get("currentUser")
	
	type Result struct {
		Type  string
		Total float64
	}
	var results []Result

	// Query SUM berdasarkan Tipe (INCOME/EXPENSE)
	config.DB.Model(&models.Transaction{}).
		Select("type, SUM(amount) as total").
		Where("user_id = ?", userID).
		Group("type").
		Scan(&results)

	var income, expense float64
	for _, r := range results {
		if r.Type == "INCOME" {
			income = r.Total
		} else {
			expense = r.Total
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"income":  income,
		"expense": expense,
		"balance": income - expense, // Saldo Bersih
	})
}

// POST /api/finance
func CreateTransaction(c *gin.Context) {
	userID, _ := c.Get("currentUser")
	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedDate, _ := time.Parse("2006-01-02", input.Date)

	trx := models.Transaction{
		UserID:      userID.(uint),
		Type:        input.Type,
		Category:    input.Category,
		Amount:      input.Amount,
		Date:        parsedDate,
		Description: input.Description,
	}

	if err := config.DB.Create(&trx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan transaksi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": trx})
}

// DELETE /api/finance/:id
func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Transaction{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Terhapus"})
}