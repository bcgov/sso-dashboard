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

type TokenStrategy interface {
	GetTokens(data url.Values, authUrl string) (string, string, error)
	IsTokenExpired(token string) (bool, error)
}

type RequestHandler struct {
	ApiBaseUrl  string
	AuthBaseUrl string

	AccessToken  string
	RefreshToken string

	password string
	username string
	clientId string

	tokenStrategy TokenStrategy
}

func NewRequestHandler(
	tokenStrategy TokenStrategy,
	ApiUrl string,
	AuthUrl string,
	password string,
	username string,
	clientId string,
) *RequestHandler {
	return &RequestHandler{
		tokenStrategy: tokenStrategy,
		ApiBaseUrl:    ApiUrl,
		AuthBaseUrl:   AuthUrl,
		password:      password,
		username:      username,
		clientId:      clientId,
	}
}

/*
Method to make a post request for a new token or to refresh a token. Pass in the relevant data for the request type,
e.g grant_type=password/client_credentials for a new token, grant_type=refresh_token for a refresh token.
*/
func (tm *RequestHandler) GetTokens(data url.Values, authUrl string) (string, string, error) {
	req, err := http.NewRequest("POST", authUrl+"/realms/master/protocol/openid-connect/token", strings.NewReader(data.Encode()))

	if err != nil {
		log.Printf("Error occurred creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error occurred sending request to API endpoint: %v", err)
		return "", "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("non 200 status code returned from token request: %v", resp.Status)
		return "", "", errors.New("non 200 status code returned from token request")
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", "", err
	}

	type TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	var tokenResponse TokenResponse

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		log.Printf("Error parsing JSON response: %v", err)
		return "", "", err
	}

	return tokenResponse.AccessToken, tokenResponse.RefreshToken, nil
}

// Method to get the token, refreshing if necessary
func (rm *RequestHandler) GetToken() (string, string, error) {
	if rm.AccessToken == "" {
		formData := url.Values{
			"grant_type": {"password"},
			"client_id":  {rm.clientId},
			"password":   {rm.password},
			"username":   {rm.username},
		}

		accessToken, refreshToken, err := rm.tokenStrategy.GetTokens(formData, rm.AuthBaseUrl)
		return accessToken, refreshToken, err
	}

	// Check expiry and refresh if necessary
	accessTokenExpired, err := rm.IsTokenExpired(rm.AccessToken)

	if err != nil {
		return "", "", err
	}

	if !accessTokenExpired {
		return rm.AccessToken, rm.RefreshToken, nil
	}

	refreshTokenExpired, err := rm.IsTokenExpired(rm.RefreshToken)

	var formData url.Values

	if refreshTokenExpired || err != nil {
		formData = url.Values{
			"grant_type": {"password"},
			"client_id":  {rm.clientId},
			"password":   {rm.password},
			"username":   {rm.username},
		}
	} else {
		formData = url.Values{
			"grant_type":    {"refresh_token"},
			"client_id":     {rm.clientId},
			"refresh_token": {rm.RefreshToken},
		}
	}

	accessToken, refreshToken, err := rm.tokenStrategy.GetTokens(formData, rm.AuthBaseUrl)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Checks if access token is expired
func (tm *RequestHandler) IsTokenExpired(token string) (bool, error) {
	if token == "" {
		return true, nil
	}

	parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			now := time.Now().UTC()
			expire := time.Unix(int64(exp), 0).UTC()
			return expire.Before(now), nil
		} else {
			fmt.Println("exp claim not found or not a float64")
			return false, errors.New("cannot read exp claim")
		}
	} else {
		log.Printf("Invalid token claims")
		return false, errors.New("Invalid token claims")
	}
}

// Method to perform an HTTP request with token management
func (rm *RequestHandler) DoRequest(req *http.Request) (*http.Response, error) {
	accessToken, refreshToken, err := rm.GetToken()

	if err != nil {
		log.Printf("Error in dorequest")
		return nil, err
	}

	rm.AccessToken = accessToken
	rm.RefreshToken = refreshToken

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
