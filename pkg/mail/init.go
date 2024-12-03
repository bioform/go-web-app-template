package mail

import (
	"fmt"

	"github.com/bioform/go-web-app-template/config"
	"github.com/bioform/go-web-app-template/pkg/env"
	"github.com/wneessen/go-mail"
)

var client *mail.Client

func Client() *mail.Client {
	return client
}

func init() {
	var err error
	cfg := config.App.Email.Smtp

	client, err = mail.NewClient(cfg.Host, options()...)
	if err != nil {
		panic(fmt.Errorf("failed to create mail client: %w", err))
	}
}

func options() []mail.Option {
	cfg := config.App.Email.Smtp
	opts := []mail.Option{
		mail.WithPort(cfg.Port),
		mail.WithSMTPAuth(cfg.SMTPAuth()),
		mail.WithUsername(cfg.Username),
		mail.WithPassword(cfg.Password),
		mail.WithTLSPolicy(cfg.TlsPolicy()),
	}

	if env.IsTest() {
		opts = append(opts, mail.WithHELO("localhost"))
	}

	return opts
}
