package main

import (
	"log"
	"time"
)

type Monitor struct {
	config         *Config
	tailscale      *TailscaleClient
	ntfy           *NtfyClient
	previousStates map[string]bool // deviceID -> online status
}

func NewMonitor(config *Config) *Monitor {
	return &Monitor{
		config:         config,
		tailscale:      NewTailscaleClient(config.TailscaleAPIKey, config.TailscaleTailnet),
		ntfy:           NewNtfyClient(config.NtfyURL, config.NtfyTopic, config.NtfyAuthToken),
		previousStates: make(map[string]bool),
	}
}

func (m *Monitor) Start() {
	ticker := time.NewTicker(m.config.PollInterval)
	defer ticker.Stop()

	// Initial poll
	m.poll()

	for range ticker.C {
		m.poll()
	}
}

func (m *Monitor) poll() {
	devices, err := m.tailscale.GetDevices()
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		return
	}

	currentStates := make(map[string]bool)
	
	for _, device := range devices {
		currentStates[device.ID] = device.Online
		
		previousOnline, existed := m.previousStates[device.ID]
		
		if !existed {
			// New device discovered
			log.Printf("New device discovered: %s (%s) - %s", device.Name, device.Hostname, onlineStatus(device.Online))
			continue
		}
		
		// Check for state change
		if previousOnline != device.Online {
			m.notifyStateChange(device)
		}
	}

	// Detect devices that disappeared
	for id, wasOnline := range m.previousStates {
		if _, exists := currentStates[id]; !exists && wasOnline {
			log.Printf("Device removed from tailnet: %s", id)
		}
	}

	m.previousStates = currentStates
}

func (m *Monitor) notifyStateChange(device Device) {
	status := "disconnected"
	priority := 3 // default
	tags := []string{"tailscale"}
	
	if device.Online {
		status = "connected"
		priority = 3
		tags = append(tags, "connected", "green_circle")
	} else {
		tags = append(tags, "disconnected", "red_circle")
	}

	message := &NtfyMessage{
		Title:    device.Name + " " + status,
		Message:  device.Hostname + " is now " + status,
		Priority: priority,
		Tags:     tags,
	}

	if err := m.ntfy.Send(message); err != nil {
		log.Printf("Error sending notification: %v", err)
	} else {
		log.Printf("Notification sent: %s is %s", device.Name, status)
	}
}

func onlineStatus(online bool) string {
	if online {
		return "online"
	}
	return "offline"
}
