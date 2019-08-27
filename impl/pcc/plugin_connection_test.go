package pcc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	. "github.com/gemfire/cloudcache-management-cf-plugin/impl/pcc"
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
)

var _ = Describe("PluginConnection", func() {

	var (
		cliConnection    *pluginfakes.FakeCliConnection
		pluginConnection impl.ConnectionProvider
	)

	BeforeEach(func() {
		cliConnection = new(pluginfakes.FakeCliConnection)
		pluginConnectionImpl, err := NewPluginConnectionProvider(cliConnection)
		pluginConnection = pluginConnectionImpl
		Expect(err).NotTo(HaveOccurred())
	})

	Context("We have a service and a service-key", func() {

		It("Returns a populated ConnectionData object", func() {
			cliConnection.CliCommandWithoutTerminalOutputReturnsOnCall(0, []string{"name", "someKey"}, nil)
			cliConnection.CliCommandWithoutTerminalOutputReturnsOnCall(1, []string{"", ""}, nil)
			connectionData, err := pluginConnection.GetConnectionData("pcc1")
			Expect(err).NotTo(HaveOccurred())
			Expect(cliConnection.CliCommandWithoutTerminalOutputCallCount).To(Equal(2))
			Expect(connectionData.Username).To(Equal("cluster_operator"))
			Expect(connectionData.Password).To(Equal("password"))
			Expect(connectionData.LocatorAddress).To(Equal("http://locator.domain.com"))
		})
	})

	Context("We don't have a service-key", func() {

		It("Returns an error indicating that there is no service-key", func() {
			cliConnection.CliCommandWithoutTerminalOutputReturnsOnCall(0, []string{"", ""}, nil)
			cliConnection.CliCommandWithoutTerminalOutputReturnsOnCall(1, []string{"", ""}, nil)
			connectionData, err := pluginConnection.GetConnectionData("pcc1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(util.NoServiceKeyMessage))
			Expect(cliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(len(connectionData.Username)).To(BeZero())
			Expect(len(connectionData.Password)).To(BeZero())
			Expect(len(connectionData.LocatorAddress)).To(BeZero())
		})
	})
})
