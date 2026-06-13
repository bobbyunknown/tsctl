package http

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"tsctl/internal/delivery/http/websocket"
	"tsctl/internal/domain"
	"tsctl/internal/usecase"
	"tsctl/pkg/config"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	tailscaleUC *usecase.TailscaleUseCase
	proxyUC     *usecase.ProxyUseCase
	wsHub       *websocket.Hub
}

func NewHandler(tailscaleUC *usecase.TailscaleUseCase, proxyUC *usecase.ProxyUseCase, wsHub *websocket.Hub) *Handler {
	return &Handler{
		tailscaleUC: tailscaleUC,
		proxyUC:     proxyUC,
		wsHub:       wsHub,
	}
}

type AuthStatusOutput struct {
	Body APIResponse
}

func (h *Handler) GetAuthStatus(ctx context.Context, input *struct{}) (*AuthStatusOutput, error) {
	status, err := h.tailscaleUC.GetAuthStatus()
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &AuthStatusOutput{
		Body: APIResponse{Success: true, Data: status},
	}, nil
}

type LogoutOutput struct {
	Body APIResponse
}

func (h *Handler) Logout(ctx context.Context, input *struct{}) (*LogoutOutput, error) {
	err := h.tailscaleUC.Logout()
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &LogoutOutput{
		Body: APIResponse{Success: true, Message: "logged out from tailnet - restart required to reconnect"},
	}, nil
}

type ServeInput struct {
	Body struct {
		Port       int  `json:"port" binding:"required"`
		Background bool `json:"background"`
	}
}

type ServeOutput struct {
	Body APIResponse
}

func (h *Handler) StartServe(ctx context.Context, input *ServeInput) (*ServeOutput, error) {
	output, err := h.tailscaleUC.StartServe(input.Body.Port, input.Body.Background)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}
	return &ServeOutput{
		Body: APIResponse{Success: true, Message: "serve started", Data: output},
	}, nil
}

type FunnelInput struct {
	Body struct {
		Port       int  `json:"port" binding:"required"`
		Background bool `json:"background"`
	}
}

type FunnelOutput struct {
	Body APIResponse
}

func (h *Handler) StartFunnel(ctx context.Context, input *FunnelInput) (*FunnelOutput, error) {
	output, err := h.tailscaleUC.StartFunnel(input.Body.Port, input.Body.Background)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}
	return &FunnelOutput{
		Body: APIResponse{Success: true, Message: "funnel started", Data: output},
	}, nil
}

type ServeStatusOutput struct {
	Body APIResponse
}

func (h *Handler) ServeStatus(ctx context.Context, input *struct{}) (*ServeStatusOutput, error) {
	output, err := h.tailscaleUC.GetServeStatus()
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &ServeStatusOutput{
		Body: APIResponse{Success: true, Data: output},
	}, nil
}

type FunnelStatusOutput struct {
	Body APIResponse
}

func (h *Handler) FunnelStatus(ctx context.Context, input *struct{}) (*FunnelStatusOutput, error) {
	output, err := h.tailscaleUC.GetFunnelStatus()
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &FunnelStatusOutput{
		Body: APIResponse{Success: true, Data: output},
	}, nil
}

type ResetServeOutput struct {
	Body APIResponse
}

func (h *Handler) ResetServe(ctx context.Context, input *struct{}) (*ResetServeOutput, error) {
	err := h.tailscaleUC.ResetServe()
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}
	return &ResetServeOutput{
		Body: APIResponse{Success: true, Message: "serve config reset"},
	}, nil
}

type StopServeInput struct {
	Port int `path:"port"`
}

type StopServeOutput struct {
	Body APIResponse
}

func (h *Handler) StopServe(ctx context.Context, input *StopServeInput) (*StopServeOutput, error) {
	err := h.tailscaleUC.StopServe(input.Port)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}
	return &StopServeOutput{
		Body: APIResponse{Success: true, Message: "serve stopped on port"},
	}, nil
}

type ResetFunnelOutput struct {
	Body APIResponse
}

func (h *Handler) ResetFunnel(ctx context.Context, input *struct{}) (*ResetFunnelOutput, error) {
	err := h.tailscaleUC.ResetFunnel()
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}
	return &ResetFunnelOutput{
		Body: APIResponse{Success: true, Message: "funnel config reset"},
	}, nil
}

type StopFunnelInput struct {
	Port int `path:"port"`
}

type StopFunnelOutput struct {
	Body APIResponse
}

