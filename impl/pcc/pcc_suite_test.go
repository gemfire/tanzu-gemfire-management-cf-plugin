package pcc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPcc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pcc Suite")
}
