package entity

import (
	"errors"
	"testing"

	entitymock "github.com/pakkasys/fluidapi/database/entity/mock"
	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUpsertEntity_NormalOperation tests the normal operation of UpsertEntity.
func TestUpsertEntity_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Test entity and projections
	entity := &TestEntity{ID: 1, Name: "Alice"}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for the entity
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(0), nil)

	_, err := Upsert(
		mockDB,
		"user",
		entity,
		inserter,
		projections,
		mockSQLUtil,
	)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestUpsertEntities_NormalOperation tests the normal operation of
// UpsertEntities.
func TestUpsertEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Test entities and projections
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for multiple entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("LastInsertId").Return(int64(0), nil)

	_, err := UpsertMany(
		mockDB,
		"user",
		entities,
		inserter,
		projections,
		mockSQLUtil,
	)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestUpsertEntities_EmptyEntities tests the case where the entities list is
// empty.
func TestUpsertEntities_EmptyEntities(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	expectedError := errors.New("must provide entities to upsert")
	mockSQLUtil.On("CheckDBError", expectedError).Return(expectedError)

	_, err := UpsertMany(
		mockDB,
		"user",
		[]*TestEntity{},
		func(e *TestEntity) ([]string, []any) {
			return []string{}, []any{}
		},
		[]util.Projection{{Column: "name", Alias: "test"}},
		mockSQLUtil,
	)

	assert.EqualError(t, err, "must provide entities to upsert")
	mockDB.AssertExpectations(t)
	mockSQLUtil.AssertExpectations(t)
}

// TestUpsertEntities_ErrorFromUpsertMany tests the case where upsertMany
// returns an error.
func TestUpsertEntities_ErrorFromUpsertMany(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Test entities and projections
	entities := []*TestEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	projections := []util.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for multiple entities
	inserter := func(e *TestEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Simulate an error in the upsertMany call
	mockDB.On("Prepare", mock.Anything).
		Return(nil, errors.New("prepare error"))
	mockSQLUtil.On("CheckDBError", mock.Anything).
		Return(errors.New("prepare error"))

	_, err := UpsertMany(
		mockDB,
		"user",
		entities,
		inserter,
		projections,
		mockSQLUtil,
	)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
	mockSQLUtil.AssertExpectations(t)
}
