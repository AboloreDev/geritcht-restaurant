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

func (c *ResendEmailClient) SendReservationConfirmation(data events.ReservationConfirmPayload) error {
	ctx := context.Background()
	email := &EmailRequest{
		Body: templates.ReservationConfirmationTemplate(
			data.FirstName,
			data.Date,
			data.TimeSlot,
			data.TableName,
		),
		To:      data.Email,
		Subject: "Reservation Confirmation Mail",
	}

	return c.SendEmail(ctx, email)
}
func (c *ResendEmailClient) SendReservationCheckInMail(data events.ReservationCheckInPayload) error {
	ctx := context.Background()
	email := &EmailRequest{
		Body: templates.ReservationCheckInTemplate(
			data.FirstName,
			data.Date,
			data.TimeSlot,
			data.TableName,
			data.PartySize,
		),
		To:      data.Email,
		Subject: "Reservation Check-In Mail",
	}

	return c.SendEmail(ctx, email)
}

func (c *ResendEmailClient) SendReservationCancellationMail(data events.ReservationCancelledPayload) error {
	ctx := context.Background()
	email := &EmailRequest{
		Body: templates.ReservationCancellationTemplate(
			data.FirstName,
			data.Date,
			data.TimeSlot,
			data.TableName,
			data.PartySize,
		),
		To:      data.Email,
		Subject: "Reservation Cancellation Mail",
	}

	return c.SendEmail(ctx, email)
}

func (c *ResendEmailClient) SendReservationReminderMail(data events.ReservationReminderPayload) error {
	ctx := context.Background()
	email := &EmailRequest{
		Body: templates.ReservationReminderTemplate(
			data.FirstName,
			data.Date,
			data.TimeSlot,
			data.TableName,
			data.PartySize,
		),
		To:      data.Email,
		Subject: "Reservation Reminder Mail",
	}

	return c.SendEmail(ctx, email)
}

func (c *ResendEmailClient) SendReservationNoShowMail(data events.ReservationNoShowPayload) error {
	ctx := context.Background()
	email := &EmailRequest{
		Body: templates.ReservationNoShowTemplate(
			data.FirstName,
			data.Date,
			data.TimeSlot,
			data.TableName,
			data.PartySize,
		),
		To:      data.Email,
		Subject: "Rservation No-Show Mail",
	}

	return c.SendEmail(ctx, email)
}
