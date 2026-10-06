package repository

import (
	"errors"

	"gorm.io/gorm"

	"calculate_lung_capacity/internal/app/ds"
)

func (r *Repository) SetLike(userID uint, categoryID uint) error {
	var like ds.PatientCategoryLike

	err := r.db.
		Where("patient_category_user_id = ? AND patient_category_id = ?", userID, categoryID).
		First(&like).Error

	if err == nil {
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	newLike := ds.PatientCategoryLike{
		UserID:     userID,
		CategoryID: categoryID,
	}

	return r.db.Create(&newLike).Error
}

func (r *Repository) RemoveLike(userID uint, categoryID uint) error {
	return r.db.
		Where("patient_category_user_id = ? AND patient_category_id = ?", userID, categoryID).
		Delete(&ds.PatientCategoryLike{}).Error
}

func (r *Repository) IsLikedByUser(userID uint, categoryID uint) (bool, error) {
	var count int64

	err := r.db.Model(&ds.PatientCategoryLike{}).
		Where("patient_category_user_id = ? AND patient_category_id = ?", userID, categoryID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
