package subscriber

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AboloreDev/geritcht-restaurant/internals/email"
	"github.com/AboloreDev/geritcht-restaurant/internals/events"
	"github.com/ThreeDotsLabs/watermill/message"
)

func (s *EventSubscriber) HandleSendVerificationMail(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.VerificationEmailPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendVerificationMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleSendPasswordReset(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.PasswordResetEmailPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendPasswordResetMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleSendPasswordChangedMail(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.PasswordChangedEmailPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendPasswordChangedMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleReservationConfirmationMail(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.ReservationConfirmPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)
	fmt.Println(data.Email)

	emailClient.SendReservationConfirmation(data)

	log.Println("Email successfully sent")

	return nil

}

func (s *EventSubscriber) HandleReservationReminderMail(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.ReservationReminderPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendReservationReminderMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleReservationCancellationMail(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.ReservationCancelledPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendReservationCancellationMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleReservationCheckInMail(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.ReservationCheckInPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendReservationCheckInMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleReservationNoShowMail(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.ReservationNoShowPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendReservationNoShowMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleOrderConfirmationMail(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.OrderConfirmationPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendOrderConfirmationMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleOrderRefundPayload(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.OrderRefundedPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendOrderRefundMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleLowStockAlert(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.LowStockAlertPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.AdminEmail)

	emailClient.SendLowStockAlertMail(data)

	log.Println("Email successfully sent")

	return nil
}

func (s *EventSubscriber) HandleWaitlistNotifier(msg *message.Message, emailClient *email.ResendEmailClient) error {
	var data events.WaitlistNotificationPayload

	err := json.Unmarshal(msg.Payload, &data)
	if err != nil {
		return err
	}

	log.Printf("Sending notification to %s", data.Email)

	emailClient.SendWaitlistNotificationMail(data)

	log.Println("Email successfully sent")

	return nil
}
