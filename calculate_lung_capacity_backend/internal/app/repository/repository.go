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
		Where("patient_category_status = ?", ds.StatusPublished).
		Where("patient_category_age >= ? AND patient_category_age <= ?", ageMin, ageMax).
		Order("patient_category_id").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) GetLikesCount(categoryID uint) (int, error) {
	var count int64

	err := r.db.Model(&ds.PatientCategoryLike{}).
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
		Where("patient_category_id = ? AND patient_category_status <> ?", id, ds.StatusDeleted).
		First(&category).Error
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) GetFirstPublishedPatientCategory() (ds.PatientCategory, error) {
	var category ds.PatientCategory

	err := r.db.
		Where("patient_category_status = ?", ds.StatusPublished).
		Order("patient_category_id").
		First(&category).Error
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) GetNextPatientCategory(afterID int) (ds.PatientCategory, error) {
	var category ds.PatientCategory

	err := r.db.
		Where("patient_category_status = ? AND patient_category_id > ?", ds.StatusPublished, afterID).
		Order("patient_category_id").
		First(&category).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.GetFirstPublishedPatientCategory()
	}
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) GetDraftPatientCategory(userID uint) (ds.PatientCategory, error) {
	var category ds.PatientCategory

	err := r.db.
		Where("patient_category_status = ? AND patient_category_creator_id = ?", ds.StatusDraft, userID).
		First(&category).Error
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) CreateDraftPatientCategory(userID uint, title string, imageURL string, videoURL string) (ds.PatientCategory, error) {
	category := ds.PatientCategory{
		PatientCategoryTitle:     title,
		PatientCategoryStatus:    ds.StatusDraft,
		PatientCategoryImageURL:  imageURL,
		PatientCategoryVideoURL:  videoURL,
		PatientCategoryCreatedAt: time.Now(),
		PatientCategoryCreatorID: userID,
	}

	err := r.db.Create(&category).Error
	if err != nil {
		return ds.PatientCategory{}, err
	}

	return category, nil
}

func (r *Repository) PublishPatientCategory(id uint, userID uint, description string, age int, height int) error {
	result := r.db.Model(&ds.PatientCategory{}).
		Where("patient_category_id = ? AND patient_category_creator_id = ? AND patient_category_status = ?",
			id, userID, ds.StatusDraft).
		Updates(map[string]interface{}{
			"patient_category_description":  description,
			"patient_category_age":          age,
			"patient_category_height":       height,
			"patient_category_status":       ds.StatusPublished,
			"patient_category_published_at": sql.NullTime{Time: time.Now(), Valid: true},
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) DeletePatientCategory(id int, userID uint) error {
	query := `UPDATE patient_categories
		SET patient_category_status = $1
		WHERE patient_category_id = $2
			AND patient_category_creator_id = $3
			AND patient_category_status <> $1
		RETURNING patient_category_id`

	row := r.db.Raw(query, ds.StatusDeleted, id, userID).Row()

	var deletedID int
	err := row.Scan(&deletedID)
	if err != nil {
		return err
	}

	return nil
}
