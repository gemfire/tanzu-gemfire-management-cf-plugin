package geode_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGeode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Geode Suite")
}
