package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// SetuP Huma API
	api := SetupHuma(router)

	// WebSocket & Avatar needs raw gin access
	v1 := router.Group("/api/v1")
	v1.GET("/ws", handler.HandleWebSocket)
	v1.GET("/avatar", handler.ProxyAvatar)

	// Huma-registered endpoints
	huma.Register(api, huma.Operation{
		OperationID: "get-auth-status",
		Method:      http.MethodGet,
		Path:        "/api/v1/auth/status",
		Tags:        []string{"auth"},
	}, handler.GetAuthStatus)

	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/logout",
		Tags:        []string{"auth"},
	}, handler.Logout)

	huma.Register(api, huma.Operation{
		OperationID: "start-serve",
		Method:      http.MethodPost,
		Path:        "/api/v1/serve",
		Tags:        []string{"tailscale"},
	}, handler.StartServe)

	huma.Register(api, huma.Operation{
		OperationID: "serve-status",
		Method:      http.MethodGet,
		Path:        "/api/v1/serve/status",
		Tags:        []string{"tailscale"},
	}, handler.ServeStatus)

	huma.Register(api, huma.Operation{
		OperationID: "reset-serve",
		Method:      http.MethodDelete,
		Path:        "/api/v1/serve",
		Tags:        []string{"tailscale"},
	}, handler.ResetServe)

	huma.Register(api, huma.Operation{
		OperationID: "stop-serve",
		Method:      http.MethodDelete,
		Path:        "/api/v1/serve/{port}",
		Tags:        []string{"tailscale"},
	}, handler.StopServe)

	huma.Register(api, huma.Operation{
		OperationID: "start-funnel",
		Method:      http.MethodPost,
		Path:        "/api/v1/funnel",
		Tags:        []string{"tailscale"},
	}, handler.StartFunnel)

	huma.Register(api, huma.Operation{
		OperationID: "funnel-status",
		Method:      http.MethodGet,
		Path:        "/api/v1/funnel/status",
		Tags:        []string{"tailscale"},
	}, handler.FunnelStatus)

	huma.Register(api, huma.Operation{
		OperationID: "reset-funnel",
		Method:      http.MethodDelete,
		Path:        "/api/v1/funnel",
		Tags:        []string{"tailscale"},
	}, handler.ResetFunnel)

	huma.Register(api, huma.Operation{
		OperationID: "stop-funnel",
		Method:      http.MethodDelete,
		Path:        "/api/v1/funnel/{port}",
		Tags:        []string{"tailscale"},
	}, handler.StopFunnel)

	huma.Register(api, huma.Operation{
		OperationID: "enable-ssh",
		Method:      http.MethodPost,
		Path:        "/api/v1/ssh/enable",
		Tags:        []string{"ssh"},
	}, handler.EnableSSH)

	huma.Register(api, huma.Operation{
		OperationID: "get-status",
		Method:      http.MethodGet,
		Path:        "/api/v1/status",
		Tags:        []string{"tailscale"},
	}, handler.Status)

	huma.Register(api, huma.Operation{
		OperationID: "start-proxy",
		Method:      http.MethodPost,
		Path:        "/api/v1/proxy/start",
		Tags:        []string{"proxy"},
	}, handler.StartProxy)

	huma.Register(api, huma.Operation{
		OperationID: "stop-proxy",
		Method:      http.MethodPost,
		Path:        "/api/v1/proxy/stop",
		Tags:        []string{"proxy"},
	}, handler.StopProxy)

	huma.Register(api, huma.Operation{
		OperationID: "stop-all-proxy",
		Method:      http.MethodDelete,
		Path:        "/api/v1/proxy",
		Tags:        []string{"proxy"},
	}, handler.StopAllProxy)

	huma.Register(api, huma.Operation{
		OperationID: "proxy-status",
		Method:      http.MethodGet,
		Path:        "/api/v1/proxy/status",
		Tags:        []string{"proxy"},
	}, handler.ProxyStatus)

	huma.Register(api, huma.Operation{
		OperationID: "get-logs",
		Method:      http.MethodGet,
		Path:        "/api/v1/logs/app",
		Tags:        []string{"logs"},
	}, handler.GetAppLogs)

	huma.Register(api, huma.Operation{
		OperationID: "clear-logs",
		Method:      http.MethodDelete,
		Path:        "/api/v1/logs/app",
		Tags:        []string{"logs"},
	}, handler.ClearLogs)

	return router
}
