package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration from environment variables
	config := &Config{
		TailscaleAPIKey:  os.Getenv("TAILSCALE_API_KEY"),
		TailscaleTailnet: os.Getenv("TAILSCALE_TAILNET"),
		NtfyURL:          os.Getenv("NTFY_URL"),
		NtfyAuthToken:    os.Getenv("NTFY_AUTH_TOKEN"),
		NtfyTopic:        os.Getenv("NTFY_TOPIC"),
		PollInterval:     getEnvDuration("POLL_INTERVAL", 60*time.Second),
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	monitor := NewMonitor(config)

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Printf("Starting ntailfy - monitoring tailnet: %s", config.TailscaleTailnet)
	log.Printf("Poll interval: %v", config.PollInterval)

	go monitor.Start()

	<-stop
	log.Println("Shutting down gracefully...")
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("Invalid duration for %s: %v, using default %v", key, err, defaultVal)
		return defaultVal
	}
	return d
}
