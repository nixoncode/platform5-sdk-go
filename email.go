package platform5

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// emailService provides access to the Email API.
type emailService struct {
	client *Client
}

// Send queues an email for delivery.
//
// An Idempotency-Key header is automatically generated to prevent
// duplicate sends on retry.
func (s *emailService) Send(ctx context.Context, req SendEmailRequest) (*SendEmailResponse, error) {
	key := uuid.New().String()
	raw, err := s.client.request(ctx, http.MethodPost, "/v1/email/send", req, key)
	if err != nil {
		return nil, err
	}
	return decode[*SendEmailResponse](raw)
}
