package geode_test

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	. "github.com/gemfire/cloudcache-management-cf-plugin/impl/geode"
)

var _ = Describe("GeodeConnection", func() {

	var (
		geodeConnection impl.ConnectionProvider
		commandData     domain.CommandData
	)

	BeforeEach(func() {
		geodeConnectionImpl, err := NewGeodeConnectionProvider()
		geodeConnection = geodeConnectionImpl
		Expect(err).NotTo(HaveOccurred())

		commandData = domain.CommandData{}
		commandData.UserCommand = domain.UserCommand{}
		commandData.UserCommand.Parameters = make(map[string]string)
	})

	Context("All co-ordinates are provided on the command line", func() {

		BeforeEach(func() {
			commandData.Target = "https://some.geode-locator.com"
			commandData.UserCommand.Parameters["-u"] = "locatorUser"
			commandData.UserCommand.Parameters["-p"] = "locatorPassword"
		})

		It("Returns a populated ConnectionData object", func() {
			err := geodeConnection.GetConnectionData(&commandData)
			Expect(err).NotTo(HaveOccurred())
			Expect(commandData.ConnnectionData.Username).To(Equal("locatorUser"))
			Expect(commandData.ConnnectionData.Password).To(Equal("locatorPassword"))
			Expect(commandData.ConnnectionData.LocatorAddress).To(Equal("https://some.geode-locator.com"))
		})
	})
})
