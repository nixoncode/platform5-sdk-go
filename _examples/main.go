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

	// Health check
	if err := client.Health(ctx); err != nil {
		log.Fatal(err)
	}

	// Send SMS
	sms, err := client.SMS().Send(ctx, platform5.SendSMSRequest{
		To:      "+254712345678",
		Message: "Your appointment is confirmed for tomorrow at 10 AM.",
		From:    "MyBrand",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SMS sent: %s (%d parts, %.2f %s)\n", sms.MessageID, sms.Parts, sms.Cost, sms.Currency)

	// Send email
	email, err := client.Email().Send(ctx, platform5.SendEmailRequest{
		To:       "user@example.com",
		Subject:  "Welcome to Platform5",
		Body:     "<h1>Hello!</h1><p>Thank you for signing up.</p>",
		BodyType: "html",
		From:     "MyBrand",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Email sent: %s\n", email.MessageID)

	// Check balance
	balance, err := client.Account().GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %.2f %s\n", balance.AvailableBalance, balance.Currency)

	// Get message status
	status, err := client.Messages().Get(ctx, sms.MessageID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Message status: %s\n", status.Status)
}
