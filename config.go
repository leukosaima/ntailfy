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
