package platform5

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer(status int, body any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range map[string]string{
			"X-Request-ID": "req-123",
		} {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
}

func TestNewClient(t *testing.T) {
	c := New(Config{APIKey: "test-key"})
	if c == nil {
		t.Fatal("client is nil")
	}
}

func TestHealthOK(t *testing.T) {
	srv := setupTestServer(http.StatusOK, map[string]any{
		"success": true,
		"message": "ok",
		"data":    nil,
	})
	defer srv.Close()

	c := New(Config{APIKey: "test", BaseURL: srv.URL})
	if err := c.Health(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestHealthError(t *testing.T) {
	srv := setupTestServer(http.StatusUnauthorized, map[string]any{
		"success": false,
		"message": "Invalid API key",
		"errors":  "bad_key",
	})
	defer srv.Close()

	c := New(Config{APIKey: "bad", BaseURL: srv.URL})
	err := c.Health(ctx)
	if err == nil {
		t.Fatal("expected error")
	}

	var p5err *Error
	if !as(err, &p5err) {
		t.Fatal("expected *Error")
	}
	if !p5err.IsUnauthorized() {
		t.Fatal("expected unauthorized")
	}
}

func TestRateLimitError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "50")
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.Header().Set("X-Request-ID", "req-123")
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"message": "Rate limited",
			"errors":  "too_many_requests",
		})
	}))
	defer srv.Close()

	c := New(Config{APIKey: "test", BaseURL: srv.URL})
	err := c.Health(ctx)

	var rlErr *RateLimitError
	if !as(err, &rlErr) {
		t.Fatal("expected *RateLimitError")
	}
	if rlErr.Limit != 50 {
		t.Fatalf("expected limit 50, got %d", rlErr.Limit)
	}
	if rlErr.Remaining != 0 {
		t.Fatalf("expected remaining 0, got %d", rlErr.Remaining)
	}
}

func TestSendsIdempotencyKey(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("Idempotency-Key")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "OK",
			"data":    map[string]any{"message_id": "m1", "to": "+2547", "sender_name": "B", "parts": 1, "cost": 1.0, "currency": "KES", "status": "queued"},
		})
	}))
	defer srv.Close()

	c := New(Config{APIKey: "test", BaseURL: srv.URL})
	_, err := c.SMS().Send(ctx, SendSMSRequest{To: "+2547", Message: "Hi", From: "B"})
	if err != nil {
		t.Fatal(err)
	}
	if gotKey == "" {
		t.Fatal("expected Idempotency-Key header")
	}
}

func TestServiceMethods(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "req-123")
		w.WriteHeader(http.StatusOK)

		var body any
		switch r.URL.Path {
		case "/v1/sms/send":
			body = map[string]any{
				"success": true, "message": "SMS queued",
				"data": map[string]any{"message_id": "m1", "to": "+2547", "sender_name": "B", "parts": 1, "cost": 1.0, "currency": "KES", "status": "queued"},
			}
		case "/v1/email/send":
			body = map[string]any{
				"success": true, "message": "Email queued",
				"data": map[string]any{"message_id": "e1", "status": "queued"},
			}
		case "/v1/balance":
			body = map[string]any{
				"success": true, "message": "Balance retrieved",
				"data": map[string]any{"available_balance": 1250.50, "current_balance": 1500.00, "currency": "KES"},
			}
		default:
			body = map[string]any{
				"success": true, "message": "Message retrieved",
				"data": map[string]any{
					"id": "m1", "to": "+2547", "sender_name": "B", "parts": 1, "cost": 1.0,
					"status": "delivered", "created_at": "2024-01-01T00:00:00Z",
				},
			}
		}
		_ = json.NewEncoder(w).Encode(body)
	}))
	defer srv.Close()

	c := New(Config{APIKey: "test", BaseURL: srv.URL})

	t.Run("SMS.Send", func(t *testing.T) {
		resp, err := c.SMS().Send(ctx, SendSMSRequest{To: "+2547", Message: "Hi", From: "B"})
		if err != nil {
			t.Fatal(err)
		}
		if resp.MessageID != "m1" {
			t.Fatalf("expected m1, got %s", resp.MessageID)
		}
	})

	t.Run("Email.Send", func(t *testing.T) {
		resp, err := c.Email().Send(ctx, SendEmailRequest{To: "a@b.com", Subject: "Hi", Body: "Hello", From: "B"})
		if err != nil {
			t.Fatal(err)
		}
		if resp.MessageID != "e1" {
			t.Fatalf("expected e1, got %s", resp.MessageID)
		}
	})

	t.Run("Messages.Get", func(t *testing.T) {
		resp, err := c.Messages().Get(ctx, "m1")
		if err != nil {
			t.Fatal(err)
		}
		if resp.ID != "m1" {
			t.Fatalf("expected m1, got %s", resp.ID)
		}
	})

	t.Run("Account.GetBalance", func(t *testing.T) {
		resp, err := c.Account().GetBalance(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if resp.AvailableBalance != 1250.50 {
			t.Fatalf("expected 1250.50, got %f", resp.AvailableBalance)
		}
	})
}
