package database

// Predicate represents the predicate of a database selector.
type Predicate string

const (
	Greater        Predicate = ">"
	GreaterOrEqual Predicate = ">="
	Equal          Predicate = "="
	NotEqual       Predicate = "!="
	Less           Predicate = "<"
	LessOrEqual    Predicate = "<="
	In             Predicate = "IN"
	NotIn          Predicate = "NOT IN"
)

// OrderDirection is used to specify the order of the result set.
type OrderDirection string

// Order directions.
const (
	OrderAsc  OrderDirection = "ASC"
	OrderDesc OrderDirection = "DESC"
)

// Order is used to specify the order of the result set.
type Order struct {
	Table     string
	Field     string
	Direction OrderDirection
}

// Orders is a list of orders
type Orders []Order

// ColumnSelector represents a column selector
type ColumnSelector struct {
	Table  string
	Column string
}

// Projection represents a projected column in a query.
type Projection struct {
	Table  string
	Column string
	Alias  string
}

// Projections is a list of projections
type Projections []Projection

// Selector represents a database selector.
type Selector struct {
	Table     string
	Column    string
	Predicate Predicate
	Value     any
}

// Selectors represents a list of database selectors.
type Selectors []Selector

// GetByField returns selector with the given field.
//
// Parameters:
//   - field: the field to search for
//
// Returns:
//   - *Selector: The selector
func (s Selectors) GetByField(field string) *Selector {
	for j := range s {
		if s[j].Column == field {
			return &s[j]
		}
	}
	return nil
}

// GetByFields returns selectors with the given fields.
//
// Parameters:
//   - fields: the fields to search for
//
// Returns:
//   - []Selector: A list of selectors
func (s Selectors) GetByFields(fields ...string) []Selector {
	var result []Selector
	for _, field := range fields {
		for i := range s {
			if s[i].Column == field {
				result = append(result, s[i])
			}
		}
	}
	return result
}

// UpdateField is the options struct used for update queries.
type UpdateField struct {
	Field string
	Value any
}

// UpdateFields is a list of update fields
type UpdateFields []UpdateField

// Page is used to specify the page of the result set.
type Page struct {
	Offset int
	Limit  int
}

// JoinType represents the type of join
type JoinType string

// Join types.
const (
	JoinTypeInner JoinType = "INNER"
	JoinTypeLeft  JoinType = "LEFT"
	JoinTypeRight JoinType = "RIGHT"
	JoinTypeFull  JoinType = "FULL"
)

// Join represents a database join clause.
type Join struct {
	Type    JoinType
	Table   string
	OnLeft  ColumnSelector
	OnRight ColumnSelector
}

// Joins is a list of joins
type Joins []Join
