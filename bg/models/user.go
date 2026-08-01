package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	EmailOrPhone string         `gorm:"size:100" json:"emailOrPhone"`
	Password     string         `gorm:"size:255;not null" json:"-"` // 不返回密码
	Code         string         `gorm:"uniqueIndex;size:50" json:"code"`
	Type         string         `gorm:"size:20;default:online" json:"type"` // online 联机, offline 脱机
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Status       int            `gorm:"default:1" json:"status"` // 1 正常, 0 禁用
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserResponse 用户信息响应（不含密码）
type UserResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	EmailOrPhone string    `json:"emailOrPhone"`
	Code         string    `json:"code"`
	Type         string    `json:"type"`
	Avatar       string    `json:"avatar"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username     string `json:"username" binding:"required,min=2,max=50"`
	EmailOrPhone string `json:"emailOrPhone" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	LoginID  string `json:"loginId" binding:"required"` // 用户名或编码
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username     string `json:"username" binding:"omitempty,min=2,max=50"`
	EmailOrPhone string `json:"emailOrPhone" binding:"omitempty"`
	Avatar       string `json:"avatar" binding:"omitempty"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Username:     u.Username,
		EmailOrPhone: u.EmailOrPhone,
		Code:         u.Code,
		Type:         u.Type,
		Avatar:       u.Avatar,
		Status:       u.Status,
		CreatedAt:    u.CreatedAt,
	}
}
