package config_test

import (
	log "github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServices(t *testing.T) {
	log.SetLevel(log.PanicLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}
