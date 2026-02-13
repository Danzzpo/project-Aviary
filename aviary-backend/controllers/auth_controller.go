package controllers

import (
	"aviary-backend/config"
	"aviary-backend/models"
	"aviary-backend/utils"
	"fmt"              // <--- BARU: Untuk format nama file
	"net/http"
	"os"               // <--- BARU: Untuk hapus file lama
	"path/filepath"    // <--- BARU: Untuk cek ekstensi file (.jpg/.png)
	"time"             // <--- BARU: Untuk cek batas 30 hari

	"github.com/gin-gonic/gin"
)

// Struct Validasi Register
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Struct Validasi Login
type LoginInput struct {
	Identity string `json:"email" binding:"required"` 
	Password string `json:"password" binding:"required"`
}

// Struct Validasi Update Profile (BARU)
// Kita pakai form, bukan json, karena ada upload file
type UpdateProfileInput struct {
	Username string `form:"username"`
	Email    string `form:"email"`
}

// 1. REGISTER
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

// 2. LOGIN
func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah: " + err.Error()})
		return
	}

	// Cari User
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

	// Generate Token
	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Set Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"user": gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"email":       user.Email,
			"profile_pic": user.ProfilePic, // Kirim info foto profil ke frontend
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

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", newAccess, 900, "/", "", false, true)
	c.SetCookie("refresh_token", newRefresh, 604800, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// 4. LOGOUT
func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout sukses"})
}

// 5. UPDATE PROFILE (FUNGSI BARU)
func UpdateProfile(c *gin.Context) {
	// A. Ambil ID User dari Middleware
	userID, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// B. Validasi Batas 30 Hari
	// Logic: Jika LastProfileUpdate ada isinya, kita cek selisihnya
	if !user.LastProfileUpdate.IsZero() {
		daysSinceUpdate := time.Since(user.LastProfileUpdate).Hours() / 24
		if daysSinceUpdate < 30 {
			remainingDays := 30 - int(daysSinceUpdate)
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Profil baru saja diubah. Tunggu %d hari lagi.", remainingDays),
			})
			return
		}
	}

	// C. Bind Input Form
	var input UpdateProfileInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// D. Cek Unik (Jika Username/Email diganti, pastikan tidak dipakai orang lain)
	var checkUser models.User
	if input.Username != "" && input.Username != user.Username {
		if err := config.DB.Where("username = ?", input.Username).First(&checkUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan orang lain!"})
			return
		}
		user.Username = input.Username
	}
	if input.Email != "" && input.Email != user.Email {
		if err := config.DB.Where("email = ?", input.Email).First(&checkUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan orang lain!"})
			return
		}
		user.Email = input.Email // Email lama otomatis tertimpa (terhapus)
	}

	// E. Handle Upload Foto
	file, err := c.FormFile("profile_pic")
	if err == nil {
		// 1. Validasi Ekstensi
		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format file harus jpg/png"})
			return
		}

		// 2. HAPUS FOTO LAMA (Agar Hemat Penyimpanan)
		if user.ProfilePic != "" {
			// user.ProfilePic misal: "/uploads/avatars/file.jpg"
			// Kita tambah "." di depan jadi: "./uploads/avatars/file.jpg"
			oldPath := "." + user.ProfilePic 
			os.Remove(oldPath) // Hapus file fisik
		}

		// 3. Simpan File Baru
		// Nama file unik: user_ID_TIMESTAMP.jpg
		filename := fmt.Sprintf("user_%d_%d%s", user.ID, time.Now().Unix(), ext)
		saveDir := "uploads/avatars"
		savePath := saveDir + "/" + filename
		
		// Buat folder jika belum ada
		os.MkdirAll(saveDir, os.ModePerm)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload foto"})
			return
		}

		// Simpan path relatif ke database (biar bisa diakses via URL)
		user.ProfilePic = "/" + savePath
	}

	// F. Simpan ke Database & Update Waktu
	user.LastProfileUpdate = time.Now()
	
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil berhasil diperbarui!",
		"user":    user,
	})
}