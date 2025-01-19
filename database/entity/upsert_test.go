package entity

import (
	"errors"
	"testing"

	"github.com/pakkasys/fluidapi/database/clause"
	entitymock "github.com/pakkasys/fluidapi/database/entity/mock"
	databasemock "github.com/pakkasys/fluidapi/database/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UpsertTestStruct struct {
	ID   int64
	Name string
}

// TestUpsert tests the Upsert function.
func TestUpsert(t *testing.T) {
	// Use TestUpsertMany to cover functionality since Upsert is a wrapper
	// that delegates to UpsertMany.
	t.Run("Delegates to UpsertMany", func(t *testing.T) {
		mockDB := new(databasemock.MockDB)
		mockStmt := new(databasemock.MockStmt)
		mockResult := new(databasemock.MockResult)
		mockErrorChecker := new(entitymock.MockErrorChecker)

		mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
		mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
		mockStmt.On("Close").Return(nil)
		mockResult.On("LastInsertId").Return(int64(1), nil)

		inserter := func(entity *UpsertTestStruct) ([]string, []any) {
			return []string{"id", "name"}, []any{entity.ID, entity.Name}
		}

		id, err := Upsert(
			mockDB,
			"user_table",
			&UpsertTestStruct{ID: 1, Name: "Alice"},
			inserter,
			[]clause.Projection{{Table: "name", Alias: "name_alias"}},
			mockErrorChecker,
		)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)

		mockDB.AssertExpectations(t)
		mockStmt.AssertExpectations(t)
		mockResult.AssertExpectations(t)
		mockErrorChecker.AssertExpectations(t)
	})
}

// TestUpsertMany tests the UpsertMany function
func TestUpsertMany(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *databasemock.MockDB,
			mockStmt *databasemock.MockStmt,
			mockResult *databasemock.MockResult,
			mockErrorChecker *entitymock.MockErrorChecker,
		)
		entities          []*UpsertTestStruct
		updateProjections []clause.Projection
		expectedID        int64
		expectedError     string
	}{
		{
			name: "Normal Operation",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockResult *databasemock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
				mockStmt.On("Close").Return(nil)
				mockResult.On("LastInsertId").Return(int64(1), nil)
			},
			entities: []*UpsertTestStruct{
				{ID: 1, Name: "Alice"},
			},
			updateProjections: []clause.Projection{
				{Table: "name", Alias: "name_alias"},
			},
			expectedID:    1,
			expectedError: "",
		},
		{
			name:       "No Entities",
			setupMocks: nil,
			entities:   []*UpsertTestStruct{},
			updateProjections: []clause.Projection{
				{Table: "name", Alias: "name_alias"},
			},
			expectedID:    0,
			expectedError: "must provide entities to upsert",
		},
		{
			name:       "No Update Projections",
			setupMocks: nil,
			entities: []*UpsertTestStruct{
				{ID: 1, Name: "Alice"},
			},
			updateProjections: []clause.Projection{},
			expectedID:        0,
			expectedError:     "must provide update projections",
		},
		{
			name:       "No Alias in Update Projections",
			setupMocks: nil,
			entities: []*UpsertTestStruct{
				{ID: 1, Name: "Alice"},
			},
			updateProjections: []clause.Projection{
				{Table: "name", Alias: ""},
			},
			expectedID:    0,
			expectedError: "must provide update projections alias",
		},
		{
			name: "Exec Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockResult *databasemock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
				mockStmt.On("Close").Return(nil)
				mockErrorChecker.On("Check", mock.Anything).Return(errors.New("exec error"))
			},
			entities: []*UpsertTestStruct{
				{ID: 1, Name: "Alice"},
			},
			updateProjections: []clause.Projection{
				{Table: "name", Alias: "name_alias"},
			},
			expectedID:    0,
			expectedError: "exec error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(databasemock.MockDB)
			mockStmt := new(databasemock.MockStmt)
			mockResult := new(databasemock.MockResult)
			mockErrorChecker := new(entitymock.MockErrorChecker)

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockResult, mockErrorChecker)
			}

			// Define inserter function
			inserter := func(entity *UpsertTestStruct) ([]string, []any) {
				return []string{"id", "name"}, []any{entity.ID, entity.Name}
			}

			// Act
			id, err := UpsertMany(
				mockDB,
				"user_table",
				tt.entities,
				inserter,
				tt.updateProjections,
				mockErrorChecker,
			)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockResult.AssertExpectations(t)
			mockErrorChecker.AssertExpectations(t)
		})
	}
}
