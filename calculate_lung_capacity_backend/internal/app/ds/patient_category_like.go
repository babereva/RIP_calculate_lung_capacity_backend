package ds

type PatientCategoryLike struct {
	PatientCategoryLikeID uint `gorm:"primaryKey;type:serial"`
	UserID                uint `gorm:"column:patient_category_user_id;type:integer;not null;uniqueIndex:idx_user_category"`
	CategoryID            uint `gorm:"column:patient_category_id;type:integer;not null;uniqueIndex:idx_user_category"`

	User     PatientCategoryUser `gorm:"foreignKey:UserID;references:PatientCategoryUserID"`
	Category PatientCategory     `gorm:"foreignKey:CategoryID;references:PatientCategoryID"`
}
