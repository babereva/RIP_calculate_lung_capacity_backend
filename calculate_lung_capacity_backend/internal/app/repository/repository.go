package repository

import (
	"database/sql"
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"calculate_lung_capacity/internal/app/ds"
)

const (
	AgeMin = 5
	AgeMax = 99
)

const CurrentPhysicianID = 1

type Repository struct {
	db *gorm.DB
}

func New(connectionString string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) GetPublishedPatientCategories(ageMin int, ageMax int) ([]ds.PatientCategory, error) {
	var categories []ds.PatientCategory

	err := r.db.
		Where("status = ?", ds.StatusPublished).
		Where("age >= ? AND age <= ?", ageMin, ageMax).
		Order("id").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) GetLikesCount(categoryID uint) (int, error) {
	var count int64

	err := r.db.Model(&ds.CategoryLike{}).
		Where("patient_category_id = ?", categoryID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *Repository) GetPatientCategoryByID(id int) (ds.PatientCategory, error) {
	var category ds.PatientCategory

	err := r.db.
		Where("id = ? AND status <> ?", id, ds.StatusDeleted).
		First(&category).Error
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) GetFirstPublishedPatientCategory() (ds.PatientCategory, error) {
	var category ds.PatientCategory

	err := r.db.
		Where("status = ?", ds.StatusPublished).
		Order("id").
		First(&category).Error
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) GetNextPatientCategory(afterID int) (ds.PatientCategory, error) {
	var category ds.PatientCategory

	err := r.db.
		Where("status = ? AND id > ?", ds.StatusPublished, afterID).
		Order("id").
		First(&category).Error

	// дошли до конца списка, начинаем сначала
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.GetFirstPublishedPatientCategory()
	}
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) GetDraftPatientCategory(physicianID uint) (ds.PatientCategory, error) {
	var category ds.PatientCategory

	err := r.db.
		Where("status = ? AND creator_id = ?", ds.StatusDraft, physicianID).
		First(&category).Error
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) CreateDraftPatientCategory(physicianID uint, title string, imageURL string, videoURL string) (ds.PatientCategory, error) {
	category := ds.PatientCategory{
		Title:     title,
		Status:    ds.StatusDraft,
		ImageURL:  imageURL,
		VideoURL:  videoURL,
		CreatedAt: time.Now(),
		CreatorID: physicianID,
	}

	err := r.db.Create(&category).Error
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) PublishPatientCategory(id uint, shortDescription string, age int, heightCm int) error {
	return r.db.Model(&ds.PatientCategory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"short_description": shortDescription,
			"age":               age,
			"height_cm":         heightCm,
			"status":            ds.StatusPublished,
			"formed_at":         sql.NullTime{Time: time.Now(), Valid: true},
		}).Error
}

func (r *Repository) DeletePatientCategory(id int) error {
	query := "UPDATE patient_categories SET status = $1 WHERE id = $2 RETURNING id"

	row := r.db.Raw(query, ds.StatusDeleted, id).Row()

	var deletedID int
	err := row.Scan(&deletedID)
	if err != nil {
		return err
	}

	return nil
}
