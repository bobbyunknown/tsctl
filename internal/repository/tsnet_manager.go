package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"tsctl/pkg/config"
	"tsctl/pkg/logger"

	"tailscale.com/ipn"
	"tailscale.com/ipn/ipnstate"
	"tailscale.com/tsnet"
)

type TsnetManager struct {
	server *tsnet.Server
	mu     sync.RWMutex
	status *ipnstate.Status
}

func NewTsnetManager(cfg *config.Config) *TsnetManager {
	return &TsnetManager{
		server: &tsnet.Server{
			Dir:       cfg.Tailscale.StateDir,
			Hostname:  cfg.Tailscale.Hostname,
			AuthKey:   cfg.Tailscale.AuthKey,
			Ephemeral: cfg.Tailscale.Ephemeral,
			Logf:      logger.Log.Printf,
		},
	}
}

func (t *TsnetManager) GetServer() *tsnet.Server {
	return t.server
}

func (t *TsnetManager) Start(ctx context.Context) error {
	logger.Log.Info("starting embedded tailscale daemon")

	status, err := t.server.Up(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tsnet: %w", err)
	}

	t.mu.Lock()
	t.status = status
	t.mu.Unlock()

	logger.Log.WithFields(map[string]interface{}{
		"hostname": status.Self.HostName,
		"ips":      status.Self.TailscaleIPs,
	}).Info("embedded daemon started successfully")

	go func() {
		time.Sleep(2 * time.Second)
		t.RestoreServeConfig(ctx)
	}()

	return nil
}

func (t *TsnetManager) RestoreServeConfig(ctx context.Context) {
	lc, err := t.server.LocalClient()
	if err != nil {
		logger.Log.WithError(err).Warn("failed to get local client for restore")
		return
	}

	existingConfig, err := lc.GetServeConfig(ctx)
	if err != nil {
		logger.Log.WithError(err).Debug("no existing serve config to restore")
		return
	}

	if existingConfig != nil && (existingConfig.TCP != nil || existingConfig.Web != nil) {
		logger.Log.Info("restored serve config from state store")
	}
}

func (t *TsnetManager) Stop() error {
	logger.Log.Info("stopping embedded daemon")
	return t.server.Close()
}

func (t *TsnetManager) Status() (string, error) {
	ctx := context.Background()
	status, err := t.getStatus(ctx)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(status)
	return string(data), err
}

func (t *TsnetManager) GetFullStatus(ctx context.Context) (*ipnstate.Status, error) {
	lc, err := t.server.LocalClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get local client: %w", err)
	}
	return lc.Status(ctx)
}

func (t *TsnetManager) getStatus(ctx context.Context) (*ipnstate.Status, error) {
	lc, err := t.server.LocalClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get local client: %w", err)
	}
	return lc.Status(ctx)
}

func (t *TsnetManager) Serve(port int, background bool) (string, error) {
	ctx := context.Background()
	config := &ipn.ServeConfig{
		TCP: map[uint16]*ipn.TCPPortHandler{
			uint16(port): {HTTPS: true},
		},
	}

	if err := t.setServeConfig(ctx, config); err != nil {
		return "", err
	}

	logger.Log.WithField("port", port).Info("serve started")
	return fmt.Sprintf("serve started on port %d", port), nil
}

func (t *TsnetManager) Funnel(port int, background bool) (string, error) {
	ctx := context.Background()

	status, err := t.getStatus(ctx)
	if err != nil {
		return "", err
	}

	dnsName := strings.TrimSuffix(status.Self.DNSName, ".")

	config := &ipn.ServeConfig{
		TCP: map[uint16]*ipn.TCPPortHandler{
			443: {HTTPS: true},
		},
		Web: map[ipn.HostPort]*ipn.WebServerConfig{
			ipn.HostPort(fmt.Sprintf("%s:443", dnsName)): {
				Handlers: map[string]*ipn.HTTPHandler{
					"/": {
						Proxy: fmt.Sprintf("http://127.0.0.1:%d", port),
					},
				},
			},
		},
		AllowFunnel: map[ipn.HostPort]bool{
			ipn.HostPort(fmt.Sprintf("%s:443", dnsName)): true,
		},
	}

	if err := t.setServeConfig(ctx, config); err != nil {
		return "", err
	}

	logger.Log.WithField("port", port).Info("funnel started")
	return fmt.Sprintf("funnel started on port %d", port), nil
}

