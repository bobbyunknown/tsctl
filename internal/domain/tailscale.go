package domain

import (
	"context"

	"tailscale.com/ipn/ipnstate"
)

type TailscaleService interface {
	Serve(port int, background bool) (string, error)
	Funnel(port int, background bool) (string, error)
	ServeStatus() (string, error)
	FunnelStatus() (string, error)
	ServeReset() error
	FunnelReset() error
	EnableSSH() error
	Status() (string, error)
	GetFullStatus(ctx context.Context) (*ipnstate.Status, error)
	Logout(ctx context.Context) error
}
