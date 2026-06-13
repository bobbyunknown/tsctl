package usecase

import (
	"fmt"

	"tsctl/internal/domain"
)

type ProxyUseCase struct {
	proxy domain.PortProxyManager
}

func NewProxyUseCase(proxy domain.PortProxyManager) *ProxyUseCase {
	return &ProxyUseCase{proxy: proxy}
}

func (u *ProxyUseCase) StartPorts(ports []int) error {
	for _, p := range ports {
		if p <= 0 || p > 65535 {
			return fmt.Errorf("invalid port: %d", p)
		}
	}
	return u.proxy.StartProxy(ports)
}

func (u *ProxyUseCase) StopPorts(ports []int) error {
	return u.proxy.StopProxy(ports)
}

func (u *ProxyUseCase) StopAll() error {
	return u.proxy.StopAll()
}

func (u *ProxyUseCase) GetStatus() []domain.PortProxyInfo {
	return u.proxy.ActivePorts()
}

func (u *ProxyUseCase) StartAutoScan(interval int, excludePorts []int) error {
	if interval <= 0 {
		interval = 5
	}
	return u.proxy.StartAutoScan(interval, excludePorts)
}

func (u *ProxyUseCase) StopAutoScan() {
	u.proxy.StopAutoScan()
}
