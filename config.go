package main

import (
	"fmt"
	"time"
)

type Config struct {
	TailscaleAPIKey  string
	TailscaleTailnet string
	NtfyURL          string
	NtfyAuthToken    string
	NtfyTopic        string
	PollInterval     time.Duration
	DeviceFilter     []string // If empty, monitor all devices
}

func (c *Config) Validate() error {
	if c.TailscaleAPIKey == "" {
		return fmt.Errorf("TAILSCALE_API_KEY is required")
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

// ShouldMonitorDevice returns true if the device should be monitored
func (c *Config) ShouldMonitorDevice(deviceName string) bool {
	if len(c.DeviceFilter) == 0 {
		return true // Monitor all devices if no filter specified
	}
	for _, name := range c.DeviceFilter {
		if name == deviceName {
			return true
		}
	}
	return false
}
