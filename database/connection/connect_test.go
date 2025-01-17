package connection

import (
	"testing"

	"github.com/pakkasys/fluidapi/database/util"
	"github.com/stretchr/testify/assert"
)

// TestConnect_UnsupportedConnectionType unsupported connection type case.
func TestConnect_UnsupportedConnectionType(t *testing.T) {
	cfg := &Config{
		User:           "user",
		Password:       "password",
		Database:       "database",
		ConnectionType: "unsupported",
		Driver:         "mysql",
	}

	dbFactory := func(driver string, dsn string) (util.DB, error) {
		assert.Fail(t, "should not be called")
		return nil, nil
	}

	db, err := Connect(cfg, dbFactory)

	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Equal(t, "unsupported connection type: unsupported", err.Error())
}
