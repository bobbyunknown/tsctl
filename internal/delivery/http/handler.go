package http

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"tsctl/internal/delivery/http/websocket"
	"tsctl/internal/usecase"
	"tsctl/pkg/config"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tailscaleUC *usecase.TailscaleUseCase
	wsHub       *websocket.Hub
}

func NewHandler(tailscaleUC *usecase.TailscaleUseCase, wsHub *websocket.Hub) *Handler {
	return &Handler{
		tailscaleUC: tailscaleUC,
		wsHub:       wsHub,
	}
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ServeRequest struct {
	Port       int  `json:"port" binding:"required"`
	Background bool `json:"background"`
}

type FunnelRequest struct {
	Port       int  `json:"port" binding:"required"`
	Background bool `json:"background"`
}

func (h *Handler) GetAuthStatus(c *gin.Context) {
	status, err := h.tailscaleUC.GetAuthStatus()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    status,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	err := h.tailscaleUC.Logout()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(200, Response{
		Success: true,
		Message: "logged out from tailnet - restart required to reconnect",
	})
}

// @Summary Start serve
// @Description Start tailscale serve on specified port
// @Tags tailscale
// @Accept json
// @Produce json
// @Param body body ServeRequest true "Serve configuration"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/serve [post]
func (h *Handler) StartServe(c *gin.Context) {
	var req ServeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, Response{Success: false, Message: err.Error()})
		return
	}

	output, err := h.tailscaleUC.StartServe(req.Port, req.Background)
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}

	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}

	c.JSON(200, Response{Success: true, Message: "serve started", Data: output})
}

// @Summary Start funnel
// @Description Start tailscale funnel on specified port
// @Tags tailscale
// @Accept json
// @Produce json
// @Param body body FunnelRequest true "Funnel configuration"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/funnel [post]
func (h *Handler) StartFunnel(c *gin.Context) {
	var req FunnelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, Response{Success: false, Message: err.Error()})
		return
	}

	output, err := h.tailscaleUC.StartFunnel(req.Port, req.Background)
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}

	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}

	c.JSON(200, Response{Success: true, Message: "funnel started", Data: output})
}

// @Summary Get serve status
// @Description Get current tailscale serve configuration
// @Tags tailscale
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/serve/status [get]
func (h *Handler) ServeStatus(c *gin.Context) {
	output, err := h.tailscaleUC.GetServeStatus()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Data: output})
}

// @Summary Get funnel status
// @Description Get current tailscale funnel configuration
// @Tags tailscale
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/funnel/status [get]
func (h *Handler) FunnelStatus(c *gin.Context) {
	output, err := h.tailscaleUC.GetFunnelStatus()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Data: output})
}

// @Summary Reset serve
// @Description Reset tailscale serve configuration
// @Tags tailscale
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/serve [delete]
func (h *Handler) ResetServe(c *gin.Context) {
	err := h.tailscaleUC.ResetServe()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}

	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}

	c.JSON(200, Response{Success: true, Message: "serve config reset"})
}

// @Summary Reset funnel
// @Description Reset tailscale funnel configuration
// @Tags tailscale
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/funnel [delete]
func (h *Handler) ResetFunnel(c *gin.Context) {
	err := h.tailscaleUC.ResetFunnel()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}

	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}

	c.JSON(200, Response{Success: true, Message: "funnel config reset"})
}

// @Summary Enable SSH
// @Description Enable tailscale SSH access
// @Tags tailscale
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/ssh/enable [post]
func (h *Handler) EnableSSH(c *gin.Context) {
	err := h.tailscaleUC.EnableSSH()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Message: "ssh enabled"})
}

// @Summary Get tailscale status
// @Description Get current tailscale connection status
// @Tags tailscale
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/status [get]
func (h *Handler) Status(c *gin.Context) {
	output, err := h.tailscaleUC.GetStatus()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Data: output})
}

// @Summary Get application logs
// @Description Get application log file content
// @Tags logs
// @Produce json
// @Param lines query int false "Number of lines to retrieve" default(100)
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/logs/app [get]
func (h *Handler) GetAppLogs(c *gin.Context) {
	lines := c.DefaultQuery("lines", "100")
	cfg := config.Get()
	logs, err := readLastLines(cfg.Logging.AppLogPath, lines)
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Data: logs})
}

func (h *Handler) ClearLogs(c *gin.Context) {
	config := config.Get()
	logPath := config.Logging.AppLogPath

	err := os.Truncate(logPath, 0)
	if err != nil {
		c.JSON(500, Response{Success: false, Message: fmt.Sprintf("Failed to clear logs: %v", err)})
		return
	}

	c.JSON(200, Response{Success: true, Message: "Logs cleared successfully"})
}

func readLastLines(filePath string, linesCount string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return nil, err
	}

	count := 100
	if len(lines) < count {
		return lines, nil
	}

	return lines[len(lines)-count:], nil
}

// ProxyAvatar proxies avatar images from external sources (Google, GitHub, etc)
// to avoid CORS issues in the browser
func (h *Handler) ProxyAvatar(c *gin.Context) {
	status, err := h.tailscaleUC.GetAuthStatus()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}

	avatarURL, ok := status["user_profile_pic"].(string)
	if !ok || avatarURL == "" {
		c.Status(404)
		return
	}

	// Fetch the avatar from external URL
	resp, err := http.Get(avatarURL)
	if err != nil {
		c.Status(500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.Status(resp.StatusCode)
		return
	}

	// Set appropriate headers
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		c.Header("Content-Type", contentType)
	}
	c.Header("Cache-Control", "public, max-age=3600")

	// Copy the image data
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.Status(500)
		return
	}
}
