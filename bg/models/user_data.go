package models

import (
	"time"

	"gorm.io/gorm"
)

// UserData 用户数据模型（按模块分类存储 JSON 数据）
type UserData struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"userId"`   // 关联用户
	Category  string         `gorm:"size:50;index;not null" json:"category"` // 数据分类（order, inventory, attendance 等）
	DataKey   string         `gorm:"size:100;index" json:"dataKey"`  // 数据键（可选，用于细分类）
	Content   string         `gorm:"type:text" json:"content"`       // JSON 格式数据内容
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SaveDataRequest 保存数据请求
type SaveDataRequest struct {
	Category string `json:"category" binding:"required"` // 数据分类
	DataKey  string `json:"dataKey"`                     // 数据键（可选）
	Content  string `json:"content" binding:"required"`  // JSON 字符串
}

// BatchSaveDataRequest 批量保存数据请求
type BatchSaveDataRequest struct {
	Items []SaveDataRequest `json:"items" binding:"required,min=1"`
}

// QueryDataRequest 查询数据请求
type QueryDataRequest struct {
	Category string `json:"category"` // 按分类筛选
	DataKey  string `json:"dataKey"`  // 按键筛选
}
