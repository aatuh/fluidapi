package mock

import (
	"context"
	"database/sql"
	"time"

	"github.com/pakkasys/fluidapi/database"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the DB interface.
type MockDB struct {
	mock.Mock
}

var _ database.DB = (*MockDB)(nil)

func (m *MockDB) Prepare(query string) (database.Stmt, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(database.Stmt), args.Error(1)
	}
}

func (m *MockDB) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) SetConnMaxLifetime(d time.Duration) {
	m.Called(d)
}

func (m *MockDB) SetConnMaxIdleTime(d time.Duration) {
	m.Called(d)
}

func (m *MockDB) SetMaxOpenConns(n int) {
	m.Called(n)
}

func (m *MockDB) SetMaxIdleConns(n int) {
	m.Called(n)
}

func (m *MockDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (database.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(database.Tx), args.Error(1)
}

func (m *MockDB) Exec(query string, args ...any) (database.Result, error) {
	calledArgs := m.Called(query, args)
	return calledArgs.Get(0).(database.Result), calledArgs.Error(1)
}

func (m *MockDB) Query(query string, args ...any) (database.Rows, error) {
	calledArgs := m.Called(query, args)
	return calledArgs.Get(0).(database.Rows), calledArgs.Error(1)
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockTx is a mock implementation of the Tx interface.
type MockTx struct {
	MockDB
}

var _ database.Tx = (*MockTx)(nil)

func (m *MockTx) Commit() error {
	return m.Called().Error(0)
}

func (m *MockTx) Rollback() error {
	return m.Called().Error(0)
}

// MockStmt is a mock implementation of the Stmt interface.
type MockStmt struct {
	mock.Mock
}

var _ database.Stmt = (*MockStmt)(nil)

func (m *MockStmt) Close() error {
	return m.Called().Error(0)
}

func (m *MockStmt) QueryRow(args ...any) database.Row {
	argsCalled := m.Called(args)
	return argsCalled.Get(0).(database.Row)
}

func (m *MockStmt) Exec(args ...any) (database.Result, error) {
	argsCalled := m.Called(args)
	if argsCalled.Get(0) == nil {
		return nil, argsCalled.Error(1)
	}
	return argsCalled.Get(0).(database.Result), argsCalled.Error(1)
}

func (m *MockStmt) Query(args ...any) (database.Rows, error) {
	argsCalled := m.Called(args)
	if argsCalled.Get(0) == nil {
		return nil, argsCalled.Error(1)
	}
	return argsCalled.Get(0).(database.Rows), argsCalled.Error(1)
}

// MockRows is a mock implementation of the Rows interface.
type MockRows struct {
	mock.Mock
}

var _ database.Rows = (*MockRows)(nil)

func (m *MockRows) Scan(dest ...any) error {
	return m.Called(dest).Error(0)
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Close() error {
	return m.Called().Error(0)
}

func (m *MockRows) Err() error {
	return m.Called().Error(0)
}

// MockRow is a mock implementation of the Row interface.
type MockRow struct {
	mock.Mock
}

var _ database.Row = (*MockRow)(nil)

func (m *MockRow) Scan(dest ...any) error {
	return m.Called(dest).Error(0)
}

func (m *MockRow) Err() error {
	return m.Called().Error(0)
}

// MockResult is a mock implementation of the Result interface.
type MockResult struct {
	mock.Mock
}

var _ database.Result = (*MockResult)(nil)

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
