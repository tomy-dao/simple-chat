package event

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)


var ChatPath = "/chat"


func mapData(payload any, data any) error {
	jsonStr, _ := json.Marshal(payload)
	err := json.Unmarshal(jsonStr, data)
	if err != nil {
		return err
	}
	return nil
}

func decodeJWT(tokenString string) (map[string]interface{}, error) {
	// Split the token into parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	// Decode the payload (second part)
	payload := parts[1]
	
	// Add padding if necessary
	if len(payload)%4 != 0 {
		payload += strings.Repeat("=", 4-len(payload)%4)
	}

	// Base64 decode
	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var claims map[string]interface{}
	err = json.Unmarshal(decoded, &claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}


