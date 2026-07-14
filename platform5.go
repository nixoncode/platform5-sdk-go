// Package platform5 provides a Go client for the Platform5 Developer API.
//
// Basic usage:
//
//	client := platform5.New(platform5.Config{
//	    APIKey: "p5_live_abc123",
//	})
//
//	resp, err := client.SMS().Send(ctx, platform5.SendSMSRequest{
//	    To:      "+254712345678",
//	    Message: "Hello from Platform5!",
//	    From:    "MyBrand",
//	})
//
// Each service is accessed via a method on [Client]:
//   - [Client.SMS] — send SMS messages
//   - [Client.Email] — send emails
//   - [Client.Messages] — check delivery status
//   - [Client.Account] — check balance
//   - [Client.Health] — health check
package platform5
