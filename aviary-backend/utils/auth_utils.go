package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// --- BAGIAN 1: PASSWORD HASHING ---

// HashPassword mengubah password asli menjadi kode acak
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash membandingkan password input dengan hash di database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// --- BAGIAN 2: TOKEN GENERATOR (DOUBLE TOKEN) ---

// GenerateTokenPair membuat DUA token sekaligus:
// 1. Access Token (Umur 15 Menit) -> Untuk request data sehari-hari
// 2. Refresh Token (Umur 7 Hari) -> Untuk minta access token baru jika habis
func GenerateTokenPair(userID uint) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "rahasia_super_aman_ganti_nanti" // Default secret
	}

	// A. Buat ACCESS TOKEN (Pendek)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // Expired 15 Menit
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tAccess, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// B. Buat REFRESH TOKEN (Panjang)
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // Expired 7 Hari
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	tRefresh, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return tAccess, tRefresh, nil
}

// ValidateToken mengecek apakah token asli atau palsu
func ValidateToken(tokenString string) (uint, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "rahasia_super_aman_ganti_nanti"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma enkripsi
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	// Ambil data user_id dari dalam token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if floatID, ok := claims["user_id"].(float64); ok {
			return uint(floatID), nil
		}
	}

	return 0, fmt.Errorf("invalid token claims")
}