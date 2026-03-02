package main

import (
	"fmt"
	"time"
)

type Config struct {
	TailscaleOAuthClientID     string
	TailscaleOAuthClientSecret string
	TailscaleOAuthScope        string
	TailscaleTailnet           string
	NtfyURL                    string
	NtfyAuthToken              string
	NtfyTopic                  string
	PollInterval               time.Duration
	DeviceFilter               []string // If empty, monitor all devices
}

func (c *Config) Validate() error {
	if c.TailscaleOAuthClientID == "" {
		return fmt.Errorf("TAILSCALE_OAUTH_CLIENT_ID is required")
	}
	if c.TailscaleOAuthClientSecret == "" {
		return fmt.Errorf("TAILSCALE_OAUTH_CLIENT_SECRET is required")
	}
	if c.TailscaleTailnet == "" {
		return fmt.Errorf("TAILSCALE_TAILNET is required")
	}
	if c.NtfyURL == "" {
		return fmt.Errorf("NTFY_URL is required")
	}
	if c.NtfyTopic == "" {
		return fmt.Errorf("NTFY_TOPIC is required")
	}
	if c.PollInterval < 10*time.Second {
		return fmt.Errorf("POLL_INTERVAL must be at least 10s")
	}
	return nil
}

// ShouldMonitorDevice returns true if the device should be monitored based on hostname
func (c *Config) ShouldMonitorDevice(hostname string) bool {
	if len(c.DeviceFilter) == 0 {
		return true // Monitor all devices if no filter specified
	}
	for _, filterHostname := range c.DeviceFilter {
		if filterHostname == hostname {
			return true
		}
	}
	return false
}
