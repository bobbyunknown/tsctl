package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"tsctl/internal/domain"
)

type TailscaleUseCase struct {
	service domain.TailscaleService
}

func NewTailscaleUseCase(service domain.TailscaleService) *TailscaleUseCase {
	return &TailscaleUseCase{service: service}
}

func (u *TailscaleUseCase) StartServe(port int, background bool) (string, error) {
	if port <= 0 || port > 65535 {
		return "", fmt.Errorf("invalid port number")
	}
	return u.service.Serve(port, background)
}

func (u *TailscaleUseCase) StartFunnel(port int, background bool) (string, error) {
	if port <= 0 || port > 65535 {
		return "", fmt.Errorf("invalid port number")
	}
	return u.service.Funnel(port, background)
}

func (u *TailscaleUseCase) GetServeStatus() (*ServeStatusResponse, error) {
	jsonStr, err := u.service.ServeStatus()
	if err != nil {
		return nil, err
	}

	var status ServeStatusResponse
	if err := json.Unmarshal([]byte(jsonStr), &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (u *TailscaleUseCase) GetFunnelStatus() (*ServeStatusResponse, error) {
	jsonStr, err := u.service.FunnelStatus()
	if err != nil {
		return nil, err
	}

	var status ServeStatusResponse
	if err := json.Unmarshal([]byte(jsonStr), &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (u *TailscaleUseCase) ResetServe() error {
	return u.service.ServeReset()
}

func (u *TailscaleUseCase) ResetFunnel() error {
	return u.service.FunnelReset()
}

func (u *TailscaleUseCase) EnableSSH() error {
	return u.service.EnableSSH()
}

func (u *TailscaleUseCase) GetStatus() (interface{}, error) {
	jsonStr, err := u.service.Status()
	if err != nil {
		return nil, err
	}

	var status map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &status); err != nil {
		return nil, err
	}
	return status, nil
}

func (u *TailscaleUseCase) GetAuthStatus() (map[string]interface{}, error) {
	ctx := context.Background()
	status, err := u.service.GetFullStatus(ctx)
	if err != nil {
		return nil, err
	}

	authenticated := status.BackendState == "Running"

	result := map[string]interface{}{
		"authenticated": authenticated,
		"backend_state": status.BackendState,
	}

	if !authenticated && status.AuthURL != "" {
		result["auth_url"] = status.AuthURL
	}

	if status.Self != nil {
		result["node_key"] = status.Self.PublicKey.String()
		result["hostname"] = status.Self.HostName
	}

	if len(status.TailscaleIPs) > 0 {
		ips := make([]string, len(status.TailscaleIPs))
		for i, ip := range status.TailscaleIPs {
			ips[i] = ip.String()
		}
		result["ips"] = ips
	}

	if status.CurrentTailnet != nil {
		result["tailnet_name"] = status.CurrentTailnet.Name
		result["tailnet_dns_suffix"] = status.CurrentTailnet.MagicDNSSuffix
	}

	if len(status.User) > 0 {
		for _, user := range status.User {
			result["user_display_name"] = user.DisplayName
			result["user_email"] = user.LoginName
			result["user_profile_pic"] = user.ProfilePicURL
			break
		}
	}

	if status.Self != nil && len(status.Self.CapMap) > 0 {
		result["is_admin"] = false
		result["is_owner"] = false

		_, isAdmin := status.Self.CapMap["https://tailscale.com/cap/is-admin"]
		_, isOwner := status.Self.CapMap["https://tailscale.com/cap/is-owner"]

		result["is_admin"] = isAdmin
		result["is_owner"] = isOwner
	}

	if status.Self != nil {
		if !status.Self.Created.IsZero() {
			result["created_at"] = status.Self.Created.Format("2006-01-02T15:04:05Z07:00")
		}
		if status.Self.KeyExpiry != nil && !status.Self.KeyExpiry.IsZero() {
			result["key_expiry"] = status.Self.KeyExpiry.Format("2006-01-02T15:04:05Z07:00")
		}
	}

	if status.Peer != nil {
		result["peer_count"] = len(status.Peer)
	}

	return result, nil
}

func (u *TailscaleUseCase) Logout() error {
	ctx := context.Background()
	return u.service.Logout(ctx)
}
