package email

import (
	"context"

	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/AboloreDev/geritcht-restaurant/internals/templates"
)

func (c *ResendEmailClient) SendVerificationMail(data events.VerificationEmailPayload) error {
	ctx := context.Background()
	email := &EmailRequest{
		Body:    templates.VerificationEmailTemplate(data.Token),
		To:      data.Email,
		Subject: "Email Verification",
	}

	return c.SendEmail(ctx, email)
}

func (c *ResendEmailClient) SendPasswordResetMail(data events.PasswordResetEmailPayload) error {
	ctx := context.Background()
	email := &EmailRequest{
		Body:    templates.PasswordResetTemplate(data.Token),
		To:      data.Email,
		Subject: "Password Reset",
	}

	return c.SendEmail(ctx, email)
}

func (c *ResendEmailClient) SendPasswordChangedMail(data events.PasswordChangedEmailPayload) error {
	ctx := context.Background()
	email := &EmailRequest{
		Body:    templates.PasswordChangedTemplate(data.FirstName),
		To:      data.Email,
		Subject: "Password Changed",
	}

	return c.SendEmail(ctx, email)
}
