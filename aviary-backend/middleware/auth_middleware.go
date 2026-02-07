package middleware

import (
	"aviary-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Coba ambil Access Token dari Cookie
		tokenString, err := c.Cookie("access_token")

		// 2. Jika Cookie kosong atau tidak ada
		if err != nil || tokenString == "" {
			// Kita kirim kode 401.
			// Kode 401 ini adalah KODE RAHASIA bagi Frontend untuk melakukan Refresh Token otomatis.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// 3. Jika ada, validasi isinya
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		// 4. Jika valid, simpan UserID ke context agar bisa dipakai Controller lain
		c.Set("currentUser", userID)
		c.Next()
	}
}