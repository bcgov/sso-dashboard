package model

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"sso-dashboard.bcgov.com/aggregator/keycloak"
	"sso-dashboard.bcgov.com/aggregator/webhooks"
)

type MockRocketChat struct {
	Messages [][]string
}

func (m *MockRocketChat) NotifyRocketChat(text string, title string, body string) error {
	message := []string{text, title, body}
	m.Messages = append(m.Messages, message)
	return nil
}

func (m *MockRocketChat) ResetMock() {
	m.Messages = [][]string{}
}

func TestTokenFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/realms/master/protocol/openid-connect/token" {
			// Fail all token requests
			w.WriteHeader(http.StatusForbidden)
		}
	}))
	defer server.Close()

	mock := &MockRocketChat{}
	ActiveSessions("DEV", server.URL, "client", "user", "pass", mock)

	if len(mock.Messages) != 1 {
		t.Errorf("Expected 1 message to be sent to RocketChat")
	}
	rcTitle := mock.Messages[0][1]
	rcError := mock.Messages[0][2]
	if rcTitle != fmt.Sprintf(RealmErrorMessage, "DEV") {
		t.Errorf("Expected rocket chat message to include the environment name")
	}
	if rcError != "non 200 status code returned from token request" {
		log.Print(rcError)
		t.Errorf("Expected rocket chat message to include the error message")
	}
}

func TestRealmNotifications(t *testing.T) {
	var realmResponse string
	var realmStatusCode int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validToken := keycloak.GenerateJWT(time.Now().Add(time.Hour).Unix())
		if r.URL.Path == "/realms/master/protocol/openid-connect/token" {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(fmt.Sprintf(`{"access_token":"%s", "refresh_token":"%s"}`, validToken, validToken)))
		}
		if r.URL.Path == "/admin/realms" {
			w.WriteHeader(realmStatusCode)
			w.Write([]byte(realmResponse))
		}
	}))
	defer server.Close()

	mock := &MockRocketChat{}

	// Test when realm response is successful. Should send a successful message to RC
	realmResponse = `[]`
	realmStatusCode = http.StatusOK
	ActiveSessions("DEV", server.URL, "client", "user", "pass", mock)

	if len(mock.Messages) != 1 {
		t.Errorf("Expected 1 message to be sent to RocketChat")
	}
	rcText := mock.Messages[0][0]
	rcTitle := mock.Messages[0][1]

	if rcText != "Session Data Loaded Successfully" {
		t.Errorf("Expected successfuly rocket chat message")
	}
	if rcTitle != "DEV" {
		t.Errorf("Expected rocket chat message to include the environment name")
	}

	// When realm response is not ok, should notify rc
	mock.ResetMock()
	realmResponse = `[]`
	realmStatusCode = http.StatusForbidden

	ActiveSessions("TEST", server.URL, "client", "user", "pass", mock)

	if len(mock.Messages) != 1 {
		t.Errorf("Expected 1 message to be sent to RocketChat")
	}
	rcTitle = mock.Messages[0][1]
	if rcTitle != fmt.Sprintf(RealmErrorMessage, "TEST") {
		print(rcTitle)
		t.Errorf("Expected rocket chat message to include the environment name")
	}
}

func TestClientNotifications(t *testing.T) {
	var clientResponse string
	var clientStatusCode int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validToken := keycloak.GenerateJWT(time.Now().Add(time.Hour).Unix())
		if r.URL.Path == "/realms/master/protocol/openid-connect/token" {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(fmt.Sprintf(`{"access_token":"%s", "refresh_token":"%s"}`, validToken, validToken)))
		}
		if r.URL.Path == "/admin/realms" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"realm": "realm 1"}, {"realm": "realm 2"}]`))
		}
		if strings.HasSuffix(r.URL.Path, "/client-session-stats") {
			w.WriteHeader(clientStatusCode)
			w.Write([]byte(clientResponse))
		}
	}))
	defer server.Close()
	mock := &MockRocketChat{}

	// Test when realm response is successful. Should send a successful message to RC
	clientResponse = `[]`
	clientStatusCode = http.StatusOK
	ActiveSessions("DEV", server.URL, "client", "user", "pass", mock)

	if len(mock.Messages) != 1 {
		t.Errorf("Expected 1 message to be sent to RocketChat")
	}
	rcText := mock.Messages[0][0]
	rcTitle := mock.Messages[0][1]

	if rcText != "Session Data Loaded Successfully" {
		t.Errorf("Expected successfuly rocket chat message")
	}
	if rcTitle != "DEV" {
		t.Errorf("Expected rocket chat message to include the environment name")
	}

	// When client response is not ok, should notify rc
	mock.ResetMock()
	clientStatusCode = http.StatusForbidden

	ActiveSessions("TEST", server.URL, "client", "user", "pass", mock)

	if len(mock.Messages) != 1 {
		t.Errorf("Expected 1 message to be sent to RocketChat")
	}
	rcTitle = mock.Messages[0][1]
	if rcTitle != fmt.Sprintf(ClientErrorMessage, "TEST") {
		t.Errorf("Expected rocket chat message to include the environment name")
	}
}

func TestClientContinue(t *testing.T) {
	var requestCount int
	var clientStatusCode int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validToken := keycloak.GenerateJWT(time.Now().Add(time.Hour).Unix())
		if r.URL.Path == "/realms/master/protocol/openid-connect/token" {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(fmt.Sprintf(`{"access_token":"%s", "refresh_token":"%s"}`, validToken, validToken)))
		}
		if strings.HasSuffix(r.URL.Path, "/client-session-stats") {
			requestCount++
			w.WriteHeader(clientStatusCode)
		}
	}))
	defer server.Close()
	rm := keycloak.NewRequestHandler(&keycloak.RequestHandler{}, server.URL, server.URL, "", "", "")

	clientStatusCode = http.StatusBadGateway
	GetClientStats(rm, []string{"realm 1", "realm 2", "realm 3"}, "dev")
	if requestCount != 3 {
		t.Errorf("Expeted all realm requests to be attempted even if one fails")
	}
}

func TestNotifyRocketChat(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	// Run a closed server simulating rocket chat being down
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()
	os.Setenv("RC_WEBHOOK", server.URL)

	rc := &webhooks.RocketChat{}

	// Notify should return error, not throw error
	err := rc.NotifyRocketChat("Test", "Title", "Body")

	if err == nil {
		t.Fatal("Expected a returned error when RC Server is down")
	}

	logOutput := logBuf.String()

	if !strings.Contains(logOutput, "Error sending rocket chat notification") {
		t.Errorf("Expected log to contain 'Error sending rocket chat notification', got: %s", logOutput)
	}
}
