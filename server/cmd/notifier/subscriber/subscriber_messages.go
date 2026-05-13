package subscriber

import (
	"encoding/json"
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
