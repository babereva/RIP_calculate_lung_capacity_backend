package repository

import (
	"gorm.io/gorm"

	"calculate_lung_capacity/internal/app/ds"
)

func (r *Repository) UpdatePatientCategoryMedia(id uint, userID uint, imageName string, videoName string) error {
	fields := map[string]interface{}{}

	if imageName != "" {
		fields["patient_category_image_url"] = imageName
	}

	if videoName != "" {
		fields["patient_category_video_url"] = videoName
	}

	if len(fields) == 0 {
		return nil
	}

	result := r.db.Model(&ds.PatientCategory{}).
		Where("patient_category_id = ? AND patient_category_creator_id = ?", id, userID).
		Updates(fields)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
