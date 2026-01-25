package usecase

import (
	"tsctl/internal/domain"
)

type DaemonUseCase struct {
	manager domain.DaemonManager
}

func NewDaemonUseCase(manager domain.DaemonManager) *DaemonUseCase {
	return &DaemonUseCase{manager: manager}
}

func (u *DaemonUseCase) Start() error {
	return u.manager.Start()
}

func (u *DaemonUseCase) Stop() error {
	return u.manager.Stop()
}

func (u *DaemonUseCase) Restart() error {
	return u.manager.Restart()
}

func (u *DaemonUseCase) Status() (*domain.DaemonStatus, error) {
	return u.manager.Status()
}
