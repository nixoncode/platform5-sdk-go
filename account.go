package platform5

import (
	"context"
	"net/http"
)

// accountService provides access to the Account API.
type accountService struct {
	client *Client
}

// GetBalance retrieves the organisation's current available balance.
func (s *accountService) GetBalance(ctx context.Context) (*BalanceResponse, error) {
	raw, err := s.client.request(ctx, http.MethodGet, "/v1/balance", nil, "")
	if err != nil {
		return nil, err
	}
	return decode[*BalanceResponse](raw)
}
