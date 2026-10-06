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
	PatientCategoryID          uint         `gorm:"primaryKey;type:serial"`
	PatientCategoryTitle       string       `gorm:"type:varchar(150);not null"`
	PatientCategoryDescription string       `gorm:"type:varchar(500)"`
	PatientCategoryStatus      string       `gorm:"type:varchar(20);not null"`
	PatientCategoryImageURL    string       `gorm:"type:varchar(255)"`
	PatientCategoryVideoURL    string       `gorm:"type:varchar(255)"`
	PatientCategoryAge         int          `gorm:"type:integer"`
	PatientCategoryHeight      int          `gorm:"type:integer"`
	PatientCategoryCreatedAt   time.Time    `gorm:"type:timestamp;not null"`
	PatientCategoryPublishedAt sql.NullTime `gorm:"type:timestamp;default:null"`
	PatientCategoryCreatorID   uint         `gorm:"type:integer;not null"`

	Creator PatientCategoryUser `gorm:"foreignKey:PatientCategoryCreatorID;references:PatientCategoryUserID"`
}
