package repository

import "fmt"


const MinioBaseURL = "http://localhost:9000/lung-capacity-media"


const (
	StatusDraft     = "черновик"
	StatusPublished = "опубликован"
	StatusDeleted   = "удалён"
)


const (
	SexMale   = "мужской"
	SexFemale = "женский"
)

type PatientCategory struct {
	ID                   int
	CategoryTitle        string
	CategoryDescription  string
	Sex                  string
	AgeFrom              int
	AgeTo                int
	HeightFromCm         int
	HeightToCm           int
	CoefficientA         float64
	CoefficientB         float64
	CoefficientC         float64
	PreviewImageKey      string
	LungCapacityVideoKey string
	CategoryStatus       string
	LikedByPhysicianIDs  []int
}

type Repository struct {
	patientCategories []PatientCategory
}

func NewRepository() (*Repository, error) {
	return &Repository{patientCategories: []PatientCategory{
		{
			ID:                   1,
			CategoryTitle:        "Расчет ДЖЕЛ для женщин",
			CategoryDescription:  "Молодые женщины без хронических заболеваний органов дыхания и без стажа курения. Измерение проводится в положении стоя, выполняются три попытки, в расчёт берётся максимальный результат. Значение не ниже 85 % от должного считается нормой, диапазон 70–85 % расценивается как умеренное снижение, ниже 70 % — как выраженное.",
			Sex:                  SexFemale,
			AgeFrom:              18,
			AgeTo:                25,
			HeightFromCm:         160,
			HeightToCm:           165,
			CoefficientA:         0.049,
			CoefficientB:         -0.019,
			CoefficientC:         -3.76,
			PreviewImageKey:      "woman_18-25_160-165.jpg",
			LungCapacityVideoKey: "woman_18-25_160-165.mp4",
			CategoryStatus:       StatusPublished,
			LikedByPhysicianIDs:  []int{1, 3, 7, 12},
		},
		{
			ID:                   2,
			CategoryTitle:        "Расчет ДЖЕЛ для женщин",
			CategoryDescription:  "Категория применима к пациенткам без хронических заболеваний органов дыхания и без стажа курения. Результат измерения считается нормальным при значении не ниже 85 % от должного. Диапазон 70–85 % расценивается как умеренное снижение, ниже 70 % — как выраженное.",
			Sex:                  SexFemale,
			AgeFrom:              30,
			AgeTo:                33,
			HeightFromCm:         170,
			HeightToCm:           175,
			CoefficientA:         0.049,
			CoefficientB:         -0.019,
			CoefficientC:         -3.76,
			PreviewImageKey:      "woman_30-33_170-175.jpg",
			LungCapacityVideoKey: "woman_30-33_170-175.mp4",
			CategoryStatus:       StatusPublished,
			LikedByPhysicianIDs:  []int{2, 4, 5, 9, 11, 14, 18},
		},
		{
			ID:                   3,
			CategoryTitle:        "Расчет ДЖЕЛ для женщин",
			CategoryDescription:  "Женщины среднего возраста. После сорока лет должная жизненная ёмкость лёгких закономерно снижается, поэтому поправка на возраст в уравнении обязательна. Перед измерением требуется покой не менее пятнадцати минут.",
			Sex:                  SexFemale,
			AgeFrom:              42,
			AgeTo:                45,
			HeightFromCm:         165,
			HeightToCm:           170,
			CoefficientA:         0.049,
			CoefficientB:         -0.019,
			CoefficientC:         -3.76,
			PreviewImageKey:      "woman_42-45_165-170.jpg",
			LungCapacityVideoKey: "woman_42-45_165-170.mp4",
			CategoryStatus:       StatusPublished,
			LikedByPhysicianIDs:  []int{6, 8, 13},
		},
		{
			ID:                   4,
			CategoryTitle:        "Расчет ДЖЕЛ для женщин",
			CategoryDescription:  "Пациентки пожилого возраста. Учитывается снижение подвижности грудной клетки и уменьшение роста с возрастом, поэтому рост измеряется непосредственно перед исследованием, а не берётся из анамнеза.",
			Sex:                  SexFemale,
			AgeFrom:              62,
			AgeTo:                65,
			HeightFromCm:         155,
			HeightToCm:           160,
			CoefficientA:         0.049,
			CoefficientB:         -0.019,
			CoefficientC:         -3.76,
			PreviewImageKey:      "woman_62-65_155-160.jpg",
			LungCapacityVideoKey: "woman_62-65_155-160.mp4",
			CategoryStatus:       StatusPublished,
			LikedByPhysicianIDs:  []int{3, 10},
		},
		{
			ID:                   5,
			CategoryTitle:        "Расчет ДЖЕЛ для мужчин",
			CategoryDescription:  "Мужчины молодого возраста без профессиональных вредностей. У мужчин коэффициент при росте выше, чем у женщин, за счёт большего объёма грудной клетки, поэтому используется отдельный набор коэффициентов.",
			Sex:                  SexMale,
			AgeFrom:              26,
			AgeTo:                29,
			HeightFromCm:         170,
			HeightToCm:           175,
			CoefficientA:         0.052,
			CoefficientB:         -0.028,
			CoefficientC:         -3.20,
			PreviewImageKey:      "man_26-29_170-175.jpg",
			LungCapacityVideoKey: "man_26-29_170-175.mp4",
			CategoryStatus:       StatusPublished,
			LikedByPhysicianIDs:  []int{1, 5, 16, 17, 20},
		},
		{
			ID:                   6,
			CategoryTitle:        "Расчет ДЖЕЛ для мужчин",
			CategoryDescription:  "Мужчины старшего возраста. Категория используется при профилактических осмотрах: сопоставление фактического результата с должным позволяет выявить снижение вентиляционной способности лёгких до появления жалоб.",
			Sex:                  SexMale,
			AgeFrom:              54,
			AgeTo:                57,
			HeightFromCm:         175,
			HeightToCm:           180,
			CoefficientA:         0.052,
			CoefficientB:         -0.028,
			CoefficientC:         -3.20,
			PreviewImageKey:      "man_54-57_175-180.jpg",
			LungCapacityVideoKey: "man_54-57_175-180.mp4",
			CategoryStatus:       StatusPublished,
			LikedByPhysicianIDs:  []int{7, 19},
		},
		{
			ID:                   7,
			CategoryTitle:        "Расчет ДЖЕЛ для мальчиков",
			CategoryDescription:  "Мальчики младшего школьного возраста. Для детей применяется отдельное регрессионное уравнение, зависящее только от роста, поэтому коэффициент при возрасте равен нулю.",
			Sex:                  SexMale,
			AgeFrom:              9,
			AgeTo:                11,
			HeightFromCm:         130,
			HeightToCm:           135,
			CoefficientA:         0.0453,
			CoefficientB:         0,
			CoefficientC:         -3.90,
			PreviewImageKey:      "boy_9-11_130-135.jpg",
			LungCapacityVideoKey: "boy_9-11_130-135.mp4",
			CategoryStatus:       StatusDraft,
			LikedByPhysicianIDs:  []int{},
		},
		{
			ID:                   8,
			CategoryTitle:        "Расчет ДЖЕЛ для девочек",
			CategoryDescription:  "Девочки дошкольного и младшего школьного возраста. Категория выведена из обращения: уравнение уточняется по новым нормативам.",
			Sex:                  SexFemale,
			AgeFrom:              7,
			AgeTo:                9,
			HeightFromCm:         120,
			HeightToCm:           125,
			CoefficientA:         0.0375,
			CoefficientB:         0,
			CoefficientC:         -3.15,
			PreviewImageKey:      "girl_7-9_120-125.jpg",
			LungCapacityVideoKey: "girl_7-9_120-125.mp4",
			CategoryStatus:       StatusDeleted,
			LikedByPhysicianIDs:  []int{},
		},
	}}, nil
}


