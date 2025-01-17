package query

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/util"
)

// UpsertMany creates an upsert query for a list of entities.
//
//   - entities: The entities to upsert.
//   - tableName: The name of the database table.
//   - insertedValue: The function used to get the columns and values to insert.
//   - updateProjections: The projections of the entities to update.
func UpsertMany[T any](
	entities []*T,
	tableName string,
	insertedValue InsertedValues[*T],
	updateProjections []util.Projection,
) (string, []any) {
	if len(entities) == 0 {
		return "", nil
	}

	updateParts := make([]string, len(updateProjections))
	for i, proj := range updateProjections {
		updateParts[i] = fmt.Sprintf(
			"`%s` = VALUES(`%s`)",
			proj.Column,
			proj.Column,
		)
	}

	insertQueryPart, allValues := InsertMany(entities, tableName, insertedValue)

	builder := strings.Builder{}
	builder.WriteString(insertQueryPart)
	if len(updateParts) != 0 {
		builder.WriteString(" ON DUPLICATE KEY UPDATE ")
		builder.WriteString(strings.Join(updateParts, ", "))
	}
	upsertQuery := builder.String()

	return upsertQuery, allValues
}
