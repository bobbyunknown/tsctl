package repository

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"tsctl/internal/domain"
	"tsctl/pkg/config"
	"tsctl/pkg/logger"
)

type daemonManager struct {
	cmd       *exec.Cmd
	mu        sync.Mutex
	startTime time.Time
	logWriter io.Writer
}

func NewDaemonManager() domain.DaemonManager {
	cfg := config.Get()
	logWriter, _ := logger.GetDaemonLogger(cfg.Logging.DaemonLogPath)
	return &daemonManager{
		logWriter: logWriter,
	}
}

func (d *daemonManager) Start() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cmd != nil && d.cmd.Process != nil {
		if err := d.cmd.Process.Signal(syscall.Signal(0)); err == nil {
			return fmt.Errorf("daemon already running")
		}
	}

	cfg := config.Get()
	d.cmd = exec.Command(cfg.Tailscale.DaemonPath)
	d.cmd.Stdout = d.logWriter
	d.cmd.Stderr = d.logWriter

	if err := d.cmd.Start(); err != nil {
		logger.Log.WithError(err).Error("failed to start daemon")
		return err
	}

	d.startTime = time.Now()
	logger.Log.WithField("pid", d.cmd.Process.Pid).Info("daemon started")

	if err := config.SetAutoStart(true); err != nil {
		logger.Log.WithError(err).Warn("failed to save auto_start state")
	}

	return nil
}

func (d *daemonManager) Stop() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cmd == nil || d.cmd.Process == nil {
		return fmt.Errorf("daemon not running")
	}

	if err := d.cmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- d.cmd.Wait()
	}()

	select {
	case <-time.After(10 * time.Second):
		if err := d.cmd.Process.Kill(); err != nil {
			logger.Log.WithError(err).Error("failed to kill daemon")
			return err
		}
		logger.Log.Warn("daemon killed after timeout")
	case err := <-done:
		if err != nil {
			logger.Log.WithError(err).Debug("daemon process exited")
		}
	}

	logger.Log.Info("daemon stopped")
	d.cmd = nil

	if err := config.SetAutoStart(false); err != nil {
		logger.Log.WithError(err).Warn("failed to save auto_start state")
	}

	return nil
}

func (d *daemonManager) Restart() error {
	if err := d.Stop(); err != nil {
		logger.Log.WithError(err).Debug("stop error during restart")
	}
	time.Sleep(1 * time.Second)
	return d.Start()
}

func (d *daemonManager) Status() (*domain.DaemonStatus, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	status := &domain.DaemonStatus{
		Running: false,
		PID:     0,
		Uptime:  0,
	}

	if d.cmd != nil && d.cmd.Process != nil {
		if err := d.cmd.Process.Signal(syscall.Signal(0)); err == nil {
			status.Running = true
			status.PID = d.cmd.Process.Pid
			status.Uptime = time.Since(d.startTime)
		}
	}

	return status, nil
}
