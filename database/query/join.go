package query

import (
	"fmt"

	"github.com/pakkasys/fluidapi/database/clause"
)

// ColumnSelectorToString returns the string representation of the ColumnnSelector.
func ColumnSelectorToString(columnnSelector clause.ColumnSelector) string {
	return fmt.Sprintf(
		"`%s`.`%s`",
		columnnSelector.Table,
		columnnSelector.Columnn,
	)
}
