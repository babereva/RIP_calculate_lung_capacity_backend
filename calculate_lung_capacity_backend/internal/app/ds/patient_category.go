package ds

import (
	"database/sql"
	"time"
)

const (
	StatusDraft     = "черновик"
	StatusPublished = "опубликован"
	StatusDeleted   = "удалён"
)


type PatientCategory struct {
	ID               uint   `gorm:"primaryKey"`
	Title            string `gorm:"type:varchar(100);not null"`
	ShortDescription string `gorm:"type:varchar(255)"`
	Status           string `gorm:"type:varchar(15);not null"`
	ImageURL         string `gorm:"type:varchar(255)"`
	VideoURL         string `gorm:"type:varchar(255)"`
	Age      int
	HeightCm int
	CreatedAt time.Time    `gorm:"not null"`
	FormedAt  sql.NullTime `gorm:"default:null"`
	CreatorID uint         `gorm:"not null"`
	Creator Physician `gorm:"foreignKey:CreatorID"`
}
