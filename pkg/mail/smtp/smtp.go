package smtp

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
)

// TODO: Implement new error handling
func NewMailer() (rolf.Mailer, error) {
	// Connect to the remote SMTP server.
	m, err := smtp.Dial(rolf.Env["SMTP_HOST"] + ":" + rolf.Env["SMTP_PORT"])
	if err != nil {
		log.Fatal(err)
		return rolf.Mailer(Mail{}), fmt.Errorf("Can't connect to mailer\n %s", err)
	}

	return rolf.Mailer(Mail{m}), nil
}

type Mail struct {
	*smtp.Client
}

// from sender to recipient with email body
func (m Mail) Send(s string, r string, b string) error {
	if err := m.Mail(s); err != nil {
		log.Fatal(err)
	}
	if err := m.Rcpt(r); err != nil {
		log.Fatal(err)
	}

	// Send the email body.
	wc, err := m.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, b)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Send the QUIT command and close the connection.
	err = m.Quit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
