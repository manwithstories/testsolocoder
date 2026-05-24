package dto

type CreateProjectRequest struct {
	Name        string  `json:"name" binding:"required,max=128"`
	Description string  `json:"description"`
	OwnerID     uint    `json:"owner_id" binding:"required"`
	CoverImage  string  `json:"cover_image"`
	RoomType    string  `json:"room_type"`
	Area        float64 `json:"area"`
	Budget      float64 `json:"budget"`
}

type UpdateProjectRequest struct {
	Name        string   `json:"name" binding:"omitempty,max=128"`
	Description string   `json:"description"`
	Status      string   `json:"status" binding:"omitempty,oneof=draft submitted approved rejected"`
	CoverImage  string   `json:"cover_image"`
	RoomType    string   `json:"room_type"`
	Area        *float64 `json:"area"`
	Budget      *float64 `json:"budget"`
}

type ProjectListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status" binding:"omitempty,oneof=draft submitted approved rejected"`
	Role     string `form:"-"`
	UserID   uint   `form:"-"`
}

type UploadImageRequest struct {
	ImageURL    string `json:"image_url" binding:"required"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

type AddCommentRequest struct {
	Content   string   `json:"content" binding:"required"`
	Type      string   `json:"type" binding:"required,oneof=text marker"`
	PositionX *float64 `json:"position_x"`
	PositionY *float64 `json:"position_y"`
	ParentID  uint     `json:"parent_id"`
}

type ProjectResponse struct {
	ID          uint    `json:"id"`
	DesignerID  uint    `json:"designer_id"`
	OwnerID     uint    `json:"owner_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	CoverImage  string  `json:"cover_image"`
	RoomType    string  `json:"room_type"`
	Area        float64 `json:"area"`
	Budget      float64 `json:"budget"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ProjectDetailResponse struct {
	ProjectResponse
	Images []ImageResponse `json:"images"`
}

type ImageResponse struct {
	ID          uint   `json:"id"`
	ProjectID   uint   `json:"project_id"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	CreatedAt   string `json:"created_at"`
}

type ProjectListResponse struct {
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	List     []ProjectResponse  `json:"list"`
}

type CommentResponse struct {
	ID        uint    `json:"id"`
	ProjectID uint    `json:"project_id"`
	UserID    uint    `json:"user_id"`
	UserRole  string  `json:"user_role"`
	Content   string  `json:"content"`
	Type      string  `json:"type"`
	PositionX float64 `json:"position_x"`
	PositionY float64 `json:"position_y"`
	ParentID  uint    `json:"parent_id"`
	CreatedAt string  `json:"created_at"`
}

type CommentListResponse struct {
	List []CommentResponse `json:"list"`
}
