package infura_test

import (
	"ethereum-api/mocks"
	"ethereum-api/pkg/app/ethereum/infura"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Http Infura Client", func() {
	var server *httptest.Server
	var mockAPI *mocks.MockInfuraAPI

	writeResponse := func(w http.ResponseWriter, response string, statusCode int) {
		w.WriteHeader(statusCode)
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}

	BeforeEach(func() {
		mockAPI = mocks.CreateMockInfuraAPI()
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				body, _ := ioutil.ReadAll(r.Body)
				response, statusCode := mockAPI.Post(string(body), r.URL.EscapedPath(), r.Header)
				writeResponse(w, response, statusCode)
				return
			}

			writeResponse(w, "Unsupported call", 400)
		})
		server = httptest.NewServer(mockHandler)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Execute", func() {
		It("should make a post call to retrieve transactions", func() {
			client := infura.NewInfuraClient(server.URL, "some-project-id")
			_, err := client.Execute("eth_getTransactionByBlockNumberAndIndex", "0x8", "0x2")

			Expect(err).Should(BeNil())
			mockAPI.AssertCalled(GinkgoT(), "Post", mock.Anything, mock.Anything, mock.Anything)
			body := mockAPI.Calls[0].Arguments[0].(string)
			path := mockAPI.Calls[0].Arguments[1].(string)
			headers := mockAPI.Calls[0].Arguments[2].(map[string][]string)
			Expect(body).Should(MatchJSON(`{"id": "1", "jsonrpc": "2.0", "method": "eth_getTransactionByBlockNumberAndIndex", "params": ["0x8", "0x2"]}`))
			Expect(path).Should(Equal("/v3/some-project-id"))
			Expect(headers).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
		})

		It("should return an error, if response code is not 2xx", func() {
			mockAPI.GivenPostResponds("\"Failed\"", 503)

			client := infura.NewInfuraClient(server.URL, "some-project-id")
			_, err := client.Execute("eth_getTransactionByBlockNumberAndIndex", "0x8", "0x2")

			Expect(err).To(Equal(fmt.Errorf("Server error. Status code: 503, Status: 503 Service Unavailable")))
		})
	})
})
