package domain

import "time"

type DaemonManager interface {
	Start() error
	Stop() error
	Restart() error
	Status() (*DaemonStatus, error)
}

type DaemonStatus struct {
	Running bool          `json:"running"`
	PID     int           `json:"pid"`
	Uptime  time.Duration `json:"uptime"`
}
