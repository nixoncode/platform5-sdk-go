package platform5

import (
	"context"
	"errors"
	"testing"
)

var ctx = context.Background()

func as[E error](err error, target *E) bool {
	return errors.As(err, target)
}

func TestNewDefaultBaseURL(t *testing.T) {
	c := New(Config{APIKey: "test"})
	if c == nil {
		t.Fatal("client is nil")
	}
}

func TestServiceLazyInit(t *testing.T) {
	c := New(Config{APIKey: "test"})
	if c.SMS() == nil || c.Email() == nil || c.Messages() == nil || c.Account() == nil {
		t.Fatal("service should not be nil")
	}
}
