package client

import (
	"bytes"
	"encoding/json"
	"local/config"
	"local/model"
	"log"
	"net/http"
	"time"
)
type SocketClient interface {
	Broadcast(message *model.BroadcastMessage)
}

type socketClient struct {
	token string
}

func (c *socketClient) Broadcast(message *model.BroadcastMessage) {
	// Create HTTP request to socket server
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Convert message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	// Create request
	req, err := http.NewRequest("POST", config.Config.SocketServerURL+"/broadcast", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending broadcast request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Broadcast request failed with status: %d", resp.StatusCode)
	}
}

func NewSocketClient() SocketClient {
	return &socketClient{token: config.Config.SocketToken}
}
