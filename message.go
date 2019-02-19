package email

import (
	"errors"
	"fmt"
	"io"
	"net/mail"
	"strings"
)

type Message struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string

	PlaintextBody string
	HTMLBody      string
}

func (msg *Message) RecipientAddresses() (addrs []string, err error) {
	for _, addr := range append(append(msg.To, msg.Cc...), msg.Bcc...) {
		parsed, err := mail.ParseAddress(addr)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, parsed.Address)
	}
	return addrs, nil
}

func (msg *Message) String() string {
	var b strings.Builder
	err := msg.Encode(&b)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

// func (msg *Message) Bytes() []byte {
// 	var b bytes.Buffer
// 	err := msg.Encode(&b)
// 	if err != nil {
// 		return nil
// 	}
// 	return b.Bytes()
// }

func (msg *Message) Encode(w io.Writer) error {
	err := msg.encodeAddressesAndSubject(w)
	if err != nil {
		return err
	}
	return msg.encodeBody(w)
}

func (msg *Message) encodeAddressesAndSubject(w io.Writer) error {
	if len(msg.To) == 0 {
		return errors.New("email.Message has no 'To' receiver")
	}
	if msg.Subject == "" {
		return errors.New("email.Message has no subject")
	}

	if msg.From != "" {
		err := msg.encodeAddresses(w, "From", msg.From)
		if err != nil {
			return err
		}
	}
	err := msg.encodeAddresses(w, "To", msg.To...)
	if err != nil {
		return err
	}
	err = msg.encodeAddresses(w, "Cc", msg.Cc...)
	if err != nil {
		return err
	}
	// err = msg.encodeAddresses(w, "Bcc", msg.Bcc...)
	// if err != nil {
	// 	return err
	// }

	_, err = fmt.Fprintf(w, "Subject: %s\r\n", encodeRFC2047(msg.Subject))
	return err
}

func (msg *Message) encodeBody(w io.Writer) error {
	// todo
	fmt.Fprint(w, "MIME-Version: 1.0\r\n")
	fmt.Fprint(w, "Content-Type: text/html; charset=\"utf-8\"\r\n")
	fmt.Fprint(w, "Content-Transfer-Encoding: base64\r\n")
	fmt.Fprint(w, "\r\n")
	fmt.Fprint(w, msg.PlaintextBody)
	return nil
}

func (msg *Message) encodeAddresses(w io.Writer, kind string, addresses ...string) error {
	if len(addresses) == 0 {
		return nil
	}
	for i := range addresses {
		addr, err := mail.ParseAddress(addresses[i])
		if err != nil {
			return err
		}
		if i == 0 {
			_, err = fmt.Fprintf(w, "%s: %s", kind, addr)
		} else {
			_, err = fmt.Fprintf(w, ", %s", addr)
		}
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w, "\r\n")
	return err
}

func encodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Name: str}
	return strings.Trim(strings.TrimSuffix(addr.String(), " <@>"), `"`)
}
