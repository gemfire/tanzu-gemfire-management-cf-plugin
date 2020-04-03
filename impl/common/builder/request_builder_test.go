package builder_test

import (
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/domain"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common"
	. "github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common/builder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestBuilder", func() {

	var (
		buildRequest      common.RequestBuilder
		restEndPoint      domain.RestEndPoint
		commandData       domain.CommandData
		expectedRegionURL string
	)

	buildRequest = BuildRequest

	BeforeEach(func() {
		restEndPoint = domain.RestEndPoint{}
		commandData = domain.CommandData{}
		expectedRegionURL = "http://localhost:7070/management/regions"
	})

	Context("Build requests using RestEndPoint and CommandData", func() {

		BeforeEach(func() {
			restEndPoint.HTTPMethod = "POST"
			restEndPoint.URL = "/regions"
			restEndPoint.CommandName = "create region"
			restEndPoint.Parameters = []domain.RestAPIParam{domain.RestAPIParam{Name: "regionConfig", In: "body", Required: true}}

			commandData.ConnnectionData.LocatorAddress = "http://localhost:7070"
			commandData.UserCommand.Command = "create region"
			commandData.UserCommand.Parameters = make(map[string]string)
		})

		Context("Request with a body", func() {

			Context("Body from file, where file is found", func() {

				It("Returns URL, bodyReader and nil error", func() {
					commandData.UserCommand.Parameters["--regionConfig"] = "@../../../testdata/request-body.json"
					request, err := buildRequest(restEndPoint, &commandData)
					Expect(err).NotTo(HaveOccurred())
					Expect(request).NotTo(BeNil())
					Expect(request.URL.String()).To(Equal(expectedRegionURL))
					Expect(request.Body).NotTo(BeNil())
				})
			})

			Context("Body from file, where file is not found", func() {
				It("Returns an error", func() {
					commandData.UserCommand.Parameters["--regionConfig"] = "@../../../testdata/notfound-body.json"
					request, err := buildRequest(restEndPoint, &commandData)
					Expect(err).To(HaveOccurred())
					Expect(request).To(BeNil())
				})
			})

			Context("Body direct from command line", func() {
				It("Returns URL, bodyReader and nil error", func() {
					commandData.UserCommand.Parameters["--regionConfig"] = `{"name": "testRegion", "type": "PARTITION"}`
					request, err := buildRequest(restEndPoint, &commandData)
					Expect(err).NotTo(HaveOccurred())
					Expect(request).NotTo(BeNil())
					Expect(request.URL.String()).To(Equal(expectedRegionURL))
					Expect(request.Body).NotTo(BeNil())
				})
			})
		})

		Context("Request with path parameters", func() {
			var expectedDeleteURL string

			BeforeEach(func() {
				restEndPoint.URL = "/regions/{regionName}/indexes/{indexName}"
				restEndPoint.CommandName = "DELETE"
				restEndPoint.Parameters = []domain.RestAPIParam{
					domain.RestAPIParam{Name: "regionName", In: "path", Required: true},
					domain.RestAPIParam{Name: "indexName", In: "path", Required: true},
				}

				commandData.UserCommand.Command = "delete region index"
				commandData.UserCommand.Parameters["--regionName"] = "testRegion"
				commandData.UserCommand.Parameters["--indexName"] = "testIndex"

				expectedDeleteURL = "http://localhost:7070/management/regions/testRegion/indexes/testIndex"
			})

			It("Returns URL, nil bodyReader and nil error", func() {
				request, err := buildRequest(restEndPoint, &commandData)
				Expect(err).NotTo(HaveOccurred())
				Expect(request).NotTo(BeNil())
				Expect(request.URL.String()).To(Equal(expectedDeleteURL))
				Expect(request.Body).To(BeNil())
			})
		})

		Context("Request with query parameters", func() {
			var expectedListURL string

			BeforeEach(func() {
				restEndPoint.URL = "/members"
				restEndPoint.CommandName = "GET"
				restEndPoint.Parameters = []domain.RestAPIParam{
					domain.RestAPIParam{Name: "group", In: "query", Required: false},
					domain.RestAPIParam{Name: "id", In: "query", Required: false},
				}

				commandData.UserCommand.Command = "list members"
				commandData.UserCommand.Parameters["--group"] = "testGroup"
				commandData.UserCommand.Parameters["--id"] = "testId"

				expectedListURL = "http://localhost:7070/management/members?group=testGroup&id=testId"
			})

			It("Returns URL, nil bodyReader and nil error", func() {
				request, err := buildRequest(restEndPoint, &commandData)
				Expect(err).NotTo(HaveOccurred())
				Expect(request).NotTo(BeNil())
				Expect(request.URL.String()).To(Equal(expectedListURL))
				Expect(request.Body).To(BeNil())
			})
		})

		Context("Request with formData", func() {
			var expectedListURL string

			BeforeEach(func() {
				restEndPoint.URL = "/deployment"
				restEndPoint.CommandName = "PUT"
				restEndPoint.Consumes = []string{"multipart/form-data"}
				restEndPoint.Parameters = []domain.RestAPIParam{
					domain.RestAPIParam{Name: "config", In: "query", Type: "string", Required: true},
					domain.RestAPIParam{Name: "file", In: "formData", Type: "file", Required: true},
				}

				commandData.UserCommand.Command = "deploy"
				commandData.UserCommand.Parameters["--config"] = "{}"
				commandData.UserCommand.Parameters["--file"] = "../../../testdata/request-body.json"

				expectedListURL = "http://localhost:7070/management/deployment?config=%7B%7D"
			})

			It("Returns URL, nil bodyReader and nil error", func() {
				request, err := buildRequest(restEndPoint, &commandData)
				Expect(err).NotTo(HaveOccurred())
				Expect(request).NotTo(BeNil())
				Expect(request.URL.String()).To(Equal(expectedListURL))
			})
		})
	})
})
