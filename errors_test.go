package platform5

import "testing"

func TestErrorHelpers(t *testing.T) {
	tests := []struct {
		err    *Error
		check  func(*Error) bool
		expect bool
	}{
		{&Error{StatusCode: 401}, (*Error).IsUnauthorized, true},
		{&Error{StatusCode: 402}, (*Error).IsInsufficientBalance, true},
		{&Error{StatusCode: 403}, (*Error).IsForbidden, true},
		{&Error{StatusCode: 404}, (*Error).IsNotFound, true},
		{&Error{StatusCode: 422}, (*Error).IsValidationFailed, true},
		{&Error{StatusCode: 400}, (*Error).IsUnauthorized, false},
	}

	for _, tt := range tests {
		got := tt.check(tt.err)
		if got != tt.expect {
			t.Errorf("status %d: got %v, want %v", tt.err.StatusCode, got, tt.expect)
		}
	}
}

func TestRateLimitErrorFields(t *testing.T) {
	err := &RateLimitError{
		Err:       Error{StatusCode: 429, Message: "too fast"},
		Limit:     50,
		Remaining: 3,
	}

	if err.Err.StatusCode != 429 {
		t.Fatalf("expected 429, got %d", err.Err.StatusCode)
	}
	if err.Limit != 50 {
		t.Fatalf("expected limit 50, got %d", err.Limit)
	}
	if err.Remaining != 3 {
		t.Fatalf("expected remaining 3, got %d", err.Remaining)
	}

	msg := err.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
}
