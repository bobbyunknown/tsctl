package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(handler *Handler, mode string) *gin.Engine {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		api.GET("/ws", handler.HandleWebSocket)
		api.GET("/avatar", handler.ProxyAvatar)

		auth := api.Group("/auth")
		{
			auth.GET("/status", handler.GetAuthStatus)
			auth.POST("/logout", handler.Logout)
		}

		api.POST("/serve", handler.StartServe)
		api.GET("/serve/status", handler.ServeStatus)
		api.DELETE("/serve", handler.ResetServe)

		api.POST("/funnel", handler.StartFunnel)
		api.GET("/funnel/status", handler.FunnelStatus)
		api.DELETE("/funnel", handler.ResetFunnel)

		ssh := api.Group("/ssh")
		{
			ssh.POST("/enable", handler.EnableSSH)
		}

		api.GET("/status", handler.Status)

		logs := api.Group("/logs")
		{
			logs.GET("/app", handler.GetAppLogs)
			logs.DELETE("/app", handler.ClearLogs)
		}
	}

	return router
}
