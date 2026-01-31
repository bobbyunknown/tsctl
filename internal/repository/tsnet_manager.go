package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

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

	return nil
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
	config := &ipn.ServeConfig{
		TCP: map[uint16]*ipn.TCPPortHandler{
			uint16(port): {HTTPS: true},
		},
		AllowFunnel: map[ipn.HostPort]bool{
			ipn.HostPort(fmt.Sprintf(":%d", port)): true,
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

	dnsName := status.Self.DNSName
	if config != nil && config.TCP != nil {
		for port := range config.TCP {
			svc := Service{
				Port:     port,
				Type:     "serve",
				LocalURL: fmt.Sprintf("https://%s:%d", dnsName, port),
			}

			if config.AllowFunnel != nil {
				portKey := ipn.HostPort(fmt.Sprintf(":%d", port))
				if config.AllowFunnel[portKey] {
					svc.Type = "funnel"
					svc.PublicURL = fmt.Sprintf("https://%s:%d", dnsName, port)
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
