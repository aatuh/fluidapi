package entity

import (
	"errors"
	"testing"

	entitymock "github.com/pakkasys/fluidapi/database/entity/mock"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSQLResult is a mock implementation of the sql.Result interface.
type MockSQLResult struct {
	mock.Mock
}

func (m *MockSQLResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSQLResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

type TestCreateEntity struct {
	ID   int
	Name string
	Age  int
}

// TestCreateEntity_NormalOperation tests the normal operation of CreateEntity.
func TestCreateEntity_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Create an Inserter function
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	// Call CreateEntity
	entity := &TestCreateEntity{ID: 1, Name: "Alice"}
	id, err := Create(entity, mockDB, "user", inserter, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestCreateEntity_InsertError tests the case where the insert function returns
// an error.
func TestCreateEntity_InsertError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Create an Inserter function
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).
		Return(nil, errors.New("prepare error"))
	mockSQLUtil.On("CheckDBError", mock.Anything).
		Return(errors.New("prepare error"))

	// Call CreateEntity
	entity := &TestCreateEntity{ID: 1, Name: "Alice"}
	_, err := Create(entity, mockDB, "user", inserter, mockSQLUtil)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
	mockSQLUtil.AssertExpectations(t)
}

// TestCreateEntities_NormalOperation tests the normal operation of
// CreateEntities with multiple entities.
func TestCreateEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Create an Inserter function
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(1), nil)

	// Call CreateEntities
	entities := []*TestCreateEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	id, err := CreateMany(entities, mockDB, "user", inserter, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestCreateEntities_EmptyEntities tests the case where no entities are passed
// to CreateEntities.
func TestCreateEntities_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Call CreateEntities with an empty list
	id, err := CreateMany([]*TestCreateEntity{}, mockDB, "user", nil, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), id)
	mockDB.AssertExpectations(t)
}

// TestCreateEntities_InsertError tests the case where an error occurs during
// insertion.
func TestCreateEntities_InsertError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Create an Inserter function
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).
		Return(nil, errors.New("prepare error"))
	mockSQLUtil.On("CheckDBError", mock.Anything).
		Return(errors.New("prepare error"))

	// Call CreateEntities
	entities := []*TestCreateEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	_, err := CreateMany(entities, mockDB, "user", inserter, mockSQLUtil)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
	mockSQLUtil.AssertExpectations(t)
}

// TestCheckInsertResult_NoError tests the case where there is no error and the
// result returns an ID.
func TestCheckInsertResult_NoError(t *testing.T) {
	mockSQLUtil := new(entitymock.MockSQLUtil)

	mockResult := new(MockSQLResult)
	mockResult.On("LastInsertId").Return(int64(123), nil)

	id, err := checkInsertResult(mockResult, nil, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), id)
	mockResult.AssertExpectations(t)
}

// TestCheckInsertResult_GeneralError tests the case where a general error is
// passed.
func TestCheckInsertResult_GeneralError(t *testing.T) {
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Simulate a general error
	generalErr := errors.New("some other error")
	mockSQLUtil.On("CheckDBError", generalErr).Return(generalErr)

	_, err := checkInsertResult(nil, generalErr, mockSQLUtil)

	assert.EqualError(t, err, "some other error")
	mockSQLUtil.AssertExpectations(t)
}

// TestCheckInsertResult_LastInsertIdError tests the case where getting the last
// insert ID returns an error.
func TestCheckInsertResult_LastInsertIdError(t *testing.T) {
	mockResult := new(MockSQLResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	mockResult.On("LastInsertId").
		Return(int64(0), errors.New("last insert ID error"))

	_, err := checkInsertResult(mockResult, nil, mockSQLUtil)

	assert.EqualError(t, err, "last insert ID error")
	mockResult.AssertExpectations(t)
}

// TestInsert_NormalOperation tests the normal operation of insert.
func TestInsert_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Inserter function for the entity
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)
	mockStmt.On("Close").Return(nil)

	_, err := insert(
		mockDB,
		&TestCreateEntity{ID: 1, Name: "Alice"},
		"user",
		inserter,
	)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestInsert_PrepareError tests the case where Prepare returns an error.
func TestInsert_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Inserter function for the entity
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Simulate an error on Prepare
	mockDB.On("Prepare", mock.Anything).
		Return(nil, errors.New("prepare error"))

	_, err := insert(
		mockDB,
		&TestCreateEntity{ID: 1, Name: "Alice"},
		"user",
		inserter,
	)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestInsert_ExecError tests the case where Exec returns an error.
func TestInsert_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Inserter function for the entity
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error on Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := insert(
		mockDB,
		&TestCreateEntity{ID: 1, Name: "Alice"},
		"user",
		inserter,
	)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestInsertMany_NormalOperation tests the normal operation of insertMany.
func TestInsertMany_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Inserter function for multiple entities
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	entities := []*TestCreateEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)
	mockStmt.On("Close").Return(nil)

	_, err := insertMany(mockDB, entities, "user", inserter)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestInsertMany_EmptyEntities tests the case where the inserted entity list is
// empty.
func TestInsertMany_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	entities := []*TestCreateEntity{}

	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Close").Return(nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, nil)

	_, err := insertMany(mockDB, entities, "user", nil)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

// TestInsertMany_PrepareError tests the case where Prepare returns an error.
func TestInsertMany_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Inserter function for multiple entities
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	entities := []*TestCreateEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := insertMany(mockDB, entities, "user", inserter)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestInsertMany_ExecError tests the case where Exec returns an error.
func TestInsertMany_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Inserter function for multiple entities
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		if entity.ID == 1 {
			return []string{"id", "name"}, []any{1, "Alice"}
		}
		return []string{"id", "name"}, []any{2, "Bob"}
	}

	entities := []*TestCreateEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error on Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := insertMany(mockDB, entities, "user", inserter)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}
