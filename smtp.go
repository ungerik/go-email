package email

import (
	"bytes"
	"fmt"
	"net/smtp"
)

func SendViaSMTP(msg *Message, host string, port uint16, username, password string) error {
	to, err := msg.RecipientAddresses()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = msg.Encode(&b, false)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", username, password, host)
	return smtp.SendMail(addr, auth, msg.From, to, b.Bytes())
}
