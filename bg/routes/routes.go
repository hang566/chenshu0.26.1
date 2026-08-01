package routes

import (
	"yunshu-bg/handlers"

	"github.com/gin-gonic/gin"
)

// UserRoutes 用户相关路由（无需认证）
func UserRoutes(rg *gin.RouterGroup, handler *handlers.UserHandler) {
	// 注册
	rg.POST("/register", handler.Register)

	// 登录
	rg.POST("/login", handler.Login)
}

// AuthUserRoutes 需要认证的用户路由
func AuthUserRoutes(rg *gin.RouterGroup, handler *handlers.UserHandler) {
	// 获取当前用户信息
	rg.GET("/info", handler.GetUserInfo)

	// 更新用户信息
	rg.PUT("/info", handler.UpdateUser)

	// 修改密码
	rg.PUT("/password", handler.ChangePassword)

	// 登出
	rg.POST("/logout", handler.Logout)
}

// DataRoutes 数据相关路由（需要认证）
func DataRoutes(rg *gin.RouterGroup, handler *handlers.DataHandler) {
	// 保存单条数据
	rg.POST("/save", handler.SaveData)

	// 批量保存数据
	rg.POST("/batch-save", handler.BatchSaveData)

	// 获取数据列表（支持 ?category=xxx&dataKey=xxx 筛选）
	rg.GET("/list", handler.GetData)

	// 获取所有分类
	rg.GET("/categories", handler.GetCategories)

	// 根据 ID 获取单条数据
	rg.GET("/:id", handler.GetDataByID)

	// 根据 ID 删除数据
	rg.DELETE("/:id", handler.DeleteData)

	// 按分类删除数据
	rg.DELETE("/category/:category", handler.DeleteDataByCategory)
}
