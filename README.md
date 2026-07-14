# platform5-sdk-go

Go SDK for the Platform5 Developer API.

## Install

```sh
go get github.com/nixoncode/platform5-sdk-go
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/nixoncode/platform5-sdk-go"
)

func main() {
    ctx := context.Background()
    client := platform5.New(platform5.Config{
        APIKey:  "p5_live_abc123def456",
        BaseURL: "http://localhost:8084",
    })

    // Send an SMS
    sms, err := client.SMS().Send(ctx, platform5.SendSMSRequest{
        To:      "+254712345678",
        Message: "Your appointment is confirmed for tomorrow at 10 AM.",
        From:    "MyBrand",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Sent: %s (%d parts)\n", sms.MessageID, sms.Parts)

    // Send an email
    client.Email().Send(ctx, platform5.SendEmailRequest{
        To:       "user@example.com",
        Subject:  "Welcome",
        Body:     "<h1>Hello!</h1>",
        BodyType: "html",
        From:     "MyBrand",
    })

    // Check message status
    status, _ := client.Messages().Get(ctx, sms.MessageID)
    fmt.Printf("Status: %s\n", status.Status)

    // Check balance
    balance, _ := client.Account().GetBalance(ctx)
    fmt.Printf("Balance: %.2f %s\n", balance.AvailableBalance, balance.Currency)
}
```

## Configuration

```go
platform5.New(platform5.Config{
    APIKey:  "p5_live_abc123def456",   // required
    BaseURL: "http://localhost:8084",  // optional, defaults to http://localhost:8084
})
```

## Services

| Service | Method | Endpoint |
|---------|--------|----------|
| `client.SMS().Send(ctx, req)` | POST | `/v1/sms/send` |
| `client.Email().Send(ctx, req)` | POST | `/v1/email/send` |
| `client.Messages().Get(ctx, id)` | GET | `/v1/messages/{id}` |
| `client.Account().GetBalance(ctx)` | GET | `/v1/balance` |
| `client.Health(ctx)` | GET | `/health` |

## Error Handling

```go
resp, err := client.SMS().Send(ctx, req)
if err != nil {
    var p5err *platform5.Error
    if errors.As(err, &p5err) {
        switch {
        case p5err.IsUnauthorized():
            // 401
        case p5err.IsInsufficientBalance():
            // 402
        case p5err.IsForbidden():
            // 403
        case p5err.IsNotFound():
            // 404
        case p5err.IsValidationFailed():
            // 422
        }
    }

    var rlErr *platform5.RateLimitError
    if errors.As(err, &rlErr) {
        fmt.Printf("Rate limited: %d/%d\n", rlErr.Remaining, rlErr.Limit)
    }
}
```

## Idempotency

`SMS.Send()` and `Email.Send()` automatically generate a UUID v4 `Idempotency-Key` header for every request to prevent duplicate charges on retry.

## Requirements

- Go 1.25+
