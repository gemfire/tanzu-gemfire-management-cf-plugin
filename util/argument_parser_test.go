package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	. "github.com/gemfire/cloudcache-management-cf-plugin/util"
)

var _ = Describe("ArgumentParser", func() {

	var (
		args        []string
		commandData domain.CommandData
		value       interface{}
		ok          bool
	)

	Context("All parameters provided are well formatted", func() {

		Context("All parameters provided are string values", func() {

			BeforeEach(func() {
				args = []string{
					"-r=regionOne",
					"-g=groupOne",
					"--data=@someJSONFile",
				}
				commandData = domain.CommandData{}
			})

			It("Parses all values as expected", func() {
				ParseArguments(args, &commandData)
				value, ok = commandData.UserCommand.Parameters["r"]
				Expect(ok).To(Equal(true))
				Expect(value).To(Equal("regionOne"))
				value, ok = commandData.UserCommand.Parameters["g"]
				Expect(ok).To(Equal(true))
				Expect(value).To(Equal("groupOne"))
				value, ok = commandData.UserCommand.Parameters["data"]
				Expect(ok).To(Equal(true))
				Expect(value).To(Equal("@someJSONFile"))
			})
		})

		Context("Some parameters are boolean values", func() {

			BeforeEach(func() {
				args = []string{
					"-r=regionOne",
					"-g=groupOne",
					"--data=@someJSONFile",
					"-j",
					"-isItTrue=false",
				}
				commandData = domain.CommandData{}
			})

			It("Parses all values as expected", func() {
				ParseArguments(args, &commandData)
				value, ok = commandData.UserCommand.Parameters["r"]
				Expect(ok).To(Equal(true))
				Expect(value).To(Equal("regionOne"))
				value, ok = commandData.UserCommand.Parameters["g"]
				Expect(ok).To(Equal(true))
				Expect(value).To(Equal("groupOne"))
				value, ok = commandData.UserCommand.Parameters["data"]
				Expect(ok).To(Equal(true))
				Expect(value).To(Equal("@someJSONFile"))
				value, ok = commandData.UserCommand.Parameters["j"]
				Expect(ok).To(Equal(true))
				Expect(value).To(Equal(true))
				value, ok = commandData.UserCommand.Parameters["isItTrue"]
				Expect(ok).To(Equal(true))
				Expect(value).To(Equal(false))
			})
		})
	})

	Context("Some parameters provided are not well formatted", func() {

		BeforeEach(func() {
			args = []string{
				"-r=regionOne",
				"-g=groupOne",
				"--data=@someJSONFile",
				"-j",
				"-isItTrue=false",
				"wrong=wrong",
				"noparam",
			}
			commandData = domain.CommandData{}
		})

		It("Ignores the poorly formatted parameters and parses the rest", func() {
			ParseArguments(args, &commandData)
			value, ok = commandData.UserCommand.Parameters["r"]
			Expect(ok).To(Equal(true))
			Expect(value).To(Equal("regionOne"))
			value, ok = commandData.UserCommand.Parameters["g"]
			Expect(ok).To(Equal(true))
			Expect(value).To(Equal("groupOne"))
			value, ok = commandData.UserCommand.Parameters["data"]
			Expect(ok).To(Equal(true))
			Expect(value).To(Equal("@someJSONFile"))
			value, ok = commandData.UserCommand.Parameters["j"]
			Expect(ok).To(Equal(true))
			Expect(value).To(Equal(true))
			value, ok = commandData.UserCommand.Parameters["isItTrue"]
			Expect(ok).To(Equal(true))
			Expect(value).To(Equal(false))
			value, ok = commandData.UserCommand.Parameters["wrong"]
			Expect(ok).To(Equal(false))
			Expect(value).To(BeNil())
			value, ok = commandData.UserCommand.Parameters["noparam"]
			Expect(ok).To(Equal(false))
			Expect(value).To(BeNil())
		})
	})

})