func (r *Repository) GetPublishedPatientCategories(patientAge int) ([]PatientCategory, error) {
	result := []PatientCategory{}
	for _, category := range r.patientCategories {
		if category.CategoryStatus != StatusPublished {
			continue
		}
		if patientAge > 0 && (patientAge < category.AgeFrom || patientAge > category.AgeTo) {
			continue
		}
		result = append(result, category)
	}
	return result, nil
}


func (r *Repository) GetPatientCategoryByID(id int) (PatientCategory, error) {
	for _, category := range r.patientCategories {
		if category.ID == id && category.CategoryStatus != StatusDeleted {
			return category, nil
		}
	}
	return PatientCategory{}, fmt.Errorf("категория пациентов %d не найдена", id)
}


func (r *Repository) GetNextPatientCategory(afterID int) (PatientCategory, error) {
	published, err := r.GetPublishedPatientCategories(0)
	if err != nil {
		return PatientCategory{}, err
	}
	if len(published) == 0 {
		return PatientCategory{}, fmt.Errorf("нет опубликованных категорий пациентов")
	}
	for i, category := range published {
		if category.ID == afterID {
			return published[(i+1)%len(published)], nil
		}
	}
	return published[0], nil
}

func (r *Repository) GetDraftPatientCategory() (PatientCategory, error) {
	for _, category := range r.patientCategories {
		if category.CategoryStatus == StatusDraft {
			return category, nil
		}
	}
	return PatientCategory{}, fmt.Errorf("черновик категории пациентов не найден")
}