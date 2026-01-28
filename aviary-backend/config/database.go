package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// KITA PAKAI SATU KONEKSI SAJA (DB_MASTER) AGAR STABIL
	dsn := os.Getenv("DB_MASTER_DSN")

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Debug Mode Aktif
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatal("❌ Gagal koneksi ke Database: ", err)
	}

	// Konfigurasi Connection Pool
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database Terhubung (Single Connection Mode)!")
}