package mocks

import "github.com/stretchr/testify/mock"

type MockInfuraClient struct {
	mock.Mock
}

func (client *MockInfuraClient) Execute(method string, params ...string) ([]byte, error) {
	args := client.Called(method, params)
	payload := args.Get(0)
	if payload == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (client *MockInfuraClient) GivenExecuteSucceeds(returnValue []byte) {
	CleanUpCallsFor(&client.Mock, "Execute")
	client.On("Execute", mock.Anything, mock.Anything).Return(returnValue, nil)
}

func (client *MockInfuraClient) GivenExecuteFails(err error) {
	CleanUpCallsFor(&client.Mock, "Execute")
	client.On("Execute", mock.Anything, mock.Anything).Return(nil, err)
}

func CreateMockInfuraClient() *MockInfuraClient {
	mockInfuraClient := new(MockInfuraClient)
	mockInfuraClient.GivenExecuteSucceeds([]byte("{}"))
	return mockInfuraClient
}
