package keycloak

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	token        string
	refreshToken string
	password     string
	username     string
	BaseUrl      string
	clientId     string
}

func NewTokenManager(clientId string, password string, username string, baseUrl string) *TokenManager {
	return &TokenManager{
		clientId: clientId,
		password: password,
		username: username,
		BaseUrl:  baseUrl,
	}
}

/*
Method to make a post request for a new token or to refresh a token. Pass in the relevant data for the request type,
e.g grant_type=password/client_credentials for a new token, grant_type=refresh_token for a refresh token.
*/
func (tm *TokenManager) getTokens(data url.Values) (string, error) {
	req, err := http.NewRequest("POST", tm.BaseUrl+"/realms/master/protocol/openid-connect/token", strings.NewReader(data.Encode()))

	if err != nil {
		log.Fatalf("Error occurred creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error occurred sending request to API endpoint: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	type TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	var tokenResponse TokenResponse

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		log.Fatalf("Error parsing JSON response: %v", err)
	}

	tm.token = tokenResponse.AccessToken
	tm.refreshToken = tokenResponse.RefreshToken
	return tm.token, nil
}

// Method to get the token, refreshing if necessary
func (tm *TokenManager) getToken() (string, error) {
	// Fetch a new token if not set yet
	if tm.token == "" {
		formData := url.Values{
			"grant_type": {"password"},
			"client_id":  {tm.clientId},
			"password":   {tm.password},
			"username":   {tm.username},
		}

		token, err := tm.getTokens(formData)
		return token, err
	}

	// Check expiry and refresh if necessary
	expired, err := tm.IsTokenExpired()
	if err != nil {
		return "", err
	}

	if !expired {
		return tm.token, nil
	}

	err = tm.refreshAccessToken()
	if err != nil {
		return "", err
	}

	return tm.token, nil
}

// Checks if access token is expired
func (tm *TokenManager) IsTokenExpired() (bool, error) {
	if tm.token == "" {
		return true, nil
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tm.token, jwt.MapClaims{})
	if err != nil {
		log.Fatalf("Error parsing token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			now := time.Now().UTC()
			expire := time.Unix(int64(exp), 0).UTC()
			return expire.Before(now), nil
		} else {
			fmt.Println("exp claim not found or not a float64")
			return false, errors.New("cannot read exp claim")
		}
	} else {
		log.Fatalf("Invalid token claims")
		return false, errors.New("Invalid token claims")
	}
}

// Method to refresh the token
func (tm *TokenManager) refreshAccessToken() error {
	formData := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {tm.clientId},
		"refresh_token": {tm.refreshToken},
	}

	tm.getTokens(formData)
	return nil
}

// Method to perform an HTTP request with token management
func (tm *TokenManager) DoRequest(req *http.Request) (*http.Response, error) {
	token, err := tm.getToken()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		// Try to refresh the token and retry the request
		err := tm.refreshAccessToken()
		if err != nil {
			return nil, err
		}

		// Update the token in the request and retry
		req.Header.Set("Authorization", "Bearer "+tm.token)
		resp, err = client.Do(req)
	}

	return resp, err
}
