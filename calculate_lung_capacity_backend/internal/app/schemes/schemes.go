package schemes

type PatientCategoryResponse struct {
	PatientCategoryID          uint   `json:"patient_category_id"`
	PatientCategoryTitle       string `json:"patient_category_title"`
	PatientCategoryDescription string `json:"patient_category_description"`
	PatientCategoryImageURL    string `json:"patient_category_image_url"`
	PatientCategoryVideoURL    string `json:"patient_category_video_url"`
	PatientCategoryAge         int    `json:"patient_category_age"`
	PatientCategoryHeight      int    `json:"patient_category_height"`
	PatientCategoryCreatedAt   string `json:"patient_category_created_at"`
	PatientCategoryPublishedAt string `json:"patient_category_published_at"`
	PatientCategoryCreatorID   uint   `json:"patient_category_creator_id"`
	LikesCount                 int    `json:"likes_count"`
	IsLiked                    int    `json:"is_liked"`
	IsMine                     int    `json:"is_mine"`
}

type PatientCategoryListResponse struct {
	PatientCategories []PatientCategoryResponse `json:"patient_categories"`
	DraftID           uint                      `json:"draft_id"`
}

type CreatePatientCategoryRequest struct {
	Title string `form:"title" binding:"required"`
}

type PublishPatientCategoryRequest struct {
	Description string `json:"description"`
	Age         int    `json:"age" binding:"required,min=1,max=120"`
	Height      int    `json:"height" binding:"required,min=1,max=250"`
}

type LikeRequest struct {
	Value int `json:"value" binding:"oneof=0 1"`
}

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=3,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	PatientCategoryUserID   uint   `json:"patient_category_user_id"`
	PatientCategoryUsername string `json:"patient_category_username"`
}