func (t *TsnetManager) ServeStatus() (string, error) {
	ctx := context.Background()
	config, err := t.getServeConfig(ctx)
	if err != nil {
		return "", err
	}

	status, err := t.getStatus(ctx)
	if err != nil {
		return "", err
	}

	type Service struct {
		Port      uint16 `json:"port"`
		Type      string `json:"type"`
		LocalURL  string `json:"local_url"`
		PublicURL string `json:"public_url,omitempty"`
	}

	response := struct {
		Services []Service `json:"services"`
	}{
		Services: []Service{},
	}

	dnsName := strings.TrimSuffix(status.Self.DNSName, ".")
	if config != nil && config.TCP != nil {
		for port := range config.TCP {
			svc := Service{
				Port:     port,
				Type:     "serve",
				LocalURL: fmt.Sprintf("https://%s:%d", dnsName, port),
			}

			if config.AllowFunnel != nil {
				portKey443 := ipn.HostPort(fmt.Sprintf("%s:443", dnsName))
				if config.AllowFunnel[portKey443] && port == 443 {
					svc.Type = "funnel"
					svc.PublicURL = fmt.Sprintf("https://%s/", dnsName)

					if config.Web != nil {
						if webCfg, ok := config.Web[portKey443]; ok {
							if handler, ok := webCfg.Handlers["/"]; ok && handler.Proxy != "" {
								svc.LocalURL = handler.Proxy
							}
						}
					}
				}
			}

			response.Services = append(response.Services, svc)
		}
	}

	data, err := json.Marshal(response)
	return string(data), err
}

func (t *TsnetManager) FunnelStatus() (string, error) {
	return t.ServeStatus()
}

func (t *TsnetManager) ServeReset() error {
	ctx := context.Background()
	logger.Log.Info("resetting serve config")
	return t.setServeConfig(ctx, new(ipn.ServeConfig))
}

func (t *TsnetManager) FunnelReset() error {
	return t.ServeReset()
}

func (t *TsnetManager) ServeStop(port int) error {
	ctx := context.Background()
	config, err := t.getServeConfig(ctx)
	if err != nil {
		return err
	}
	if config == nil || config.TCP == nil {
		return nil
	}

	port16 := uint16(port)
	if _, exists := config.TCP[port16]; exists {
		delete(config.TCP, port16)
		logger.Log.WithField("port", port).Info("stopped serve on port")
		return t.setServeConfig(ctx, config)
	}
	return nil
}

func (t *TsnetManager) FunnelStop(port int) error {
	ctx := context.Background()
	config, err := t.getServeConfig(ctx)
	if err != nil {
		return err
	}
	if config == nil || config.Web == nil {
		return nil
	}

	status, err := t.getStatus(ctx)
	if err != nil {
		return err
	}
	dnsName := strings.TrimSuffix(status.Self.DNSName, ".")
	portKey443 := ipn.HostPort(fmt.Sprintf("%s:443", dnsName))

	if webCfg, exists := config.Web[portKey443]; exists {
		// remove handler for this port
		targetProxy := fmt.Sprintf("http://127.0.0.1:%d", port)
		modified := false
		for path, handler := range webCfg.Handlers {
			if handler.Proxy == targetProxy {
				delete(webCfg.Handlers, path)
				modified = true
			}
		}

		if modified {
			if len(webCfg.Handlers) == 0 {
				delete(config.Web, portKey443)
				if config.AllowFunnel != nil {
					delete(config.AllowFunnel, portKey443)
				}
			}
			logger.Log.WithField("port", port).Info("stopped funnel on port")
			return t.setServeConfig(ctx, config)
		}
	}
	return nil
}

func (t *TsnetManager) EnableSSH() error {
	ctx := context.Background()
	lc, err := t.server.LocalClient()
	if err != nil {
		return fmt.Errorf("failed to get local client: %w", err)
	}

	logger.Log.Info("enabling SSH")
	_, err = lc.EditPrefs(ctx, &ipn.MaskedPrefs{
		Prefs:     ipn.Prefs{RunSSH: true},
		RunSSHSet: true,
	})
	return err
}

func (t *TsnetManager) setServeConfig(ctx context.Context, cfg *ipn.ServeConfig) error {
	lc, err := t.server.LocalClient()
	if err != nil {
		return fmt.Errorf("failed to get local client: %w", err)
	}

	logger.Log.WithField("config", cfg).Debug("setting serve config")
	return lc.SetServeConfig(ctx, cfg)
}

func (t *TsnetManager) getServeConfig(ctx context.Context) (*ipn.ServeConfig, error) {
	lc, err := t.server.LocalClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get local client: %w", err)
	}
	return lc.GetServeConfig(ctx)
}

func (t *TsnetManager) Logout(ctx context.Context) error {
	lc, err := t.server.LocalClient()
	if err != nil {
		return fmt.Errorf("failed to get local client: %w", err)
	}

	logger.Log.Info("logging out from tailnet")
	return lc.Logout(ctx)
}
