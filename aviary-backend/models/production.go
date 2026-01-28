package models

import "time"

// Model Produksi (Satu periode bertelur/Clutch)
type Production struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PairID    uint      `gorm:"not null" json:"pair_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Status    string    `gorm:"type:enum('ACTIVE','COMPLETED');default:'ACTIVE'" json:"status"`
	
	// Relasi ke Telur
	Eggs      []Egg     `gorm:"foreignKey:ProductionID" json:"eggs"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Model Telur
type Egg struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ProductionID uint       `gorm:"not null" json:"production_id"`
	
	EggOrder     int        `json:"egg_order"` // Telur ke-1, ke-2, dst
	LaidDate     time.Time  `json:"laid_date"` // Tanggal Keluar
	EstHatchDate time.Time  `json:"est_hatch_date"` // Estimasi Menetas
	HatchDate    *time.Time `json:"hatch_date"` // Tanggal Menetas Asli
	
	// Status Telur
	Status       string     `gorm:"type:enum('PENDING','FERTILE','INFERTILE','DIS','HATCHED','BROKEN');default:'PENDING'" json:"status"`
	
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}