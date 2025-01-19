package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pakkasys/fluidapi/database/clause"
)

const (
	in    = "IN"
	is    = "IS"
	isNot = "IS NOT"
	null  = "NULL"
)

// ProcessSelectors processes selectors and returns columns and values.
//
//   - selectors: the selectors to process
func ProcessSelectors(selectors []clause.Selector) ([]string, []any) {
	var whereColumns []string
	var whereValues []any
	for _, selector := range selectors {
		col, vals := processSelector(selector)
		whereColumns = append(whereColumns, col)
		whereValues = append(whereValues, vals...)
	}
	return whereColumns, whereValues
}

func processSelector(selector clause.Selector) (string, []any) {
	if selector.Predicate == in {
		return processInSelector(selector)
	}
	return processDefaultSelector(selector)
}

func processInSelector(selector clause.Selector) (string, []any) {
	value := reflect.ValueOf(selector.Value)
	if value.Kind() == reflect.Slice {
		placeholders, values := createPlaceholdersAndValues(value)
		column := fmt.Sprintf(
			"`%s`.`%s` %s (%s)",
			selector.Table,
			selector.Field,
			in,
			placeholders,
		)
		return column, values
	}
	// If value is not a slice, treat as a single value
	return fmt.Sprintf(
		"`%s`.`%s` %s (?)",
		selector.Table,
		selector.Field,
		in,
	), []any{selector.Value}
}

func processDefaultSelector(selector clause.Selector) (string, []any) {
	if selector.Value == nil {
		return processNullSelector(selector)
	}
	if selector.Table == "" {
		return fmt.Sprintf(
			"`%s` %s ?",
			selector.Field,
			selector.Predicate,
		), []any{selector.Value}
	} else {
		return fmt.Sprintf(
			"`%s`.`%s` %s ?",
			selector.Table,
			selector.Field,
			selector.Predicate,
		), []any{selector.Value}
	}
}

func processNullSelector(selector clause.Selector) (string, []any) {
	if selector.Predicate == "=" {
		return buildNullClause(selector, is), nil
	}
	if selector.Predicate == "!=" {
		return buildNullClause(selector, isNot), nil
	}
	return "", nil
}

func buildNullClause(selector clause.Selector, clause string) string {
	if selector.Table == "" {
		return fmt.Sprintf("`%s` %s %s", selector.Field, clause, null)
	}
	return fmt.Sprintf(
		"`%s`.`%s` %s %s",
		selector.Table,
		selector.Field,
		clause,
		null,
	)
}

func createPlaceholdersAndValues(value reflect.Value) (string, []any) {
	placeholderCount := value.Len()
	placeholders := createPlaceholders(placeholderCount)
	values := make([]any, placeholderCount)
	for i := 0; i < placeholderCount; i++ {
		values[i] = value.Index(i).Interface()
	}
	return placeholders, values
}

func createPlaceholders(count int) string {
	return strings.TrimSuffix(strings.Repeat("?,", count), ",")
}
