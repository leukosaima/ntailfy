package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type NtfyClient struct {
	baseURL   string
	topic     string
	authToken string
	client    *http.Client
}

type NtfyMessage struct {
	Title    string   `json:"title,omitempty"`
	Message  string   `json:"message"`
	Priority int      `json:"priority,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

func NewNtfyClient(baseURL, topic, authToken string) *NtfyClient {
	return &NtfyClient{
		baseURL:   strings.TrimSuffix(baseURL, "/"),
		topic:     topic,
		authToken: authToken,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (n *NtfyClient) Send(msg *NtfyMessage) error {
	url := fmt.Sprintf("%s/%s", n.baseURL, n.topic)
	
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	if n.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+n.authToken)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy returned status %d", resp.StatusCode)
	}

	return nil
}
