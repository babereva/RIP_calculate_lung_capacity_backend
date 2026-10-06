package ds

type PatientCategoryUser struct {
	PatientCategoryUserID   uint   `gorm:"primaryKey;type:serial"`
	PatientCategoryUsername string `gorm:"type:varchar(50);unique;not null"`
	PatientCategoryPassword string `gorm:"type:varchar(50);not null"`
}
