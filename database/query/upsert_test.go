package query

import (
	"testing"

	"github.com/pakkasys/fluidapi/database/clause"
	"github.com/stretchr/testify/assert"
)

type TestUpsertEntity struct {
	ID   int
	Name string
	Age  int
}

// TestUpsertManyQuery_NormalOperation tests upsertManyQuery with multiple
// entities and projections.
func TestUpsertManyQuery_NormalOperation(t *testing.T) {
	// Test entities and projections
	entities := []*TestUpsertEntity{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
	}
	projections := []clause.Projection{
		{Column: "name", Alias: "test"},
		{Column: "age", Alias: "test"},
	}

	// Inserter function for the entities
	inserter := func(e *TestUpsertEntity) ([]string, []any) {
		return []string{"id", "name", "age"}, []any{e.ID, e.Name, e.Age}
	}

	query, values := UpsertMany(entities, "user", inserter, projections)

	expectedQuery := "INSERT INTO `user` (`id`, `name`, `age`) VALUES (?, ?, ?), (?, ?, ?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `age` = VALUES(`age`)"
	expectedValues := []any{1, "Alice", 30, 2, "Bob", 25}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpsertManyQuery_SingleEntity tests upsertManyQuery with a single entity.
func TestUpsertManyQuery_SingleEntity(t *testing.T) {
	// Test entity and projections
	entities := []*TestUpsertEntity{
		{ID: 1, Name: "Alice"},
	}
	projections := []clause.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function for the entity
	inserter := func(e *TestUpsertEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	query, values := UpsertMany(entities, "user", inserter, projections)

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`)"
	expectedValues := []any{1, "Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpsertManyQuery_EmptyEntities tests upsertManyQuery with no entities.
func TestUpsertManyQuery_EmptyEntities(t *testing.T) {
	// Test with no entities
	entities := []*TestUpsertEntity{}
	projections := []clause.Projection{
		{Column: "name", Alias: "test"},
	}

	// Inserter function (not used here since entities is empty)
	inserter := func(e *TestUpsertEntity) ([]string, []any) {
		return []string{}, []any{}
	}

	query, values := UpsertMany(entities, "user", inserter, projections)

	assert.Equal(t, "", query)
	assert.Equal(t, []any(nil), values)
}

// TestUpsertManyQuery_MissingUpdateProjections tests upsertManyQuery with no
// update projections.
func TestUpsertManyQuery_MissingUpdateProjections(t *testing.T) {
	// Test entities
	entities := []*TestUpsertEntity{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	// Inserter function for the entities
	inserter := func(e *TestUpsertEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{e.ID, e.Name}
	}

	// Call the function with an empty projection list
	query, values := UpsertMany(
		entities,
		"user",
		inserter,
		[]clause.Projection{},
	)

	assert.Equal(
		t,
		"INSERT INTO `user` (`id`, `name`) VALUES (?, ?), (?, ?)",
		query,
	)
	assert.Equal(t, []any{1, "Alice", 2, "Bob"}, values)
}
