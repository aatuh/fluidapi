package query

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/clause"
)

// ProjectionToString returns the string representation of a projection.
func ProjectionToString(projection clause.Projection) string {
	builder := strings.Builder{}

	if projection.Table == "" {
		builder.WriteString(fmt.Sprintf("`%s`", projection.Column))
	} else {
		builder.WriteString(fmt.Sprintf(
			"`%s`.`%s`",
			projection.Table,
			projection.Column,
		))
	}

	if projection.Alias != "" {
		builder.WriteString(fmt.Sprintf(" AS `%s`", projection.Alias))
	}

	return builder.String()
}
