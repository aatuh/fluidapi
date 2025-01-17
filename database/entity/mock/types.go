package mock

import "github.com/stretchr/testify/mock"

// MockErrorChecker is a mock implementation of the ErrorChecker interface.
type MockErrorChecker struct {
	mock.Mock
}

func (m *MockErrorChecker) Check(err error) error {
	args := m.Called(err)
	return args.Error(0)
}
