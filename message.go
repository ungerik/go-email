package email

import (
	"crypto/rand"
	"encoding/base64"
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

	Plaintext string
	HTML      []byte
}

// RecipientAddresses returns the address parts of the To, Cc, Bcc slices
func (msg *Message) RecipientAddresses() (rcpt []string, err error) {
	for _, addr := range append(append(msg.To, msg.Cc...), msg.Bcc...) {
		parsed, err := mail.ParseAddress(addr)
		if err != nil {
			return nil, err
		}
		rcpt = append(rcpt, parsed.Address)
	}
	return rcpt, nil
}

// String returns the encoded message or any encoding error as string.
func (msg *Message) String() string {
	var b strings.Builder
	err := msg.Encode(&b, true)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

func (msg *Message) Encode(w io.Writer, withBcc bool) error {
	// https://th-h.de/net/usenet/faqs/headerfaq/

	if len(msg.To) == 0 {
		return errors.New("email.Message has no 'To' receiver")
	}
	if msg.Subject == "" {
		return errors.New("email.Message has no subject")
	}

	// Header

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
	if withBcc {
		err = msg.encodeAddresses(w, "Bcc", msg.Bcc...)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(w, "Subject: %s\r\n", encodeRFC2047(msg.Subject))
	if err != nil {
		return err
	}

	// Body

	switch {
	case msg.Plaintext != "" && len(msg.HTML) > 0:
		// HTML and Plaintext
		boundary := randomString()
		_, err = fmt.Fprintf(w, "MIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n", boundary)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "\r\n--%s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", boundary, msg.Plaintext)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "\r\n--%s\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: base64\r\n\r\n%s\r\n--%s--", boundary, base64.RawURLEncoding.EncodeToString(msg.HTML), boundary)

	case msg.Plaintext == "" && len(msg.HTML) > 0:
		// HTML
		_, err = fmt.Fprintf(w, "MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", msg.HTML)
	default:
		// Plaintext
		_, err = fmt.Fprintf(w, "MIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", msg.Plaintext)
	}

	return err
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

// randomString returns a 120 bit randum number
// encoded as URL compatible base64 string with a length of 20 characters.
func randomString() string {
	var buffer [15]byte
	b := buffer[:]
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
