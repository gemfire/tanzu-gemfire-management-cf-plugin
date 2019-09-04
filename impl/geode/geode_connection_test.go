package geode_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	. "github.com/gemfire/cloudcache-management-cf-plugin/impl/geode"
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
)

var _ = Describe("GeodeConnection", func() {

	var (
		geodeConnection impl.ConnectionProvider
		args            []string
		locatorURL      string
		userName        string
		password        string
	)

	BeforeEach(func() {
		geodeConnectionImpl, err := NewGeodeConnectionProvider()
		geodeConnection = geodeConnectionImpl
		Expect(err).NotTo(HaveOccurred())
	})

	Context("All co-ordinates are provided on the command line", func() {

		BeforeEach(func() {
			locatorURL = "https://some.geode-locator.com"
			userName = "-u=locatorUser"
			password = "-p=locatorPassword"
			args = []string{locatorURL, userName, password}
		})

		It("Returns a populated ConnectionData object", func() {
			connectionData, err := geodeConnection.GetConnectionData(args)
			Expect(err).NotTo(HaveOccurred())
			Expect(connectionData.Username).To(Equal("locatorUser"))
			Expect(connectionData.Password).To(Equal("locatorPassword"))
			Expect(connectionData.LocatorAddress).To(Equal("https://some.geode-locator.com"))
		})
	})

	Context("URI is missing from arguments", func() {

		BeforeEach(func() {
			locatorURL = "some.geode-locator.com"
			args = []string{locatorURL}
		})

		It("Returns an error indicating that URL is missing", func() {
			connectData, err := geodeConnection.GetConnectionData(args)
			Expect(err).To(HaveOccurred())
			Expect(len(connectData.LocatorAddress)).To(BeZero())
			Expect(err.Error()).To(Equal(util.NoEndpointFoundMessage))
		})
	})
})
