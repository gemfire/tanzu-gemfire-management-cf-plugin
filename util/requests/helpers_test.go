package requests_test

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/util/requests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Helpers", func() {

	Context("GetTargetAndClusterCommand", func() {
		var (
			args []string
		)

		BeforeEach(func() {

		})

		Context("with no target in environment", func() {
			It("returns no target, no command", func() {
				args = []string{"program"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal(""))
				Expect(userCommand.Command).To(Equal(""))
			})

			It("returns target but no command", func() {
				args = []string{"program", "target"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal(""))
			})

			It("returns target and multiple word command", func() {
				args = []string{"program", "target", "list", "members"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(0))
			})

			It("returns target, multiple word command and options ", func() {
				args = []string{"program", "target", "list", "members", "-h"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(1))
				Expect(userCommand.Parameters["-h"]).To(Equal("true"))
				Expect(userCommand.Parameters["-foo"]).To(Equal(""))
			})

			It("returns target, multiple word command and option with values ", func() {
				args = []string{"program", "target", "list", "members", "-t", "abc"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(1))
				Expect(userCommand.Parameters["-t"]).To(Equal("abc"))
				Expect(userCommand.Parameters["-foo"]).To(Equal(""))
			})

			It("returns target, multiple word command, option without value and option with values ", func() {
				args = []string{"program", "target", "list", "members", "-h", "-t", "abc"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(2))
				Expect(userCommand.Parameters["-t"]).To(Equal("abc"))
				Expect(userCommand.Parameters["-h"]).To(Equal("true"))
				Expect(userCommand.Parameters["-foo"]).To(Equal(""))
			})

			It("returns target, multiple word command, option with value and option without values ", func() {
				args = []string{"program", "target", "list", "members", "-t", "abc", "-h"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(2))
				Expect(userCommand.Parameters["-t"]).To(Equal("abc"))
				Expect(userCommand.Parameters["-h"]).To(Equal("true"))
				Expect(userCommand.Parameters["-foo"]).To(Equal(""))
			})
		})

		// for now, if you have target in the environment variable, you can not override it in
		// the individual command
		Context("with target in environment", func() {
			BeforeEach(func() {
				err := os.Setenv("CFPCC", "target")
				Expect(err).To(BeNil())
			})

			It("returns target but no command", func() {
				args = []string{"program", "target"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal(""))
			})

			It("returns target and command", func() {
				args = []string{"program", "command"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("command"))
			})

			It("returns target and multiple word command", func() {
				args = []string{"program", "target", "list", "members"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(0))
			})

			It("returns target and multiple word command", func() {
				args = []string{"program", "list", "members"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(0))
			})

			It("returns target, multiple word command and options ", func() {
				args = []string{"program", "target", "list", "members", "-h"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(1))
				Expect(userCommand.Parameters["-h"]).To(Equal("true"))
				Expect(userCommand.Parameters["-foo"]).To(Equal(""))
			})

			It("returns target, multiple word command and options ", func() {
				args = []string{"program", "list", "members", "-h"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(1))
				Expect(userCommand.Parameters["-h"]).To(Equal("true"))
				Expect(userCommand.Parameters["-foo"]).To(Equal(""))
			})

			It("returns target, multiple word command and option with values ", func() {
				args = []string{"program", "list", "members", "-t", "abc"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(1))
				Expect(userCommand.Parameters["-t"]).To(Equal("abc"))
				Expect(userCommand.Parameters["-foo"]).To(Equal(""))
			})

			It("returns target, multiple word command, option without value and option with values ", func() {
				args = []string{"program", "list", "members", "-h", "-t", "abc"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(2))
				Expect(userCommand.Parameters["-t"]).To(Equal("abc"))
				Expect(userCommand.Parameters["-h"]).To(Equal("true"))
				Expect(userCommand.Parameters["-foo"]).To(Equal(""))
			})

			It("returns target, multiple word command, option with value and option without values ", func() {
				args = []string{"program", "list", "members", "-t", "abc", "-h"}
				target, userCommand := requests.GetTargetAndClusterCommand(args)
				Expect(target).To(Equal("target"))
				Expect(userCommand.Command).To(Equal("list members"))
				Expect(len(userCommand.Parameters)).To(Equal(2))
				Expect(userCommand.Parameters["-t"]).To(Equal("abc"))
				Expect(userCommand.Parameters["-h"]).To(Equal("true"))
				Expect(userCommand.Parameters["-foo"]).To(Equal(""))
			})
		})
	})

	Context("GetEndPoints", func() {

		It("Does nothing", func() {
			Expect(nil).To(BeNil())
		})
	})

	Context("RequestToEndPoint", func() {

		It("Does nothing", func() {
			Expect(nil).To(BeNil())
		})
	})
})
