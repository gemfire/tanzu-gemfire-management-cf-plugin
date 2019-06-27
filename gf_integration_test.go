package main

import (
	//"code.cloudfoundry.org/cli/cf/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//. "github.com/onsi/gomega"
	"github.com/gemfire/cloudcache-management-cf-plugin/cfservice")
var _ = Describe("cf gf plugin integration tests", func() {
	BeforeEach(func(){
		cfserv := cfservice.Cf{}
		_,err := cfserv.Cmd("create-service", "p-cloudecache", "dev-plan", "serviceForTesting")
		Expect(err).To(BeNil())
		_,err = cfserv.Cmd("create-service-key","serviceForTesting", "key")
		Expect(err).To(BeNil())

	})
})
