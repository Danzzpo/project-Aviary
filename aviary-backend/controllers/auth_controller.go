package controllers

import (
	"aviary-backend/config"
	"aviary-backend/models"
	"aviary-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Struct Validasi
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// --- PERUBAHAN 1: STRUCT LOGIN ---
// Kita ubah nama field jadi "Identity" dan HAPUS binding "email"
// agar bisa menerima inputan berupa Username biasa.
type LoginInput struct {
	Identity string `json:"email" binding:"required"` // <--- HAPUS ",email" DI SINI
	Password string `json:"password" binding:"required"`
}

// 1. REGISTER (User Baru)
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hash password"})
		return
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         "breeder",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau Username sudah terdaftar!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registrasi berhasil!", "data": user.Username})
}

// 2. LOGIN (DIGANTI TOTAL AGAR BISA USERNAME/EMAIL)
func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	// Validasi Input (Sekarang Username biasa tidak akan kena Error 400)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah: " + err.Error()})
		return
	}

	// --- PERUBAHAN 2: QUERY DATABASE ---
	// Logika: Cari user yang Email-nya cocok ATAU Username-nya cocok
	if err := config.DB.Where("email = ? OR username = ?", input.Identity, input.Identity).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akun tidak ditemukan"})
		return
	}

	// Cek Password
	match := utils.CheckPasswordHash(input.Password, user.PasswordHash)
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// Generate 2 Token
	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Set Cookie (HttpOnly)
	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/", "", false, true)

	// Response JSON
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// 3. REFRESH TOKEN
func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session habis, login ulang"})
		return
	}

	userID, err := utils.ValidateToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	newAccess, newRefresh, _ := utils.GenerateTokenPair(userID)

	c.SetCookie("access_token", newAccess, 900, "/", "", false, true)
	c.SetCookie("refresh_token", newRefresh, 604800, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// 4. LOGOUT
func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout sukses"})
}