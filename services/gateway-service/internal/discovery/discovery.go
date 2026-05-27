package discovery

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"smartcommunity-microservices/pkg/config"
)

// serviceEntry holds the resolved URL and its source for a single service.
type serviceEntry struct {
	URL    string // e.g. "http://user-service:8001"
	Source string // "nacos" or "config"
}

// Resolver resolves service names to base URLs.
// It tries Nacos first and falls back to local config, tracking the source per service.
type Resolver struct {
	nacosCfg config.NacosConfig
	localCfg map[string]string
	log      *slog.Logger

	mu       sync.RWMutex
	services map[string]serviceEntry
}

// NewResolver creates a Resolver that tries Nacos then falls back to local config.
func NewResolver(nacosCfg config.NacosConfig, localServices map[string]string, log *slog.Logger) *Resolver {
	r := &Resolver{
		nacosCfg: nacosCfg,
		localCfg: localServices,
		log:      log,
		services: make(map[string]serviceEntry),
	}
	// Seed with local config as baseline
	for name, url := range localServices {
		r.services[name] = serviceEntry{URL: url, Source: "config"}
	}
	// Best-effort Nacos fetch on startup
	r.refreshNacos()
	return r
}

// Resolve returns the base URL for a service (e.g. "http://user-service:8001").
// Returns empty string if not found.
func (r *Resolver) Resolve(serviceName string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if entry, ok := r.services[serviceName]; ok {
		return entry.URL
	}
	return ""
}

// Source returns the overall source: "nacos" if ANY service came from Nacos, else "config".
// For per-service source, use Services().
func (r *Resolver) Source() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, entry := range r.services {
		if entry.Source == "nacos" {
			return "nacos"
		}
	}
	return "config"
}

// Services returns all known service mappings with their source.
func (r *Resolver) Services() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]string, len(r.services))
	for name, entry := range r.services {
		result[name] = entry.URL + " (" + entry.Source + ")"
	}
	return result
}

// StartRefresh starts a background goroutine that refreshes Nacos every 30s.
func (r *Resolver) StartRefresh(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				r.refreshNacos()
			}
		}
	}()
}

func (r *Resolver) refreshNacos() {
	if !r.nacosCfg.Enabled {
		return
	}

	knownServices := []string{
		"user-service", "mall-service", "community-service", "agent-service",
	}

	nacosOK := 0
	nacosResults := make(map[string]string)

	for _, name := range knownServices {
		url, err := r.queryNacos(name)
		if err != nil {
			r.log.Warn("nacos lookup failed", "service", name, "error", err)
			continue
		}
		if url != "" {
			nacosResults[name] = url
			nacosOK++
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if nacosOK > 0 {
		// Update services from Nacos results; fall back to config for missing ones
		for _, name := range knownServices {
			if url, ok := nacosResults[name]; ok {
				r.services[name] = serviceEntry{URL: url, Source: "nacos"}
			} else if cfgURL, ok := r.localCfg[name]; ok {
				// Nacos didn't return this service — fall back to config
				r.services[name] = serviceEntry{URL: cfgURL, Source: "config"}
				r.log.Warn("nacos missing service, falling back to config", "service", name)
			}
		}
		r.log.Info("nacos services resolved", "count", nacosOK)
	} else {
		// Nacos returned nothing — revert all services to config
		for _, name := range knownServices {
			if cfgURL, ok := r.localCfg[name]; ok {
				r.services[name] = serviceEntry{URL: cfgURL, Source: "config"}
			}
		}
		r.log.Warn("nacos returned 0 services, all reverted to local config")
	}
}

func (r *Resolver) queryNacos(serviceName string) (string, error) {
	baseURL := r.nacosCfg.BaseURL()
	url := fmt.Sprintf("%s/nacos/v1/ns/instance/list?serviceName=%s&namespaceId=%s&groupName=%s&healthyOnly=true",
		baseURL, serviceName, r.nacosCfg.Namespace, r.nacosCfg.Group)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("nacos status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip, port := parseInstance(string(body))
	if ip == "" || port == 0 {
		return "", nil
	}
	return fmt.Sprintf("http://%s:%d", ip, port), nil
}

// parseInstance extracts the first ip:port from Nacos instance list JSON.
func parseInstance(jsonStr string) (string, int) {
	hostsIdx := indexOf(jsonStr, `"hosts"`)
	if hostsIdx < 0 {
		return "", 0
	}
	slice := jsonStr[hostsIdx:]

	ipIdx := indexOf(slice, `"ip"`)
	if ipIdx < 0 {
		return "", 0
	}
	ipVal := extractQuotedValue(slice[ipIdx:])
	if ipVal == "" {
		return "", 0
	}

	portIdx := indexOf(slice, `"port"`)
	if portIdx < 0 {
		return ipVal, 0
	}
	portStr := extractNumber(slice[portIdx+6:])
	port := 0
	for _, c := range portStr {
		if c >= '0' && c <= '9' {
			port = port*10 + int(c-'0')
		} else {
			break
		}
	}
	return ipVal, port
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func extractQuotedValue(s string) string {
	start := indexOf(s, `"`)
	if start < 0 {
		return ""
	}
	colonIdx := indexOf(s[start:], `:`)
	if colonIdx < 0 {
		return ""
	}
	afterColon := s[start+colonIdx+1:]
	q1 := indexOf(afterColon, `"`)
	if q1 < 0 {
		return ""
	}
	rest := afterColon[q1+1:]
	q2 := indexOf(rest, `"`)
	if q2 < 0 {
		return ""
	}
	return rest[:q2]
}

func extractNumber(s string) string {
	for i, c := range s {
		if (c >= '0' && c <= '9') || c == '-' {
			continue
		}
		if i > 0 {
			return s[:i]
		}
		return ""
	}
	return s
}
