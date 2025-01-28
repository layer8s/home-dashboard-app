package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Declare a new variable with the type embed.FS (embedded file system) to hold email templates.

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	client *sendgrid.Client
	sender string
	logger *slog.Logger
}

func New(apiKey, senderEmail string, logger *slog.Logger) *Mailer {
	client := sendgrid.NewSendClient(apiKey)
	return &Mailer{
		client: client,
		sender: senderEmail,
		logger: logger,
	}
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Follow the same pattern to execute the "plainBody" template and store the result
	// in the plainBody variable.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// And likewise with the "htmlBody" template.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}
	from := mail.NewEmail("Fantasy Football Archive", m.sender)
	to := mail.NewEmail("Example User", recipient)
	plainTextContent := plainBody.String()
	htmlContent := htmlBody.String()
	message := mail.NewSingleEmail(from, subject.String(), to, plainTextContent, htmlContent)

	for i := 1; i <= 3; i++ {
		response, err := m.client.Send(message)
		if nil == err {
			// If the email is sent successfully, break out of the loop.
			m.logger.Info("email sent successfully",
				"status", response.StatusCode,
				"recipient", recipient,
				"Body", response.Body,
				"Headers", response.Headers)
			return nil

		}

		m.logger.Error("failed to send email",
			"error", err,
			"recipient", recipient,
			"template", templateFile,
			"Sending attempt", i,
			"Retrying in : ", 500*time.Millisecond)
		time.Sleep(500 * time.Millisecond)

	}
	return fmt.Errorf("error sending email: %w", err)

}
