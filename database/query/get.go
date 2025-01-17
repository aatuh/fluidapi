package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/internal"
	util "github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/page"
)

type RowScanner[T any] func(row util.Row, entity *T) error
type RowScannerMultiple[T any] func(rows util.Rows, entity *T) error

// GetOptions is the options struct used for queries.
type GetOptions struct {
	Selectors   []util.Selector
	Orders      []util.Order
	Page        *page.Page
	Joins       []util.Join
	Projections []util.Projection
	Lock        bool
}

func Get[T any](
	tableName string,
	rowScanner RowScanner[T],
	preparer util.Preparer,
	dbOptions *GetOptions,
) (*T, error) {
	query, whereValues := buildBaseGetQuery(tableName, dbOptions)

	return GetWithQuery(
		tableName,
		rowScanner,
		preparer,
		query,
		whereValues,
	)
}

func GetWithQuery[T any](
	tableName string,
	rowScanner RowScanner[T],
	preparer util.Preparer,
	query string,
	params []any,
) (*T, error) {
	entity, err := querySingle(
		preparer,
		query,
		params,
		rowScanner,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return entity, nil
}

func GetMany[T any](
	tableName string,
	rowScannerMultiple RowScannerMultiple[T],
	preparer util.Preparer,
	dbOptions *GetOptions,
) ([]T, error) {
	query, whereValues := buildBaseGetQuery(tableName, dbOptions)

	ent, e := GetManyWithQuery(
		tableName,
		rowScannerMultiple,
		preparer,
		query,
		whereValues,
	)

	return ent, e
}

func GetManyWithQuery[T any](
	tableName string,
	rowScannerMultiple RowScannerMultiple[T],
	preparer util.Preparer,
	query string,
	params []any,
) ([]T, error) {
	entities, err := queryMultiple(preparer, query, params, rowScannerMultiple)
	if err != nil {
		if err == sql.ErrNoRows {
			return []T{}, nil
		}
		return nil, err
	}
	return entities, nil
}

func queryMultiple[T any](
	preparer util.Preparer,
	query string,
	params []any,
	rowScannerMultiple RowScannerMultiple[T],
) ([]T, error) {
	rows, statement, err := Query(preparer, query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer statement.Close()

	entities, err := rowsToEntities(rows, rowScannerMultiple)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func querySingle[T any](
	preparer util.Preparer,
	query string,
	params []any,
	rowScanner RowScanner[T],
) (*T, error) {
	statement, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var entity T

	row := statement.QueryRow(params...)
	err = rowScanner(row, &entity)
	if err != nil {
		return nil, err
	}
	if err := row.Err(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func projectionsToStrings(projections []util.Projection) []string {
	if len(projections) == 0 {
		return []string{"*"}
	}

	projectionStrings := make([]string, len(projections))
	for i, projection := range projections {
		projectionStrings[i] = projection.String()
	}
	return projectionStrings
}

func joinClause(joins []util.Join) string {
	var joinClause string
	for _, join := range joins {
		if joinClause != "" {
			joinClause += " "
		}
		joinClause += fmt.Sprintf(
			"%s JOIN `%s` ON %s = %s",
			join.Type,
			join.Table,
			join.OnLeft.String(),
			join.OnRight.String(),
		)
	}
	return joinClause
}

func whereClause(selectors []util.Selector) (string, []any) {
	whereColumns, whereValues := internal.ProcessSelectors(selectors)

	var whereClause string
	if len(whereColumns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumns, " AND ")
	}

	return strings.Trim(whereClause, " "), whereValues
}

func buildBaseGetQuery(
	tableName string,
	dbOptions *GetOptions,
) (string, []any) {
	whereClause, whereValues := whereClause(dbOptions.Selectors)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(
		"SELECT %s",
		strings.Join(projectionsToStrings(dbOptions.Projections), ","),
	))
	builder.WriteString(fmt.Sprintf(" FROM `%s`", tableName))
	if len(dbOptions.Joins) != 0 {
		builder.WriteString(" " + joinClause(dbOptions.Joins))
	}
	if whereClause != "" {
		builder.WriteString(" " + whereClause)
	}
	if len(dbOptions.Orders) != 0 {
		builder.WriteString(" " + getOrderClauseFromOrders(dbOptions.Orders))
	}
	if dbOptions.Page != nil {
		builder.WriteString(" " + getLimitOffsetClauseFromPage(dbOptions.Page))
	}
	if dbOptions.Lock {
		builder.WriteString(" FOR UPDATE")
	}

	return builder.String(), whereValues
}

func getLimitOffsetClauseFromPage(page *page.Page) string {
	if page == nil {
		return ""
	}

	return fmt.Sprintf(
		"LIMIT %d OFFSET %d",
		page.Limit,
		page.Offset,
	)
}

func getOrderClauseFromOrders(orders []util.Order) string {
	if len(orders) == 0 {
		return ""
	}

	orderClause := "ORDER BY"
	for _, readOrder := range orders {
		if readOrder.Table == "" {
			orderClause += fmt.Sprintf(
				" `%s` %s,",
				readOrder.Field,
				readOrder.Direction,
			)
		} else {
			orderClause += fmt.Sprintf(
				" `%s`.`%s` %s,",
				readOrder.Table,
				readOrder.Field,
				readOrder.Direction,
			)
		}
	}

	return strings.TrimSuffix(orderClause, ",")
}

func rowsToEntities[T any](
	rows util.Rows,
	rowScannerMultiple RowScannerMultiple[T],
) ([]T, error) {
	if rowScannerMultiple == nil {
		return nil, fmt.Errorf("must provide rowScannerMultiple")
	}

	var entities []T
	for rows.Next() {
		var entity T
		err := rowScannerMultiple(rows, &entity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}
