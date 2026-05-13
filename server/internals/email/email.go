package email

import (
	"context"
	"fmt"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/resend/resend-go/v3"
)

type ResendEmailClient struct {
	client *resend.Client
	cfg    *config.ResendConfig
}

type EmailRequest struct {
	Body    string
	Subject string
	To      string
}

func NewResendEmailClient(ctx context.Context, cfg *config.ResendConfig) *ResendEmailClient {
	client := resend.NewClient(cfg.ResendAPIKey)

	return &ResendEmailClient{
		client: client,
		cfg:    cfg,
	}
}

func (c *ResendEmailClient) SendEmail(ctx context.Context, email *EmailRequest) error {
	params := &resend.SendEmailRequest{
		From:    c.cfg.ResendFromMail,
		To:      []string{email.To},
		Subject: email.Subject,
		Html:    email.Body,
	}

	sent, err := c.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("Email sent successfully %s", sent.Id)

	return nil
}
