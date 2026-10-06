package repository

import (
	"calculate_lung_capacity/internal/app/ds"
)

func (r *Repository) CreateUser(username string, password string) (ds.PatientCategoryUser, error) {
	user := ds.PatientCategoryUser{
		PatientCategoryUsername: username,
		PatientCategoryPassword: password,
	}

	err := r.db.Create(&user).Error
	if err != nil {
		return ds.PatientCategoryUser{}, err
	}

	return user, nil
}

func (r *Repository) GetUserByUsername(username string) (ds.PatientCategoryUser, error) {
	var user ds.PatientCategoryUser

	err := r.db.
		Where("patient_category_username = ?", username).
		First(&user).Error
	if err != nil {
		return ds.PatientCategoryUser{}, err
	}

	return user, nil
}
