package platform5

// Config configures the Platform5 API client.
type Config struct {
	// APIKey is your Platform5 Developer API key (starts with p5_).
	APIKey string

	// BaseURL is the Platform5 API base URL.
	// Defaults to https://api.advantaplus.ai.
	BaseURL string
}

// apiResponse is the standard JSON envelope returned by all Platform5 API endpoints.
type apiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Errors  any    `json:"errors"`
}

// SendSMSRequest is the request body for POST /v1/sms/send.
type SendSMSRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
	From    string `json:"from"`
}

// SendSMSResponse is returned on a successful SMS send.
type SendSMSResponse struct {
	MessageID  string  `json:"message_id"`
	To         string  `json:"to"`
	SenderName string  `json:"sender_name"`
	Parts      int     `json:"parts"`
	Cost       float64 `json:"cost"`
	Currency   string  `json:"currency"`
	Status     string  `json:"status"`
}

// SendEmailRequest is the request body for POST /v1/email/send.
type SendEmailRequest struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	BodyType string `json:"body_type,omitempty"`
	From     string `json:"from"`
}

// SendEmailResponse is returned on a successful email send.
type SendEmailResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// MessageStatusResponse contains the current delivery status of a message.
type MessageStatusResponse struct {
	ID          string  `json:"id"`
	To          string  `json:"to"`
	SenderName  string  `json:"sender_name"`
	Parts       int     `json:"parts"`
	Cost        float64 `json:"cost"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	SentAt      *string `json:"sent_at,omitempty"`
	DeliveredAt *string `json:"delivered_at,omitempty"`
	Error       *string `json:"error,omitempty"`
}

// BalanceResponse contains the organisation's current and available balance.
type BalanceResponse struct {
	AvailableBalance float64 `json:"available_balance"`
	CurrentBalance   float64 `json:"current_balance"`
	Currency         string  `json:"currency"`
}
