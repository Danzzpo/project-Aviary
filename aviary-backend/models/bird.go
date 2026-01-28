package models

import "time"

type Bird struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserID     uint   `gorm:"index;not null" json:"user_id"`
	
	// TAG JSON WAJIB HURUF KECIL
	RingNumber string `gorm:"type:varchar(50);not null" json:"ring_number"`
	Species    string `gorm:"type:varchar(50)" json:"species"`
	Mutation   string `gorm:"type:varchar(100)" json:"mutation"`
	Gender     string `gorm:"type:enum('M','F','UNKNOWN');default:'UNKNOWN'" json:"gender"`
	Status     string `gorm:"type:enum('AVAILABLE','PAIRED','SOLD','DEAD');default:'AVAILABLE'" json:"status"`
	
	SireID     *uint  `json:"sire_id"`
	DamID      *uint  `json:"dam_id"`
	
	// Relasi Self-Join
	Sire       *Bird  `gorm:"foreignKey:SireID;references:ID" json:"sire"`
	Dam        *Bird  `gorm:"foreignKey:DamID;references:ID" json:"dam"`

	DOB        *time.Time `json:"dob"`
	Notes      string     `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}