func (h *Handler) StopFunnel(ctx context.Context, input *StopFunnelInput) (*StopFunnelOutput, error) {
	err := h.tailscaleUC.StopFunnel(input.Port)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	status, _ := h.tailscaleUC.GetServeStatus()
	event := map[string]interface{}{"type": "serve_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}
	return &StopFunnelOutput{
		Body: APIResponse{Success: true, Message: "funnel stopped on port"},
	}, nil
}

type EnableSSHOutput struct {
	Body APIResponse
}

func (h *Handler) EnableSSH(ctx context.Context, input *struct{}) (*EnableSSHOutput, error) {
	err := h.tailscaleUC.EnableSSH()
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &EnableSSHOutput{
		Body: APIResponse{Success: true, Message: "ssh enabled"},
	}, nil
}

type StatusOutput struct {
	Body APIResponse
}

func (h *Handler) Status(ctx context.Context, input *struct{}) (*StatusOutput, error) {
	output, err := h.tailscaleUC.GetStatus()
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &StatusOutput{
		Body: APIResponse{Success: true, Data: output},
	}, nil
}

type AppLogsInput struct {
	Lines int `query:"lines" default:"100"`
}

type AppLogsOutput struct {
	Body APIResponse
}

func (h *Handler) GetAppLogs(ctx context.Context, input *AppLogsInput) (*AppLogsOutput, error) {
	cfg := config.Get()
	logs, err := readLastLines(cfg.Logging.AppLogPath, fmt.Sprintf("%d", input.Lines))
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return &AppLogsOutput{
		Body: APIResponse{Success: true, Data: logs},
	}, nil
}

type ClearLogsOutput struct {
	Body APIResponse
}

func (h *Handler) ClearLogs(ctx context.Context, input *struct{}) (*ClearLogsOutput, error) {
	cfg := config.Get()
	logPath := cfg.Logging.AppLogPath
	err := os.Truncate(logPath, 0)
	if err != nil {
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to clear logs: %v", err))
	}
	return &ClearLogsOutput{
		Body: APIResponse{Success: true, Message: "Logs cleared successfully"},
	}, nil
}

type ProxyStartInput struct {
	Body struct {
		Mode         string `json:"mode" binding:"required"`
		Port         int    `json:"port,omitempty"`
		Ports        []int  `json:"ports,omitempty"`
		ExcludePorts []int  `json:"exclude_ports,omitempty"`
		ScanInterval int    `json:"scan_interval,omitempty"`
	}
}

type ProxyStartOutput struct {
	Body APIResponse
}

func (h *Handler) StartProxy(ctx context.Context, input *ProxyStartInput) (*ProxyStartOutput, error) {
	req := input.Body
	switch req.Mode {
	case "single":
		if req.Port <= 0 {
			return nil, huma.Error400BadRequest("port is required for single mode")
		}
		if err := h.proxyUC.StartPorts([]int{req.Port}); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
	case "multi":
		if len(req.Ports) == 0 {
			return nil, huma.Error400BadRequest("ports is required for multi mode")
		}
		if err := h.proxyUC.StartPorts(req.Ports); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
	case "all":
		interval := req.ScanInterval
		if interval <= 0 {
			interval = 5
		}
		if err := h.proxyUC.StartAutoScan(interval, req.ExcludePorts); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
	default:
		return nil, huma.Error400BadRequest("mode must be single, multi, or all")
	}
	status := h.proxyUC.GetStatus()
	event := map[string]interface{}{"type": "proxy_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}
	return &ProxyStartOutput{
		Body: APIResponse{Success: true, Message: fmt.Sprintf("proxy started (mode: %s)", req.Mode), Data: status},
	}, nil
}

type ProxyStopInput struct {
	Body struct {
		Ports []int `json:"ports" binding:"required"`
	}
}

type ProxyStopOutput struct {
	Body APIResponse
}

func (h *Handler) StopProxy(ctx context.Context, input *ProxyStopInput) (*ProxyStopOutput, error) {
	if err := h.proxyUC.StopPorts(input.Body.Ports); err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	status := h.proxyUC.GetStatus()
	event := map[string]interface{}{"type": "proxy_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}
	return &ProxyStopOutput{
		Body: APIResponse{Success: true, Message: "proxy stopped", Data: status},
	}, nil
}

type StopAllProxyOutput struct {
	Body APIResponse
}

func (h *Handler) StopAllProxy(ctx context.Context, input *struct{}) (*StopAllProxyOutput, error) {
	h.proxyUC.StopAutoScan()
	if err := h.proxyUC.StopAll(); err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	status := domain.ProxyStatusResponse{
		IsAutoScanActive: false,
		Proxies:          []domain.PortProxyInfo{},
	}
	event := map[string]interface{}{"type": "proxy_status_changed", "data": status}
	if data, err := json.Marshal(event); err == nil {
		h.wsHub.Broadcast(data)
	}
	return &StopAllProxyOutput{
		Body: APIResponse{Success: true, Message: "all proxies stopped"},
	}, nil
}

type ProxyStatusOutput struct {
	Body APIResponse
}

func (h *Handler) ProxyStatus(ctx context.Context, input *struct{}) (*ProxyStatusOutput, error) {
	status := h.proxyUC.GetStatus()
	return &ProxyStatusOutput{
		Body: APIResponse{Success: true, Data: status},
	}, nil
}

// ProxyAvatar needs raw gin.Context so stay as Gin handler

func (h *Handler) ProxyAvatar(c *gin.Context) {
	status, err := h.tailscaleUC.GetAuthStatus()
	if err != nil {
		c.JSON(500, APIResponse{Success: false, Message: err.Error()})
		return
	}

	avatarURL, ok := status["user_profile_pic"].(string)
	if !ok || avatarURL == "" {
		c.Status(404)
		return
	}

	resp, err := http.Get(avatarURL)
	if err != nil {
		c.Status(500)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		c.Status(resp.StatusCode)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		c.Header("Content-Type", contentType)
	}
	c.Header("Cache-Control", "public, max-age=3600")

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.Status(500)
		return
	}
}

func readLastLines(filePath string, linesCount string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return nil, err
	}

	count := 100
		_, _ = fmt.Sscanf(linesCount, "%d", &count)
	if len(lines) < count {
		return lines, nil
	}

	return lines[len(lines)-count:], nil
}
