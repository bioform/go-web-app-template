package server_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/bioform/go-web-app-template/config"
	"github.com/bioform/go-web-app-template/pkg/logging"
	"github.com/bioform/go-web-app-template/pkg/mail"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	logging.InitLogger()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var _ = BeforeSuite(func() {
	smtp := config.App.Email.Smtp
	// You can pass empty smtpmock.ConfigurationAttr{}. It means that smtpmock will use default settings
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		PortNumber:        smtp.Port,
		LogToStdout:       true,
		LogServerActivity: true,
	})

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start smtpmock server: %v", err)
	}

	// Server's port will be assigned dynamically after server.Start()
	// for case when portNumber wasn't specified
	hostAddress, portNumber := "127.0.0.1", server.PortNumber()

	// Possible SMTP-client stuff for iteration with mock server
	address := fmt.Sprintf("%s:%d", hostAddress, portNumber)
	log.Printf("SMTP server started on %s", address)

	DeferCleanup(server.Stop)
	DeferCleanup(mail.Client().Close)
})
