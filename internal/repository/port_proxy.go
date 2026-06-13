package repository

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"tsctl/internal/domain"
	"tsctl/pkg/logger"

	"tailscale.com/tsnet"
)

type portProxy struct {
	listener  net.Listener
	cancel    context.CancelFunc
	port      int
	status    string
	errMsg    string
	localAddr string
	tsnetAddr string
	bytesIn   atomic.Int64
	bytesOut  atomic.Int64
	connCount atomic.Int64
}

type PortProxyManager struct {
	server     *tsnet.Server
	mu         sync.RWMutex
	proxies    map[int]*portProxy
	scanCancel context.CancelFunc
	scanning   bool
}

func NewPortProxyManager(server *tsnet.Server) *PortProxyManager {
	return &PortProxyManager{
		server:  server,
		proxies: make(map[int]*portProxy),
	}
}

func (p *PortProxyManager) StartProxy(ports []int) error {
	for _, port := range ports {
		if err := p.startOne(port); err != nil {
			logger.Log.WithError(err).WithField("port", port).Warn("failed to start proxy")
		}
	}
	return nil
}

func (p *PortProxyManager) startOne(port int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.proxies[port]; exists {
		return nil
	}

	ln, err := p.server.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on tsnet port %d: %w", port, err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	proxy := &portProxy{
		listener:  ln,
		cancel:    cancel,
		port:      port,
		status:    "active",
		localAddr: fmt.Sprintf("localhost:%d", port),
		tsnetAddr: ln.Addr().String(),
	}

	p.proxies[port] = proxy

	go p.acceptLoop(ctx, ln, port, proxy)

	logger.Log.WithFields(map[string]interface{}{
		"port":       port,
		"tsnet_addr": ln.Addr().String(),
	}).Info("port proxy started")

	return nil
}

func (p *PortProxyManager) acceptLoop(ctx context.Context, ln net.Listener, port int, proxy *portProxy) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				logger.Log.WithError(err).WithField("port", port).Warn("accept error on proxy listener")
				return
			}
		}
		go p.handleConn(ctx, conn, port, proxy)
	}
}

func (p *PortProxyManager) handleConn(_ context.Context, tsConn net.Conn, port int, proxy *portProxy) {
	defer func() { _ = tsConn.Close() }()

	localConn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 5*time.Second)
	if err != nil {
		logger.Log.WithError(err).WithField("port", port).Warn("failed to connect to local port")
		return
	}
	defer func() { _ = localConn.Close() }()

	proxy.connCount.Add(1)

	var wg sync.WaitGroup
	wg.Add(2)

	// tsnet -> local
	go func() {
		defer wg.Done()
		n, _ := io.Copy(localConn, tsConn)
		proxy.bytesIn.Add(n)
		// half-close to unblock the other direction
		if tc, ok := localConn.(*net.TCPConn); ok {
			_ = tc.CloseWrite()
		}
	}()

	// local -> tsnet
	go func() {
		defer wg.Done()
		n, _ := io.Copy(tsConn, localConn)
		proxy.bytesOut.Add(n)
		if tc, ok := tsConn.(*net.TCPConn); ok {
			_ = tc.CloseWrite()
		}
	}()

	wg.Wait()
}

func (p *PortProxyManager) StopProxy(ports []int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, port := range ports {
		p.stopLocked(port)
	}
	return nil
}

func (p *PortProxyManager) stopLocked(port int) {
	proxy, exists := p.proxies[port]
	if !exists {
		return
	}

	proxy.cancel()
	_ = proxy.listener.Close()
	delete(p.proxies, port)

	logger.Log.WithField("port", port).Info("port proxy stopped")
}

func (p *PortProxyManager) StopAll() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for port := range p.proxies {
		p.stopLocked(port)
	}
	return nil
}

func (p *PortProxyManager) ActivePorts() []domain.PortProxyInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	infos := make([]domain.PortProxyInfo, 0, len(p.proxies))
	for _, proxy := range p.proxies {
		infos = append(infos, domain.PortProxyInfo{
			Port:      proxy.port,
			Status:    proxy.status,
			LocalAddr: proxy.localAddr,
			TsnetAddr: proxy.tsnetAddr,
			Error:     proxy.errMsg,
			BytesIn:   proxy.bytesIn.Load(),
			BytesOut:  proxy.bytesOut.Load(),
			ConnCount: proxy.connCount.Load(),
		})
	}
	return infos
}

