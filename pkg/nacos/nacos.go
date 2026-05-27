package nacos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"smartcommunity-microservices/pkg/config"
)

func RegisterService(ctx context.Context, cfg config.NacosConfig, serviceName string, ip string, port int, metadata map[string]string, registerIPOverride ...string) error {
	if !cfg.Enabled {
		return nil
	}
	registerIP := resolveRegisterIP(cfg, ip, registerIPOverride...)


	form := url.Values{}
	form.Set("serviceName", serviceName)
	form.Set("ip", registerIP)
	form.Set("port", strconv.Itoa(port))
	form.Set("namespaceId", cfg.Namespace)
	form.Set("groupName", cfg.Group)
	form.Set("healthy", "true")
	form.Set("enabled", "true")
	form.Set("ephemeral", "true")
	if len(metadata) > 0 {
		form.Set("metadata", encodeMetadata(metadata))
	}

	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, cfg.BaseURL()+"/nacos/v1/ns/instance", nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = form.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("nacos register failed: status=%d body=%s", resp.StatusCode, string(body))
	}
	go keepAlive(ctx, cfg, serviceName, registerIP, port, metadata)
	return nil
}

type beatInfo struct {
	IP          string            `json:"ip"`
	Port        int               `json:"port"`
	ServiceName string            `json:"serviceName"`
	Cluster     string            `json:"cluster"`
	Weight      float64           `json:"weight"`
	Metadata    map[string]string `json:"metadata"`
	Scheduled   bool              `json:"scheduled"`
}

func keepAlive(ctx context.Context, cfg config.NacosConfig, serviceName string, ip string, port int, metadata map[string]string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = sendBeat(ctx, cfg, serviceName, ip, port, metadata)
		}
	}
}

func sendBeat(ctx context.Context, cfg config.NacosConfig, serviceName string, ip string, port int, metadata map[string]string) error {
	beat := beatInfo{
		IP:          ip,
		Port:        port,
		ServiceName: serviceName,
		Cluster:     "DEFAULT",
		Weight:      1,
		Metadata:    metadata,
		Scheduled:   true,
	}
	payload, err := json.Marshal(beat)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Set("serviceName", serviceName)
	form.Set("ip", ip)
	form.Set("port", strconv.Itoa(port))
	form.Set("namespaceId", cfg.Namespace)
	form.Set("groupName", cfg.Group)
	form.Set("ephemeral", "true")
	form.Set("beat", string(payload))

	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPut, cfg.BaseURL()+"/nacos/v1/ns/instance/beat", nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = form.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("nacos beat failed: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

func resolveRegisterIP(cfg config.NacosConfig, ip string, override ...string) string {
	if len(override) > 0 && override[0] != "" {
		return override[0]
	}
	if ip != "" && ip != "0.0.0.0" && ip != "::" {
		return ip
	}
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), time.Second)
	if err != nil {
		return ip
	}
	defer conn.Close()
	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok && addr.IP != nil {
		return addr.IP.String()
	}
	return ip
}

func GetConfig(ctx context.Context, cfg config.NacosConfig, dataID string) (string, error) {
	if !cfg.Enabled {
		return "", nil
	}
	values := url.Values{}
	values.Set("dataId", dataID)
	values.Set("group", cfg.Group)
	values.Set("tenant", cfg.Namespace)
	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, cfg.BaseURL()+"/nacos/v1/cs/configs?"+values.Encode(), nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("nacos config failed: status=%d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func LoadConfig(ctx context.Context, cfg config.NacosConfig, dataID string) (string, error) {
	return GetConfig(ctx, cfg, dataID)
}

func encodeMetadata(metadata map[string]string) string {
	values := url.Values{}
	for k, v := range metadata {
		values.Set(k, v)
	}
	return values.Encode()
}
