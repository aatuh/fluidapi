package connection

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/database/util/mock"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
)

// TestConnect tests the Connect function
func TestConnect(t *testing.T) {
	tests := []struct {
		name          string
		cfg           Config
		mockDBSetup   func(mockDB *mock.MockDB)
		mockFactory   func(mockDB *mock.MockDB) DriverFactory
		expectedError error
	}{
		{
			name: "Successful Connection",
			cfg: Config{
				Driver:          MySQL,
				ConnMaxLifetime: 30 * time.Minute,
				ConnMaxIdleTime: 10 * time.Minute,
				MaxOpenConns:    100,
				MaxIdleConns:    10,
			},
			mockDBSetup: func(mockDB *mock.MockDB) {
				mockDB.On("Ping").Return(nil).Once()
				mockDB.On("SetConnMaxLifetime", 30*time.Minute).Once()
				mockDB.On("SetConnMaxIdleTime", 10*time.Minute).Once()
				mockDB.On("SetMaxOpenConns", 100).Once()
				mockDB.On("SetMaxIdleConns", 10).Once()
			},
			mockFactory: func(mockDB *mock.MockDB) DriverFactory {
				return func(driver string, dsn string) (util.DB, error) {
					return mockDB, nil
				}
			},
			expectedError: nil,
		},
		{
			name: "Driver Factory Fails",
			cfg: Config{
				Driver: MySQL,
			},
			mockDBSetup: nil,
			mockFactory: func(mockDB *mock.MockDB) DriverFactory {
				return func(driver string, dsn string) (util.DB, error) {
					return nil, errors.New("failed to create driver")
				}
			},
			expectedError: errors.New("failed to open database: failed to create driver"),
		},
		{
			name: "Ping Fails",
			cfg: Config{
				Driver:          MySQL,
				ConnMaxLifetime: 30 * time.Minute,
				ConnMaxIdleTime: 10 * time.Minute,
				MaxOpenConns:    100,
				MaxIdleConns:    10,
			},
			mockDBSetup: func(mockDB *mock.MockDB) {
				mockDB.On("SetConnMaxLifetime", 30*time.Minute).Once()
				mockDB.On("SetConnMaxIdleTime", 10*time.Minute).Once()
				mockDB.On("SetMaxOpenConns", 100).Once()
				mockDB.On("SetMaxIdleConns", 10).Once()
				mockDB.On("Ping").Return(errors.New("ping failed")).Once()
			},
			mockFactory: func(mockDB *mock.MockDB) DriverFactory {
				return func(driver string, dsn string) (util.DB, error) {
					return mockDB, nil
				}
			},
			expectedError: errors.New("failed to ping database: ping failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(mock.MockDB)
			if tt.mockDBSetup != nil {
				tt.mockDBSetup(mockDB)
			}

			dbFactory := tt.mockFactory(mockDB)

			// Act
			db, err := Connect(tt.cfg, dbFactory, "dsn-string")

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}

			if tt.mockDBSetup != nil {
				mockDB.AssertExpectations(t)
			}
		})
	}
}

// TestGetDSN tests the GetDSN function
func TestGetDSN(t *testing.T) {
	tests := []struct {
		name          string
		cfg           Config
		expectedDSN   string
		expectedError error
	}{
		{
			name: "MySQL DSN",
			cfg: Config{
				Driver:         MySQL,
				User:           "user",
				Password:       "pass",
				ConnectionType: TCP,
				Host:           "localhost",
				Port:           3306,
				Database:       "testdb",
				Parameters:     "charset=utf8mb4&parseTime=True",
			},
			expectedDSN:   "user:pass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True",
			expectedError: nil,
		},
		{
			name: "Postgres DSN",
			cfg: Config{
				Driver:     Postgres,
				User:       "user",
				Password:   "pass",
				Host:       "localhost",
				Port:       5432,
				Database:   "testdb",
				Parameters: "sslmode=disable",
			},
			expectedDSN:   "postgres://user:pass@localhost:5432/testdb?sslmode=disable",
			expectedError: nil,
		},
		{
			name: "SQLite3 DSN",
			cfg: Config{
				Driver:     SQLite3,
				Database:   "testdb.sqlite",
				Parameters: "mode=memory",
			},
			expectedDSN:   "testdb.sqlite?mode=memory",
			expectedError: nil,
		},
		{
			name: "Unsupported Driver",
			cfg: Config{
				Driver: "unknown",
			},
			expectedDSN:   "",
			expectedError: fmt.Errorf("unsupported driver: unknown"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn, err := GetDSN(tt.cfg)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, dsn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dsn)
				assert.Equal(t, tt.expectedDSN, *dsn)
			}
		})
	}
}

// TestConfigureConnection tests the configureConnection function.
func TestConfigureConnection(t *testing.T) {
	// Create a mock database
	mockDB := new(utilmock.MockDB)

	// Expected configuration
	cfg := Config{
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
		MaxOpenConns:    100,
		MaxIdleConns:    10,
	}

	// Setup expectations on the mock
	mockDB.On("SetConnMaxLifetime", cfg.ConnMaxLifetime).Once()
	mockDB.On("SetConnMaxIdleTime", cfg.ConnMaxIdleTime).Once()
	mockDB.On("SetMaxOpenConns", cfg.MaxOpenConns).Once()
	mockDB.On("SetMaxIdleConns", cfg.MaxIdleConns).Once()

	// Call the function under test
	configureConnection(mockDB, cfg)

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}
