package email

import (
	"bytes"
	"encoding/base64"

	gmail "google.golang.org/api/gmail/v1"
)

func NewGMailMessage(msg *Message) (*gmail.Message, error) {
	var b bytes.Buffer
	err := msg.Encode(&b)
	if err != nil {
		return nil, err
	}
	return &gmail.Message{Raw: base64.RawURLEncoding.EncodeToString(b.Bytes())}, nil
}

func SendViaGMail(msg *Message, service *gmail.Service) error {
	gmailMessage, err := NewGMailMessage(msg)
	if err != nil {
		return err
	}
	_, err = service.Users.Messages.Send("me", gmailMessage).Do()
	return err
}
