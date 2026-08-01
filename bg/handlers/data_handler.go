package handlers

import (
	"net/http"

	"yunshu-bg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DataHandler 数据处理器
type DataHandler struct {
	DB *gorm.DB
}

// NewDataHandler 创建数据处理器
func NewDataHandler(db *gorm.DB) *DataHandler {
	return &DataHandler{DB: db}
}

// SaveData 保存单条数据（存在则更新）
func (h *DataHandler) SaveData(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.SaveDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "detail": err.Error()})
		return
	}

	// 查找是否已存在同分类同键的数据
	var existing models.UserData
	result := h.DB.Where("user_id = ? AND category = ? AND data_key = ?", userID, req.Category, req.DataKey).First(&existing)

	if result.Error == nil {
		// 更新
		existing.Content = req.Content
		if err := h.DB.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新数据失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": existing})
	} else {
		// 新建
		data := models.UserData{
			UserID:   userID,
			Category: req.Category,
			DataKey:  req.DataKey,
			Content:  req.Content,
		}
		if err := h.DB.Create(&data).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存数据失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功", "data": data})
	}
}

// BatchSaveData 批量保存数据
func (h *DataHandler) BatchSaveData(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.BatchSaveDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "detail": err.Error()})
		return
	}

	tx := h.DB.Begin()
	successCount := 0
	for _, item := range req.Items {
		var existing models.UserData
		result := tx.Where("user_id = ? AND category = ? AND data_key = ?", userID, item.Category, item.DataKey).First(&existing)
		if result.Error == nil {
			existing.Content = item.Content
			if err := tx.Save(&existing).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批量保存失败"})
				return
			}
		} else {
			data := models.UserData{
				UserID:   userID,
				Category: item.Category,
				DataKey:  item.DataKey,
				Content:  item.Content,
			}
			if err := tx.Create(&data).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批量保存失败"})
				return
			}
		}
		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "批量保存成功", "data": gin.H{"count": successCount}})
}

// GetData 获取数据列表（支持按分类筛选）
func (h *DataHandler) GetData(c *gin.Context) {
	userID := c.GetUint("userID")

	category := c.Query("category")
	dataKey := c.Query("dataKey")

	query := h.DB.Where("user_id = ?", userID)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if dataKey != "" {
		query = query.Where("data_key = ?", dataKey)
	}

	var dataList []models.UserData
	if err := query.Find(&dataList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": dataList})
}

// GetDataByID 根据 ID 获取单条数据
func (h *DataHandler) GetDataByID(c *gin.Context) {
	userID := c.GetUint("userID")
	dataID := c.Param("id")

	var data models.UserData
	if err := h.DB.Where("id = ? AND user_id = ?", dataID, userID).First(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "数据不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": data})
}

// DeleteData 删除数据
func (h *DataHandler) DeleteData(c *gin.Context) {
	userID := c.GetUint("userID")
	dataID := c.Param("id")

	result := h.DB.Where("id = ? AND user_id = ?", dataID, userID).Delete(&models.UserData{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "数据不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// DeleteDataByCategory 按分类删除数据
func (h *DataHandler) DeleteDataByCategory(c *gin.Context) {
	userID := c.GetUint("userID")
	category := c.Param("category")

	result := h.DB.Where("user_id = ? AND category = ?", userID, category).Delete(&models.UserData{})
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功", "data": gin.H{"count": result.RowsAffected}})
}

// GetCategories 获取用户所有数据分类
func (h *DataHandler) GetCategories(c *gin.Context) {
	userID := c.GetUint("userID")

	var categories []string
	h.DB.Model(&models.UserData{}).Where("user_id = ?", userID).Distinct("category").Pluck("category", &categories)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": categories})
}
