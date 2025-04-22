package controller

import "time"

// ========================= 通用响应结构 =========================

// 基础响应结构
type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"操作成功"`
	Data    interface{} `json:"data,omitempty"`
}

// 错误响应
type ErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message" example:"发生错误"`
	Data    interface{} `json:"data,omitempty"`
}

// 通用分页数据结构
type PagedData struct {
	Total int64 `json:"total" example:"100"`
	Page  int   `json:"page" example:"1"`
	Limit int   `json:"limit" example:"10"`
}

// ========================= 用户相关 =========================

// 用户DTO
type UserDTO struct {
	Id          uint      `json:"id" example:"1"`
	Username    string    `json:"username" example:"user123"`
	DisplayName string    `json:"display_name" example:"张三"`
	Role        int       `json:"role" example:"1"`
	Status      int       `json:"status" example:"1"`
	Email       string    `json:"email" example:"user@example.com"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"user123"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// 登录响应
type LoginResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message,omitempty" example:"登录成功"`
	Data    TokenDTO `json:"data"`
}

// Token数据
type TokenDTO struct {
	Token string  `json:"token" example:"eyJhbGciOiJIUzI..."`
	User  UserDTO `json:"user"`
}

// 用户列表响应
type UsersResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message,omitempty"`
	Data    UsersListData `json:"data"`
}

// 用户列表数据
type UsersListData struct {
	Users []UserDTO `json:"users"`
	PagedData
}

// ========================= 项目相关 =========================

// 项目成员请求
type ProjectMemberRequest struct {
	ProjectID uint `json:"projectID" binding:"required" example:"1"`
	UserID    uint `json:"userID" binding:"required" example:"2"`
	Role      int  `json:"role" binding:"required" example:"1"`
}

// 项目成员DTO
type MemberDTO struct {
	ProjectID uint   `json:"projectId" example:"1"`
	UserID    uint   `json:"userId" example:"2"`
	Username  string `json:"username" example:"user123"`
	Role      int    `json:"role" example:"1"`
}

// 项目成员关系DTO
type ProjectMembershipDTO struct {
	ID        uint    `json:"id" example:"1"`
	UserID    uint    `json:"userId" example:"2"`
	ProjectID uint    `json:"projectId" example:"1"`
	Role      int     `json:"role" example:"1"`
	User      UserDTO `json:"user"`
}

// 项目成员响应
type ProjectMembershipResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// 角色更新请求
type RoleUpdateRequest struct {
	Role int `json:"role" binding:"required" example:"2"`
}

// 项目DTO
type ProjectDTO struct {
	ID          uint        `json:"id" example:"1"`
	Name        string      `json:"name" example:"项目名称"`
	Description string      `json:"description" example:"项目描述"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
	Users       []MemberDTO `json:"users,omitempty"`
}

// 项目列表响应
type ProjectsResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message,omitempty"`
	Data    ProjectListData `json:"data"`
}

// 项目列表数据
type ProjectListData struct {
	Projects []ProjectDTO `json:"projects"`
	PagedData
}

// 项目响应
type ProjectResponse struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message,omitempty"`
	Data    ProjectDTO `json:"data"`
}

// ========================= 数据集相关 =========================

// 数据集DTO
type DatasetDTO struct {
	ID               uint      `json:"id" example:"1"`
	Name             string    `json:"name" example:"数据集名称"`
	Description      string    `json:"description" example:"数据集描述"`
	StorageType      string    `json:"storage_type" example:"database"`
	TemplateType     string    `json:"template_type" example:"instruction_io"`
	EntryCount       int64     `json:"entry_count" example:"100"`
	TotalSize        int64     `json:"total_size" example:"1024"`
	ProjectID        uint      `json:"project_id" example:"1"`
	UserID           uint      `json:"user_id" example:"1"`
	SchemaDefinition string    `json:"schema_definition,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

// 数据集响应
type DatasetResponse struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message" example:""`
	Data    DatasetDTO `json:"data"`
}

// 数据集列表响应
type DatasetsResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:""`
	Data    DatasetsListData `json:"data"`
}

// 数据集列表数据
type DatasetsListData struct {
	Datasets []DatasetDTO `json:"datasets"`
	PagedData
}

// 数据集条目DTO
type DatasetEntryDTO struct {
	ID          uint      `json:"id" example:"1"`
	DatasetID   uint      `json:"dataset_id" example:"1"`
	EntryIndex  int       `json:"entry_index" example:"0"`
	Instruction string    `json:"instruction" example:"指令内容"`
	Input       string    `json:"input" example:"输入内容"`
	Output      string    `json:"output" example:"输出内容"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// 数据集条目响应
type DatasetEntryResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:""`
	Data    DatasetEntryDTO `json:"data"`
}

// 数据集条目列表响应
type DatasetEntriesResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:""`
	Data    DatasetEntriesListData `json:"data"`
}

// 数据集条目列表数据
type DatasetEntriesListData struct {
	Entries []DatasetEntryDTO `json:"entries"`
	PagedData
}

// ========================= 笔记本相关 =========================

// 笔记本DTO
type NotebookDTO struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"笔记本名称"`
	Status    string    `json:"status" example:"running"`
	ProjectID uint      `json:"project_id" example:"1"`
	UserID    uint      `json:"user_id" example:"1"`
	Image     string    `json:"image" example:"pytorch:latest"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// 笔记本响应
type NotebookResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:""`
	Data    NotebookDTO `json:"data"`
}

// 笔记本列表响应
type NotebooksResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:""`
	Data    NotebooksListData `json:"data"`
}

// 笔记本列表数据
type NotebooksListData struct {
	Notebooks []NotebookDTO `json:"notebooks"`
	PagedData
}

// ========================= 训练任务相关 =========================

// 训练任务DTO
type TrainingJobDTO struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"训练任务名称"`
	Status    string    `json:"status" example:"running"`
	ProjectID uint      `json:"project_id" example:"1"`
	UserID    uint      `json:"user_id" example:"1"`
	Image     string    `json:"image" example:"pytorch:latest"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// 训练任务响应
type TrainingJobResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"Training Job created successfully"`
	Data    TrainingJobDTO `json:"data"`
}

// 训练任务列表响应
type TrainingJobsResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:""`
	Data    TrainingJobsListData `json:"data"`
}

// 训练任务列表数据
type TrainingJobsListData struct {
	TrainingJobs []TrainingJobDTO `json:"training_jobs"`
	PagedData
}

// ========================= 模型部署相关 =========================

// Triton部署DTO
type TritonDeployDTO struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"模型服务名称"`
	Status      string    `json:"status" example:"running"`
	ProjectID   uint      `json:"project_id" example:"1"`
	UserID      uint      `json:"user_id" example:"1"`
	ModelPath   string    `json:"model_path" example:"/models/resnet50"`
	ModelFormat string    `json:"model_format" example:"onnx"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// Triton部署响应
type TritonDeployResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:""`
	Data    TritonDeployDTO `json:"data"`
}

// Triton部署列表数据
type TritonDeployListData struct {
	Deploys []TritonDeployDTO `json:"deploys"`
	PagedData
}

// ========================= 第三方认证相关 =========================

// GitHub认证响应
type GitHubOAuthResponse struct {
	Success     bool        `json:"success" example:"true"`
	Message     string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	AccessToken string      `json:"access_token,omitempty"`
}

// 微信认证响应
type WechatLoginResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
