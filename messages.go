package platform5

import (
	"context"
	"net/http"
)

// messagesService provides access to the Message Status API.
type messagesService struct {
	client *Client
}

// Get retrieves the current delivery status of a message by its ID.
func (s *messagesService) Get(ctx context.Context, id string) (*MessageStatusResponse, error) {
	raw, err := s.client.request(ctx, http.MethodGet, "/v1/messages/"+id, nil, "")
	if err != nil {
		return nil, err
	}
	return decode[*MessageStatusResponse](raw)
}
