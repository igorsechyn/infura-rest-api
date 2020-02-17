package transactions_test

import (
	"ethereum-api/mocks"
	"ethereum-api/pkg/app"
	"ethereum-api/pkg/app/errors"
	"ethereum-api/pkg/app/reporter"
	"ethereum-api/pkg/app/transactions"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

const validInfuraResponse = `{"jsonrpc":"2.0","id":"1","result":{"blockHash":"0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35","blockNumber":"0x5bad55","from":"0xc837f51a0efa33f8eca03570e3d01a4b2cf97ffd","gas":"0x15f90","gasPrice":"0x14b8d03a00","hash":"0x311be6a9b58748717ac0f70eb801d29973661aaf1365960d159e4ec4f4aa2d7f","input":"0x","nonce":"0x4241","r":"0xe9ef2f6fcff76e45fac6c2e8080094370082cfb47e8fde0709312f9aa3ec06ad","s":"0x421ebc4ebe187c173f13b1479986dcbff5c4997c0dfeb1fd149a982ad4bcdfe7","to":"0xf49bd0367d830850456d2259da366a054038dc46","transactionIndex":"0x1","v":"0x25","value":"0x1bafa9ee16e78000"}}`

var _ = Describe("Transactions Service", func() {
	var (
		mockInfuraClient *mocks.MockInfuraClient
		mockLogger       *mocks.MockLoggerSink
	)

	BeforeEach(func() {
		mockInfuraClient = mocks.CreateMockInfuraClient()
		mockLogger = mocks.CreateMockLoggerSink()
	})

	Describe("retrive transaction", func() {
		Context("by block number and index", func() {

			whenGetTransactionByBlockNumberAndIndexIsCalled := func(blockNumber, index string) (*transactions.Transaction, error) {
				reporter := reporter.New(mockLogger)
				service := app.CreateTransactionsService(mockInfuraClient, reporter)
				return service.GetTransactionByBlockNumberAndIndex(blockNumber, index)
			}

			DescribeTable("block and index value parsing",
				func(blockNumber, index string, params []string, expectedError error) {
					_, err := whenGetTransactionByBlockNumberAndIndexIsCalled(blockNumber, index)
					if expectedError != nil {
						Expect(err).Should(BeEquivalentTo(expectedError))
					} else {
						Expect(err).Should(BeNil())
						mockInfuraClient.AssertCalled(GinkgoT(), "Execute", "eth_getTransactionByBlockNumberAndIndex", params)
					}
				},
				Entry("should call infura client to retrieve transaction with numeric block number", "8", "1", []string{"0x8", "0x1"}, nil),
				Entry("should call infura client to retrieve transaction from latest block number", "latest", "1", []string{"latest", "0x1"}, nil),
				Entry("should call infura client to retrieve transaction from latest block number", "earliest", "1", []string{"earliest", "0x1"}, nil),
				Entry("should call infura client to retrieve transaction from latest block number", "pending", "1", []string{"pending", "0x1"}, nil),
				Entry("should return an error, if parsing block value fails", "invalid value", "1", []string{}, fmt.Errorf("Error parsing block number value: strconv.ParseInt: parsing \"invalid value\": invalid syntax: %w", errors.BadInputErr)),
				Entry("should return an error, if parsing index value fails", "8", "invalid value", []string{}, fmt.Errorf("Error parsing index number value: strconv.ParseInt: parsing \"invalid value\": invalid syntax: %w", errors.BadInputErr)),
			)

			It("should log an error, if retrieving transaction fails", func() {
				mockInfuraClient.GivenExecuteFails(fmt.Errorf("Some RPC Error"))

				_, err := whenGetTransactionByBlockNumberAndIndexIsCalled("8", "2")

				Expect(err).ShouldNot(BeNil())
				mockLogger.AssertCalled(GinkgoT(), "Error", fmt.Errorf("Error retrieving transation: Some RPC Error"), map[string]interface{}{"code": "transactions.get.blocknumberandindex.error", "block": "8", "index": "2"})
			})

			It("should return an error, if retrieving transaction fails", func() {
				mockInfuraClient.GivenExecuteFails(fmt.Errorf("Some RPC Error"))

				_, err := whenGetTransactionByBlockNumberAndIndexIsCalled("8", "2")

				Expect(err).To(Equal(fmt.Errorf("Error retrieving transation: Some RPC Error")))
			})

			It("should log an error, if parsing response fails", func() {
				mockInfuraClient.GivenExecuteSucceeds([]byte("{invalid response"))

				_, err := whenGetTransactionByBlockNumberAndIndexIsCalled("8", "2")

				Expect(err).ShouldNot(BeNil())
				mockLogger.AssertCalled(GinkgoT(), "Error", fmt.Errorf("Error parsing response: invalid character 'i' looking for beginning of object key string"), map[string]interface{}{"code": "transactions.get.blocknumberandindex.error", "block": "8", "index": "2"})
			})

			It("should return an error, if retrieving transaction fails", func() {
				mockInfuraClient.GivenExecuteSucceeds([]byte("{invalid response"))

				_, err := whenGetTransactionByBlockNumberAndIndexIsCalled("8", "2")

				Expect(err).To(Equal(fmt.Errorf("Error parsing response: invalid character 'i' looking for beginning of object key string")))
			})

			It("should return transaction, if it was found", func() {
				mockInfuraClient.GivenExecuteSucceeds([]byte(validInfuraResponse))

				response, err := whenGetTransactionByBlockNumberAndIndexIsCalled("8", "2")

				Expect(err).Should(BeNil())
				expectedTransaction := transactions.Transaction{
					BlockHash:   "0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35",
					BlockNumber: "0x5bad55",
					From:        "0xc837f51a0efa33f8eca03570e3d01a4b2cf97ffd",
					Gas:         "0x15f90",
					GasPrice:    "0x14b8d03a00",
					Hash:        "0x311be6a9b58748717ac0f70eb801d29973661aaf1365960d159e4ec4f4aa2d7f",
					Input:       "0x",
					Nonce:       "0x4241",
					To:          "0xf49bd0367d830850456d2259da366a054038dc46",
					Index:       "0x1",
					Value:       "0x1bafa9ee16e78000",
				}
				Expect(response).To(Equal(&expectedTransaction))
			})

			It("should log an info, if transaction was found", func() {
				mockInfuraClient.GivenExecuteSucceeds([]byte(validInfuraResponse))

				_, err := whenGetTransactionByBlockNumberAndIndexIsCalled("8", "2")

				Expect(err).Should(BeNil())
				mockLogger.AssertCalled(GinkgoT(), "Info", fmt.Sprintf("Found transaction in block 8 at index 2"), map[string]interface{}{"code": "transactions.get.blocknumberandindex", "block": "8", "index": "2"})
			})

			It("should return nil, if transaction was not found", func() {
				mockInfuraClient.GivenExecuteSucceeds([]byte(`{"jsonrpc":"2.0","id":"1","result":null}`))

				response, err := whenGetTransactionByBlockNumberAndIndexIsCalled("8", "2")

				Expect(err).To(BeNil())
				Expect(response).To(BeNil())
			})

			It("should log an info, if transaction was not found", func() {
				mockInfuraClient.GivenExecuteSucceeds([]byte(`{"jsonrpc":"2.0","id":"1","result":null}`))

				_, err := whenGetTransactionByBlockNumberAndIndexIsCalled("8", "2")

				Expect(err).Should(BeNil())
				mockLogger.AssertCalled(GinkgoT(), "Info", fmt.Sprintf("No transaction found in block 8 at index 2"), map[string]interface{}{"code": "transactions.get.blocknumberandindex.notfound", "block": "8", "index": "2"})
			})
		})
	})
})
