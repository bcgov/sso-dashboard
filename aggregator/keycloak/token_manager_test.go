package keycloak

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

type MockTokenProvider struct {
	TokenRefreshed    bool
	NewTokenRequested bool
}

func (tm *MockTokenProvider) ResetMock() {
	tm.TokenRefreshed = false
	tm.NewTokenRequested = false
}

func (m *MockTokenProvider) GetTokens(data url.Values, url string) (string, string, error) {
	grantType := data.Get("grant_type")

	if grantType == "refresh_token" {
		m.TokenRefreshed = true
	}

	if grantType == "password" {
		m.NewTokenRequested = true
	}

	return "", "", nil
}

func (m *MockTokenProvider) IsTokenExpired(token string) (bool, error) {
	return false, nil
}

func TestTokenManagerHandler(t *testing.T) {
	mockTokenProvider := &MockTokenProvider{}
	handler := NewRequestHandler(mockTokenProvider, "", "", "", "", "")

	req, _ := http.NewRequest("GET", "http://somedomain.com", nil)

	// Make a request with expired access token and valid refresh token, expect refresh callout
	handler.AccessToken = GenerateJWT(time.Now().Add(-time.Hour).Unix())
	handler.RefreshToken = GenerateJWT(time.Now().Add(time.Hour).Unix())
	handler.DoRequest(req)

	if !mockTokenProvider.TokenRefreshed {
		t.Errorf("expected RefreshToken to be called, but it was not")
	}

	// Make a request with expired access token and an expired refresh token, expect new token requested.
	mockTokenProvider.ResetMock()
	handler.AccessToken = ""
	handler.RefreshToken = GenerateJWT(time.Now().Add(-time.Hour).Unix())
	handler.DoRequest(req)

	if !mockTokenProvider.NewTokenRequested {
		t.Errorf("expected a new token to be requested, but it was not")
	}
	if mockTokenProvider.TokenRefreshed {
		t.Errorf("expected only a new token to be requested, but a refresh callout was made")
	}

	// Make a request with a valid access token, expect no refresh callout
	mockTokenProvider.ResetMock()
	handler.AccessToken = GenerateJWT(time.Now().Add(time.Hour).Unix())
	handler.DoRequest(req)

	if mockTokenProvider.TokenRefreshed || mockTokenProvider.NewTokenRequested {
		t.Errorf("expected existing token to be used, but a new token was requested")
	}

	// Make a request with no tokens, expect new token callout
	mockTokenProvider.ResetMock()
	handler.AccessToken = ""
	handler.RefreshToken = ""
	handler.DoRequest(req)

	if !mockTokenProvider.NewTokenRequested {
		t.Errorf("expected a new token to be requested, but it was not")
	}
}
