package mock

import "github.com/stretchr/testify/mock"

// MockSQLUtil is a mock implementation of the ErrorChecker interface.
type MockSQLUtil struct {
	mock.Mock
}

func (m *MockSQLUtil) Check(err error) error {
	args := m.Called(err)
	return args.Error(0)
}
