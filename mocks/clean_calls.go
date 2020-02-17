package mocks

import "github.com/stretchr/testify/mock"

func CleanUpCallsFor(mockInstance *mock.Mock, method string) {
	calls := []*mock.Call{}
	for _, call := range mockInstance.ExpectedCalls {
		if call.Method != method {
			calls = append(calls, call)
		}
	}
	mockInstance.ExpectedCalls = calls
}
