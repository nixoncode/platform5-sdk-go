package platform5

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// airtimeService provides access to the Airtime API.
type airtimeService struct {
	client *Client
}

// Buy purchases airtime for a phone number.
//
// An Idempotency-Key header is automatically generated to prevent
// duplicate charges on retry.
func (s *airtimeService) Buy(ctx context.Context, req BuyAirtimeRequest) (*BuyAirtimeResponse, error) {
	key := uuid.New().String()
	raw, err := s.client.request(ctx, http.MethodPost, "/v1/airtime/buy", req, key)
	if err != nil {
		return nil, err
	}
	return decode[*BuyAirtimeResponse](raw)
}
