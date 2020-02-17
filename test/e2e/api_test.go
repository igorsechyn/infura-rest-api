// +build integration

package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	apiClient *APIClient
)

var _ = Describe("Ethereum API", func() {
	BeforeSuite(func() {
		client, err := NewAPIClient()
		if err != nil {
			Fail("Failed to create a client")
		}

		apiClient = client
	})

	Describe("transactions", func() {
		It("should respond with 200 for /healthcheck", func() {
			status, err := apiClient.Healthcheck()

			Expect(err).Should(BeNil())
			Expect(status).Should(Equal(200))
		})

		It("should get transaction by block number and index", func() {
			transaction, status, err := apiClient.GetTransaction("2343", "1", map[string]string{})

			Expect(err).Should(BeNil())
			Expect(status).Should(Equal(200))
			Expect(transaction).ShouldNot(BeNil())
		})

		It("should return 400 status code for wrong block and index values", func() {
			transaction, status, err := apiClient.GetTransaction("wrong", "1", map[string]string{})

			Expect(err).Should(BeNil())
			Expect(status).Should(Equal(400))
			Expect(transaction).Should(BeNil())
		})
	})
})