func (p *PortProxyManager) StartAutoScan(interval int, excludePorts []int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.scanning {
		return nil
	}

	excludeSet := make(map[int]bool, len(excludePorts))
	for _, ep := range excludePorts {
		excludeSet[ep] = true
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.scanCancel = cancel
	p.scanning = true

	go p.scanLoop(ctx, interval, excludeSet)

	logger.Log.WithFields(map[string]interface{}{
		"interval":      interval,
		"exclude_ports": excludePorts,
	}).Info("auto port scan started")

	return nil
}

func (p *PortProxyManager) IsAutoScanActive() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.scanning
}

func (p *PortProxyManager) scanLoop(ctx context.Context, interval int, excludeSet map[int]bool) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.runScan(excludeSet)
		}
	}
}

func (p *PortProxyManager) runScan(excludeSet map[int]bool) {
	found, err := scanLocalPorts()
	if err != nil {
		logger.Log.WithError(err).Warn("auto scan failed to list local ports")
		return
	}

	foundSet := make(map[int]bool, len(found))
	for _, port := range found {
		if !excludeSet[port] {
			foundSet[port] = true
		}
	}

	p.mu.RLock()
	activeSet := make(map[int]bool, len(p.proxies))
	for port := range p.proxies {
		activeSet[port] = true
	}
	p.mu.RUnlock()

	// start proxies for newly discovered ports
	for port := range foundSet {
		if !activeSet[port] {
			if err := p.startOne(port); err != nil {
				logger.Log.WithError(err).WithField("port", port).Warn("auto scan: failed to start proxy")
			}
		}
	}

	// stop proxies for ports no longer listening
	toStop := []int{}
	for port := range activeSet {
		if !foundSet[port] {
			toStop = append(toStop, port)
		}
	}
	if len(toStop) > 0 {
		p.mu.Lock()
		for _, port := range toStop {
			p.stopLocked(port)
		}
		p.mu.Unlock()
	}
}

func (p *PortProxyManager) StopAutoScan() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.scanCancel != nil {
		p.scanCancel()
		p.scanCancel = nil
	}
	p.scanning = false

	logger.Log.Info("auto port scan stopped")
}

// scanLocalPorts returns the list of TCP ports currently listening on localhost.
func scanLocalPorts() ([]int, error) {
	if runtime.GOOS == "darwin" {
		return scanLocalPortsDarwin()
	}
	return scanLocalPortsLinux()
}

var lsofPortRe = regexp.MustCompile(`\*:(\d+)\s+\(LISTEN\)`)

func scanLocalPortsDarwin() ([]int, error) {
	out, err := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P").Output()
	if err != nil {
		return nil, fmt.Errorf("lsof failed: %w", err)
	}

	seen := make(map[int]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		m := lsofPortRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		port, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		seen[port] = true
	}

	ports := make([]int, 0, len(seen))
	for port := range seen {
		ports = append(ports, port)
	}
	return ports, nil
}

// /proc/net/tcp hex-encodes local_address as hex_ip:hex_port; state 0A = LISTEN.
func scanLocalPortsLinux() ([]int, error) {
	out, err := exec.Command("cat", "/proc/net/tcp").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/net/tcp: %w", err)
	}

	seen := make(map[int]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	// skip header line
	if scanner.Scan() {
		_ = scanner.Text()
	}

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		// state field: 0A = TCP_LISTEN
		if fields[3] != "0A" {
			continue
		}
		// local_address: "XXXXXXXX:PPPP"
		parts := strings.SplitN(fields[1], ":", 2)
		if len(parts) != 2 {
			continue
		}
		portHex := parts[1]
		port64, err := strconv.ParseInt(portHex, 16, 32)
		if err != nil {
			continue
		}
		seen[int(port64)] = true
	}

	ports := make([]int, 0, len(seen))
	for port := range seen {
		ports = append(ports, port)
	}
	return ports, nil
}
