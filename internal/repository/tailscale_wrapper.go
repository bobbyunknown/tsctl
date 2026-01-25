package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"

	"tsctl/internal/domain"
	"tsctl/pkg/config"
	"tsctl/pkg/logger"
)

type tailscaleWrapper struct {
	binaryPath string
	socketPath string
}

func NewTailscaleWrapper() domain.TailscaleService {
	cfg := config.Get()
	return &tailscaleWrapper{
		binaryPath: cfg.Tailscale.BinaryPath,
		socketPath: cfg.Tailscale.SocketPath,
	}
}

func (t *tailscaleWrapper) runCommand(args ...string) (string, error) {
	cmd := exec.Command(t.binaryPath, args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("TS_SOCKET=%s", t.socketPath))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger.Log.WithField("args", args).Debug("executing tailscale command")

	if err := cmd.Run(); err != nil {
		logger.Log.WithError(err).WithField("stderr", stderr.String()).Error("tailscale command failed")
		return "", fmt.Errorf("%s: %s", err.Error(), stderr.String())
	}

	return stdout.String(), nil
}

func (t *tailscaleWrapper) Serve(port int, background bool) (string, error) {
	args := []string{"serve"}
	if background {
		args = append(args, "--bg")
	}
	args = append(args, strconv.Itoa(port))

	output, err := t.runCommand(args...)
	if err != nil {
		return "", err
	}

	logger.Log.WithField("port", port).Info("serve started")
	return output, nil
}

func (t *tailscaleWrapper) Funnel(port int, background bool) (string, error) {
	args := []string{"funnel"}
	if background {
		args = append(args, "--bg")
	}
	args = append(args, strconv.Itoa(port))

	output, err := t.runCommand(args...)
	if err != nil {
		return "", err
	}

	logger.Log.WithField("port", port).Info("funnel started")
	return output, nil
}

func (t *tailscaleWrapper) ServeStatus() (string, error) {
	return t.runCommand("serve", "status")
}

func (t *tailscaleWrapper) FunnelStatus() (string, error) {
	return t.runCommand("funnel", "status")
}

func (t *tailscaleWrapper) ServeReset() error {
	_, err := t.runCommand("serve", "reset")
	if err != nil {
		return err
	}
	logger.Log.Info("serve config reset")
	return nil
}

func (t *tailscaleWrapper) FunnelReset() error {
	_, err := t.runCommand("funnel", "reset")
	if err != nil {
		return err
	}
	logger.Log.Info("funnel config reset")
	return nil
}

func (t *tailscaleWrapper) EnableSSH() error {
	_, err := t.runCommand("up", "--ssh")
	if err != nil {
		return err
	}
	logger.Log.Info("ssh enabled")
	return nil
}

func (t *tailscaleWrapper) Status() (string, error) {
	return t.runCommand("status")
}
