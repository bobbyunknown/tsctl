package usecase

import (
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

func (u *TailscaleUseCase) GetServeStatus() (string, error) {
	return u.service.ServeStatus()
}

func (u *TailscaleUseCase) GetFunnelStatus() (string, error) {
	return u.service.FunnelStatus()
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

func (u *TailscaleUseCase) GetStatus() (string, error) {
	return u.service.Status()
}
