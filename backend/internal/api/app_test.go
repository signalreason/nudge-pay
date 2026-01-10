package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"nudgepay/internal/api"
	"nudgepay/internal/config"
	"nudgepay/internal/db"
)

type registerResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
	Org struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"org"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type createResponse struct {
	ID string `json:"id"`
}

type outboxResponse struct {
	Outbox []struct {
		ID string `json:"id"`
	} `json:"outbox"`
}

func newTestApp(t *testing.T) (*fiber.App, func()) {
	t.Helper()
	database, err := db.New(":memory:")
	if err != nil {
		t.Fatalf("db error: %v", err)
	}
	cfg := config.Config{JWTSecret: "test-secret", WorkerEnabled: false}
	app := api.NewApp(database, cfg)
	cleanup := func() {
		_ = database.Close()
	}
	return app, cleanup
}

func TestRegisterAndLogin(t *testing.T) {
	app, cleanup := newTestApp(t)
	defer cleanup()

	registerBody := map[string]string{
		"email":    "owner@example.com",
		"password": "password123",
		"org_name": "Studio One",
	}
	resp := performRequest(t, app, "POST", "/api/auth/register", registerBody, "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var reg registerResponse
	decodeJSON(t, resp, &reg)
	if reg.Token == "" {
		t.Fatalf("expected token")
	}

	loginBody := map[string]string{"email": "owner@example.com", "password": "password123"}
	loginResp := performRequest(t, app, "POST", "/api/auth/login", loginBody, "")
	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", loginResp.StatusCode)
	}
	var login loginResponse
	decodeJSON(t, loginResp, &login)
	if login.Token == "" {
		t.Fatalf("expected login token")
	}
}

func TestInvoiceReminderFlow(t *testing.T) {
	app, cleanup := newTestApp(t)
	defer cleanup()

	registerBody := map[string]string{
		"email":    "owner@example.com",
		"password": "password123",
		"org_name": "Studio One",
	}
	resp := performRequest(t, app, "POST", "/api/auth/register", registerBody, "")
	var reg registerResponse
	decodeJSON(t, resp, &reg)
	if reg.Token == "" {
		t.Fatalf("expected token")
	}

	clientBody := map[string]string{
		"name":    "Jamie Client",
		"email":   "client@example.com",
		"company": "ClientCo",
	}
	clientResp := performRequest(t, app, "POST", "/api/clients", clientBody, reg.Token)
	if clientResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", clientResp.StatusCode)
	}
	var client createResponse
	decodeJSON(t, clientResp, &client)

		invoiceBody := map[string]interface{}{
			"client_id":        client.ID,
			"number":           "INV-100",
			"amount_cents":     125000,
			"currency":         "usd",
			"due_date":         time.Now().UTC().Add(-24 * time.Hour).Format("2006-01-02"),
			"reminder_offsets": []int{0},
		}
	invoiceResp := performRequest(t, app, "POST", "/api/invoices", invoiceBody, reg.Token)
	if invoiceResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", invoiceResp.StatusCode)
	}

	sendResp := performRequest(t, app, "POST", "/api/reminders/send-due", nil, reg.Token)
	if sendResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", sendResp.StatusCode)
	}

	outboxResp := performRequest(t, app, "GET", "/api/outbox", nil, reg.Token)
	if outboxResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", outboxResp.StatusCode)
	}
	var outbox outboxResponse
	decodeJSON(t, outboxResp, &outbox)
	if len(outbox.Outbox) == 0 {
		t.Fatalf("expected outbox entry")
	}
}

func performRequest(t *testing.T, app *fiber.App, method, path string, body interface{}, token string) *http.Response {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode error: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, dst interface{}) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		t.Fatalf("decode error: %v", err)
	}
}
