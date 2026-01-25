package http

import (
	"bufio"
	"io"
	"os"

	"tsctl/internal/usecase"
	"tsctl/pkg/config"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	daemonUC    *usecase.DaemonUseCase
	tailscaleUC *usecase.TailscaleUseCase
}

func NewHandler(daemonUC *usecase.DaemonUseCase, tailscaleUC *usecase.TailscaleUseCase) *Handler {
	return &Handler{
		daemonUC:    daemonUC,
		tailscaleUC: tailscaleUC,
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

// @Summary Start daemon
// @Description Start the tailscaled daemon process
// @Tags daemon
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/daemon/start [post]
func (h *Handler) StartDaemon(c *gin.Context) {
	err := h.daemonUC.Start()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Message: "daemon started"})
}

// @Summary Stop daemon
// @Description Stop the tailscaled daemon process
// @Tags daemon
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/daemon/stop [post]
func (h *Handler) StopDaemon(c *gin.Context) {
	err := h.daemonUC.Stop()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Message: "daemon stopped"})
}

// @Summary Restart daemon
// @Description Restart the tailscaled daemon process
// @Tags daemon
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/daemon/restart [post]
func (h *Handler) RestartDaemon(c *gin.Context) {
	err := h.daemonUC.Restart()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Message: "daemon restarted"})
}

// @Summary Get daemon status
// @Description Get the current status of tailscaled daemon
// @Tags daemon
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/daemon/status [get]
func (h *Handler) DaemonStatus(c *gin.Context) {
	status, err := h.daemonUC.Status()
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Data: status})
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

// @Summary Get daemon logs
// @Description Get tailscaled daemon log file content
// @Tags logs
// @Produce json
// @Param lines query int false "Number of lines to retrieve" default(100)
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/logs/daemon [get]
func (h *Handler) GetDaemonLogs(c *gin.Context) {
	lines := c.DefaultQuery("lines", "100")
	cfg := config.Get()
	logs, err := readLastLines(cfg.Logging.DaemonLogPath, lines)
	if err != nil {
		c.JSON(500, Response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, Response{Success: true, Data: logs})
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
