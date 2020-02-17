package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockInfuraAPI struct {
	mock.Mock
}

func (api *MockInfuraAPI) Post(body string, path string, header map[string][]string) (string, int) {
	args := api.Called(body, path, header)
	return args.String(0), args.Int(1)
}

func (api *MockInfuraAPI) GivenPostResponds(payload string, statusCode int) {
	CleanUpCallsFor(&api.Mock, "Post")
	api.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(payload, statusCode)
}

func CreateMockInfuraAPI() *MockInfuraAPI {
	mockServer := new(MockInfuraAPI)
	mockServer.GivenPostResponds("{}", 200)
	return mockServer
}
