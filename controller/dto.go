package controller

// 响应结构体定义
type ProjectMembershipResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// 请求结构体定义


type RoleUpdateRequest struct {
	Role int `json:"role" binding:"required"`
}

// 成员DTO
type MemberDTO struct {
	ProjectID int    `json:"projectId"`
	UserID    int    `json:"userId"`
	Username  string `json:"username"`
	Role      int    `json:"role"`
}

type ProjectMembershipDTO struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"userId"`
	ProjectID uint    `json:"projectId"`
	Role      int     `json:"role"`
	User      UserDTO `json:"user"`
}

// UserDTO 结构体
type UserDTO struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        int    `json:"role"`
	Status      int    `json:"status"`
	Email       string `json:"email"`
}
