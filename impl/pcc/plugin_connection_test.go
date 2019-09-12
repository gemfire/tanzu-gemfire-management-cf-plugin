package pcc_test

import (
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl"
	. "github.com/gemfire/cloudcache-management-cf-plugin/impl/pcc"
)

var _ = Describe("PluginConnection", func() {

	var (
		cliConnection          *pluginfakes.FakeCliConnection
		pluginConnection       impl.ConnectionProvider
		goodServiceKeyResponse []string
		commandData            domain.CommandData
	)

	BeforeEach(func() {
		cliConnection = new(pluginfakes.FakeCliConnection)
		pluginConnectionImpl, err := NewPluginConnectionProvider(cliConnection)
		pluginConnection = pluginConnectionImpl
		Expect(err).NotTo(HaveOccurred())

		commandData = domain.CommandData{}
		commandData.UserCommand = domain.UserCommand{}
		commandData.UserCommand.Parameters = make(map[string]string)
		commandData.Target = "pcc1"
	})

	Context("We have a service and a service-key", func() {
		BeforeEach(func() {
			goodServiceKeyResponseString := `Getting key pcc1ServiceKey for service instance pcc1 as admin...
				
{
 "distributed_system_id": "0",
 "gfsh_login_string": "connect --url=https://cloudcache-45371efd-f4ca-4549-a5f2-e06330aa53dc.sys.riverbank.cf-app.com/gemfire/v1 --user=cluster_operator_M5Scgeb0b6yp5f99E6SA8w --password=AMmxU9H6J5KSCYDLccipIw --skip-ssl-validation",
 "locators": [
  "10.0.8.6[55221]",
  "10.0.8.5[55221]",
  "10.0.8.7[55221]"
 ],
 "urls": {
  "gfsh": "https://cloudcache-45371efd-f4ca-4549-a5f2-e06330aa53dc.sys.riverbank.cf-app.com/gemfire/v1",
  "pulse": "https://cloudcache-45371efd-f4ca-4549-a5f2-e06330aa53dc.sys.riverbank.cf-app.com/pulse"
 },
 "users": [
  {
   "password": "AMmxU9H6J5KSCYDLccipIw",
   "roles": [
    "cluster_operator"
   ],
   "username": "cluster_operator_M5Scgeb0b6yp5f99E6SA8w"
  },
  {
   "password": "PVTrpLgX7K53rthdvd67CQ",
   "roles": [
    "developer"
   ],
   "username": "developer_3VTqJTIkftQX3pBJcSW1w"
  }
 ],
 "wan": {
  "sender_credentials": {
   "active": {
    "password": "nPYPcJBLoI4vbFFODV7ULg",
    "username": "gateway_sender_qZDiVscHzRqF34RFsWTrVQ"
   }
  }
 }
}`
			goodServiceKeyResponse = strings.Split(goodServiceKeyResponseString, "\n")
		})

		It("Returns a populated ConnectionData object", func() {
			cliConnection.CliCommandWithoutTerminalOutputReturnsOnCall(0, []string{"name", "pcc1ServiceKey"}, nil)
			cliConnection.CliCommandWithoutTerminalOutputReturnsOnCall(1, goodServiceKeyResponse, nil)
			err := pluginConnection.GetConnectionData(&commandData)
			Expect(err).NotTo(HaveOccurred())
			Expect(cliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(2))
			Expect(commandData.ConnnectionData.Username).To(Equal("cluster_operator_M5Scgeb0b6yp5f99E6SA8w"))
			Expect(commandData.ConnnectionData.Password).To(Equal("AMmxU9H6J5KSCYDLccipIw"))
			Expect(commandData.ConnnectionData.LocatorAddress).To(Equal("https://cloudcache-45371efd-f4ca-4549-a5f2-e06330aa53dc.sys.riverbank.cf-app.com"))
		})
	})

	Context("We don't have a service-key", func() {

		It("Returns an error indicating that there is no service-key", func() {
			cliConnection.CliCommandWithoutTerminalOutputReturnsOnCall(0, []string{"", ""}, nil)
			cliConnection.CliCommandWithoutTerminalOutputReturnsOnCall(1, []string{"", ""}, nil)
			err := pluginConnection.GetConnectionData(&commandData)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf(common.NoServiceKeyMessage, commandData.Target, commandData.Target)))
			Expect(cliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(len(commandData.ConnnectionData.Username)).To(BeZero())
			Expect(len(commandData.ConnnectionData.Password)).To(BeZero())
			Expect(len(commandData.ConnnectionData.LocatorAddress)).To(BeZero())
		})
	})
})
