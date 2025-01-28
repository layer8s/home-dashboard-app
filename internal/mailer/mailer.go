package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Declare a new variable with the type embed.FS (embedded file system) to hold email templates.

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	client *sendgrid.Client
	sender string
}

func New(apiKey, senderEmail string) *Mailer {
	client := sendgrid.NewSendClient(apiKey)
	return &Mailer{
		client: client,
		sender: senderEmail,
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
	response, err := m.client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}
