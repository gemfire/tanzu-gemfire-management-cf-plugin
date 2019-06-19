package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCloudcacheManagementCfPlugin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CloudcacheManagementCfPlugin Suite")
}
