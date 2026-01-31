package repository

import (
	"context"
	"os"
	"testing"

	"tsctl/pkg/config"
	"tsctl/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.Init(os.DevNull, "error", "json")
	os.Exit(m.Run())
}

func TestNewTsnetManager(t *testing.T) {
	cfg := &config.Config{
		Tailscale: config.TailscaleConfig{
			StateDir:  t.TempDir(),
			Hostname:  "test-node",
			AuthKey:   "",
			Ephemeral: true,
		},
	}

	mgr := NewTsnetManager(cfg)

	if mgr == nil {
		t.Fatal("expected non-nil manager")
	}

	if mgr.server == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestTsnetManager_StartStop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	cfg := &config.Config{
		Tailscale: config.TailscaleConfig{
			StateDir:  t.TempDir(),
			Hostname:  "test-node",
			AuthKey:   "",
			Ephemeral: true,
		},
	}

	mgr := NewTsnetManager(cfg)

	ctx := context.Background()
	err := mgr.Start(ctx)
	if err != nil {
		t.Skipf("cannot start tsnet (needs tailscale backend): %v", err)
	}
	defer mgr.Stop()

	status, err := mgr.Status()
	if err != nil {
		t.Fatalf("Status() failed: %v", err)
	}

	if status == "" {
		t.Error("expected non-empty status")
	}
}
