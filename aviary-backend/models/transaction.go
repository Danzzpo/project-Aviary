package models

import "time"

type Transaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	
	// Tipe: INCOME (Masuk) atau EXPENSE (Keluar)
	Type        string    `gorm:"type:enum('INCOME','EXPENSE');not null" json:"type"`
	
	// Kategori: JUAL_BURUNG, PAKAN, OBAT, PERLENGKAPAN, LAINNYA
	Category    string    `gorm:"type:varchar(50);not null" json:"category"`
	
	Amount      float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Date        time.Time `json:"date"`
	Description string    `gorm:"type:text" json:"description"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}