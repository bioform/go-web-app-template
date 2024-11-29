package server_test

import (
	"testing"

	"github.com/bioform/go-web-app-template/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	logging.InitLogger()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}
