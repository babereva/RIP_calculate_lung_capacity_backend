package ds


type CategoryLike struct {
	ID                uint `gorm:"primaryKey"`
	PhysicianID       uint `gorm:"not null;uniqueIndex:idx_physician_category"`
	PatientCategoryID uint `gorm:"not null;uniqueIndex:idx_physician_category"`
	Physician       Physician       `gorm:"foreignKey:PhysicianID"`
	PatientCategory PatientCategory `gorm:"foreignKey:PatientCategoryID"`
}
