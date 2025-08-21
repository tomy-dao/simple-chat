package client

import (
	"encoding/json"
	"local/config"
	"net/http"
	"time"
)

type User struct {
	ID          uint           `json:"id"`
	UserName    string         `json:"username"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Response[T any] struct {
	Data  *T      `json:"data"`
	Error string `json:"error"`
}

func GetMe(token string) (*User, error) {
	req, err := http.NewRequest("GET", config.Config.BackendServerURL+"/api/v1/me", nil)
	
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response Response[User]
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}



