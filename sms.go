package platform5

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// smsService provides access to the SMS API.
type smsService struct {
	client *Client
}

// Send queues an SMS message for delivery.
//
// An Idempotency-Key header is automatically generated to prevent
// duplicate charges on retry.
func (s *smsService) Send(ctx context.Context, req SendSMSRequest) (*SendSMSResponse, error) {
	key := uuid.New().String()
	raw, err := s.client.request(ctx, http.MethodPost, "/v1/sms/send", req, key)
	if err != nil {
		return nil, err
	}
	return decode[*SendSMSResponse](raw)
}
