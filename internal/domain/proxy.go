package domain

// PortProxyManager manages TCP port proxying through the tailscale network
type PortProxyManager interface {
	// StartProxy starts proxying the given ports through tsnet
	StartProxy(ports []int) error
	// StopProxy stops proxying the given ports
	StopProxy(ports []int) error
	// StopAll stops all active proxies
	StopAll() error
	// ActivePorts returns list of currently proxied ports
	ActivePorts() []PortProxyInfo
	// StartAutoScan starts automatic port discovery and proxying
	StartAutoScan(interval int, excludePorts []int) error
	// StopAutoScan stops automatic port discovery
	StopAutoScan()
	// IsAutoScanActive returns true if auto scan is currently running
	IsAutoScanActive() bool
}

// ProxyStatusResponse represents the overall status of the proxy subsystem
type ProxyStatusResponse struct {
	IsAutoScanActive bool            `json:"is_auto_scan_active"`
	Proxies          []PortProxyInfo `json:"proxies"`
}

// PortProxyInfo represents info about a proxied port
type PortProxyInfo struct {
	Port      int    `json:"port"`
	Status    string `json:"status"`
	LocalAddr string `json:"local_addr"`
	TsnetAddr string `json:"tsnet_addr"`
	Error     string `json:"error,omitempty"`
	BytesIn   int64  `json:"bytes_in"`
	BytesOut  int64  `json:"bytes_out"`
	ConnCount int64  `json:"conn_count"`
}
