package email

import (
	"bytes"
	"context"
	"embed"
	"html/template"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/wneessen/go-mail"
)

//go:embed templates/*
var emailTemplate embed.FS

// SendConfirmationEmail sends a confirmation email to the newly registered user
func SendConfirmationEmail(ctx context.Context, user model.User) error {
	// Parse the email template
	tmpl, err := template.ParseFS(emailTemplate, "templates/email_confirmation.html")
	if err != nil {
		return err
	}

	// Prepare the data for the template
	data := struct {
		Name             string
		ConfirmationLink string
	}{
		Name:             user.Name,
		ConfirmationLink: "https://example.com/confirm_email",
	}

	// Execute the template with user data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	// Create a new email message
	m := mail.NewMsg()

	// Set the sender address
	if err := m.From("your_email@example.com"); err != nil {
		return err
	}

	// Add recipient
	if err := m.To(user.Email); err != nil {
		return err
	}

	// Set the email subject
	m.Subject("Email Confirmation")

	// Set the email body
	m.SetBodyString(mail.TypeTextHTML, body.String())

	// Create a new SMTP client
	client, err := mail.NewClient("smtp.example.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername("your_email@example.com"),
		mail.WithPassword("your_password"),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return err
	}

	// Send the email
	if err := client.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
