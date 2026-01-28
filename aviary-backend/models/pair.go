package models

import "time"

type Pair struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	CageName    string    `gorm:"type:varchar(50)" json:"cage_name"`
	
	SireID      uint      `gorm:"not null" json:"sire_id"`
	DamID       uint      `gorm:"not null" json:"dam_id"`
	
	// BAGIAN INI SANGAT PENTING:
	Sire        *Bird     `gorm:"foreignKey:SireID;references:ID" json:"sire"`
	Dam         *Bird     `gorm:"foreignKey:DamID;references:ID" json:"dam"`

	PairingDate time.Time `json:"pairing_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      string    `gorm:"type:enum('ACTIVE','HISTORY');default:'ACTIVE'" json:"status"`
	
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}