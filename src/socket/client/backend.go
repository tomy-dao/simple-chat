package client

import (
	"encoding/json"
	"local/config"
	"net/http"
)

func GetMe(token string) (any, error) {
	req, err := http.NewRequest("GET", config.Config.BackendServerURL+"/me", nil)
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

	var me any
	err = json.NewDecoder(resp.Body).Decode(&me)
	if err != nil {
		return nil, err
	}

	return me, nil
}



