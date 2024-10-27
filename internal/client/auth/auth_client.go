package auth

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

type Response struct {
	Status    string `json:"status"`
	UserLogin string `json:"user-login"`
}

var authServiceURL = getAuthServiceURL()

func getAuthServiceURL() string {
	_ = godotenv.Load()
	return os.Getenv("AUTH_SERVICE_URL")
}

func ValidateToken(token string) (*Response, error) {
	req, err := http.NewRequest("GET", authServiceURL, nil)
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

	var authResponse Response
	if err = render.DecodeJSON(resp.Body, &authResponse); err != nil {
		return nil, err
	}
	if authResponse.Status != "OK" {
		return nil, fmt.Errorf("failed request to auth service")
	}
	return &authResponse, nil
}
