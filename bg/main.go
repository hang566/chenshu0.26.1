package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"yunshu-bg/config"
	"yunshu-bg/database"
	"yunshu-bg/handlers"
	"yunshu-bg/middleware"
	"yunshu-bg/models"
	"yunshu-bg/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// 初始化数据库
	db, err := database.Init(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移数据库
	if err := db.AutoMigrate(&models.User{}, &models.UserData{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化服务器
	r := gin.Default()

	// 配置中间件
	r.Use(middleware.CORS())

	// API 路由
	api := r.Group("/api")
	{
		// 用户相关路由
		userHandler := handlers.NewUserHandler(db)
		userRoutes := api.Group("/user")
		{
			routes.UserRoutes(userRoutes, userHandler)
		}

		// 需要认证的路由
		authRoutes := api.Group("")
		authRoutes.Use(middleware.Auth())
		{
			userRoutes := authRoutes.Group("/user")
			{
				routes.AuthUserRoutes(userRoutes, userHandler)
			}

			// 数据存取路由
			dataHandler := handlers.NewDataHandler(db)
			dataRoutes := authRoutes.Group("/data")
			{
				routes.DataRoutes(dataRoutes, dataHandler)
			}
		}
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 服务器配置
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("服务器启动在 http://localhost:%d", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
