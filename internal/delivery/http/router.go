package http

import (
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

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		daemon := api.Group("/daemon")
		{
			daemon.POST("/start", handler.StartDaemon)
			daemon.POST("/stop", handler.StopDaemon)
			daemon.POST("/restart", handler.RestartDaemon)
			daemon.GET("/status", handler.DaemonStatus)
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
			logs.GET("/daemon", handler.GetDaemonLogs)
		}
	}

	return router
}
