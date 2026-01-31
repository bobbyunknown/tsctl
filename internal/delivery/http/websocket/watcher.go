package websocket

import (
	"context"
	"encoding/json"
	"time"

	"tsctl/internal/domain"
	"tsctl/pkg/logger"
)

type WSEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type StatusWatcher struct {
	hub     *Hub
	service domain.TailscaleService
}

func NewStatusWatcher(hub *Hub, service domain.TailscaleService) *StatusWatcher {
	return &StatusWatcher{
		hub:     hub,
		service: service,
	}
}

func (w *StatusWatcher) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var lastAuthStatus string

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			status, err := w.service.GetFullStatus(ctx)
			if err != nil {
				logger.Log.WithError(err).Debug("failed to get status")
				continue
			}

			currentStatus := status.BackendState
			if currentStatus != lastAuthStatus {
				event := WSEvent{
					Type: "auth_status_changed",
					Data: map[string]interface{}{
						"authenticated": currentStatus == "Running",
						"backend_state": currentStatus,
						"auth_url":      status.AuthURL,
						"node_key":      "",
						"hostname":      "",
						"ips":           []string{},
					},
				}

				if status.Self != nil {
					event.Data.(map[string]interface{})["node_key"] = status.Self.PublicKey.String()
					event.Data.(map[string]interface{})["hostname"] = status.Self.HostName
				}

				if len(status.TailscaleIPs) > 0 {
					ips := make([]string, len(status.TailscaleIPs))
					for i, ip := range status.TailscaleIPs {
						ips[i] = ip.String()
					}
					event.Data.(map[string]interface{})["ips"] = ips
				}

				data, _ := json.Marshal(event)
				w.hub.Broadcast(data)

				lastAuthStatus = currentStatus
			}
		}
	}
}
