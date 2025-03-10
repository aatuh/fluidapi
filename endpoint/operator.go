package endpoint

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/core"
	"github.com/pakkasys/fluidapi/database"
)

// Predicate is a string representation of a filtering predicate.
type Predicate string

// String returns the string representation of the predicate.
func (p Predicate) String() string {
	return string(p)
}

// Predicates is a slice of Predicate values.
type Predicates []Predicate

// String returns a string representation of the predicates.
func (p Predicates) String() string {
	str := make([]string, len(p))
	for i, predicate := range p {
		str[i] = predicate.String()
	}
	return strings.Join(str, ",")
}

// StrSlice returns a slice of strings representing the predicates.
func (p Predicates) StrSlice() []string {
	str := make([]string, len(p))
	for i, predicate := range p {
		str[i] = predicate.String()
	}
	return str
}

// Predicates for filtering data.
const (
	Greater        Predicate = ">"
	Gt             Predicate = "gt"
	GreaterOrEqual Predicate = ">="
	Ge             Predicate = "ge"
	Equal          Predicate = "="
	Eq             Predicate = "eq"
	NotEqual       Predicate = "!="
	Ne             Predicate = "ne"
	Less           Predicate = "<"
	Lt             Predicate = "LT"
	LessOrEqual    Predicate = "<="
	Le             Predicate = "le"
	In             Predicate = "in"
	NotIn          Predicate = "not_in"
)

// ToDBPredicates maps API-level predicates to database predicates.
var ToDBPredicates = map[Predicate]database.Predicate{
	Greater:        database.Greater,
	Gt:             database.Greater,
	GreaterOrEqual: database.GreaterOrEqual,
	Ge:             database.GreaterOrEqual,
	Equal:          database.Equal,
	Eq:             database.Equal,
	NotEqual:       database.NotEqual,
	Ne:             database.NotEqual,
	Less:           database.Less,
	Lt:             database.Less,
	LessOrEqual:    database.LessOrEqual,
	Le:             database.LessOrEqual,
	In:             database.In,
	NotIn:          database.NotIn,
}

// AllPredicates is a slice of all available predicates.
var AllPredicates = []Predicate{
	Greater,
	Gt,
	GreaterOrEqual,
	Ge,
	Equal,
	Eq,
	NotEqual,
	Ne,
	Less,
	Lt,
	LessOrEqual,
	Le,
	In,
	NotIn,
}

// OnlyEqualPredicates is a slice of predicates that only allow equality.
var OnlyEqualPredicates = []Predicate{
	Equal,
	Eq,
}

// EqualAndNotEqualPredicates is a slice of predicates that allow both equality
// and inequality.
var EqualAndNotEqualPredicates = []Predicate{
	Equal,
	Eq,
	NotEqual,
	Ne,
}

// OnlyGreaterPredicates is a slice of predicates that only allow greater
// values.
var OnlyGreaterPredicates = []Predicate{
	GreaterOrEqual,
	Ge,
	Greater,
	Gt,
}

// OnlyLessPredicates is a slice of predicates that only allow less values.
var OnlyLessPredicates = []Predicate{
	LessOrEqual,
	Le,
	Less,
	Lt,
}

// OnlyInAndNotInPredicates is a slice of predicates that only allow
// IN and NOT_IN
var OnlyInAndNotInPredicates = []Predicate{
	In,
	NotIn,
}

// DBField is used to translate between API field and database field.
type DBField struct {
	Table  string
	Column string
}

// MaxPageLimitExceededErrorData is the data for the MaxPageLimitExceededError
// error.
type MaxPageLimitExceededErrorData struct {
	MaxLimit int `json:"max_limit"`
}

// MaxPageLimitExceededError is returned when a page limit is exceeded.
var MaxPageLimitExceededError = core.NewAPIError("MAX_PAGE_LIMIT_EXCEEDED")

// Page represents a pagination input.
type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// ToDBPage converts a Page to database Page.
func (p *Page) ToDBPage() *database.Page {
	if p == nil {
		return nil
	}
	return &database.Page{
		Offset: p.Offset,
		Limit:  p.Limit,
	}
}

// InvalidOrderFieldErrorData is the data for the InvalidOrderFieldError error.
type InvalidOrderFieldErrorData struct {
	Field string `json:"field"`
}

// InvalidOrderFieldError is returned when a field is not allowed.
var InvalidOrderFieldError = core.NewAPIError("INVALID_ORDER_FIELD")

// OrderDirection is used to specify the direction of the order.
type OrderDirection string

// String returns the string representation of the order direction.
func (o OrderDirection) String() string {
	return string(o)
}

// Available order directions.
const (
	DirectionAsc        OrderDirection = "asc"
	DirectionAscending  OrderDirection = "ascending"
	DirectionDesc       OrderDirection = "desc"
	DirectionDescending OrderDirection = "descending"
)

// DirectionsToDB is a map of order directions to database order directions.
var DirectionsToDB = map[OrderDirection]database.OrderDirection{
	DirectionAsc:        database.OrderAsc,
	DirectionAscending:  database.OrderAsc,
	DirectionDesc:       database.OrderDesc,
	DirectionDescending: database.OrderDesc,
}

// Orders is a map of field names to order directions.
type Orders map[string]OrderDirection

// TranslateToDBOrders translates the provided orders into database orders.
// It also returns an error if any of the orders are invalid.
//
//   - orders: The list of orders to translate.
//   - allowedOrderFields: The list of allowed order fields.
//   - apiToDBFieldMap: The mapping of API field names to database field names.
func (o Orders) TranslateToDBOrders(
	apiToDBFieldMap map[string]DBField,
) ([]database.Order, error) {
	dbOrders, err := o.dedup().ToDBOrders(apiToDBFieldMap)
	if err != nil {
		return nil, fmt.Errorf("TranslateToDBOrders: %w", err)
	}
	return dbOrders, nil
}

// dedup deduplicates the provided orders.
func (o Orders) dedup() Orders {
	dedup := map[string]OrderDirection{}
	existing := make(map[string]bool)
	for field := range o {
		order := o[field]
		if !existing[field] {
			dedup[field] = order
			existing[field] = true
		}
	}
	return dedup
}

// ToDBOrders translates the provided orders into database orders.
// It returns an error if any of the orders are invalid.
//
//   - apiToDBFieldMap: The mapping of API field names to database field names.
func (o Orders) ToDBOrders(
	apiToDBFieldMap map[string]DBField,
) ([]database.Order, error) {
	dbOrders := []database.Order{}

	for field, direction := range o {
		translatedField := apiToDBFieldMap[field]

		// Translate field.
		dbColumn := translatedField.Column
		if dbColumn == "" {
			return nil, InvalidOrderFieldError.
				WithData(
					InvalidOrderFieldErrorData{Field: field},
				).
				WithMessage(fmt.Sprintf(
					"cannot translate field: %s", field,
				))
		}

		lowerDir := OrderDirection(strings.ToLower(string(direction)))
		dbOrders = append(dbOrders, database.Order{
			Table:     translatedField.Table,
			Field:     dbColumn,
			Direction: DirectionsToDB[lowerDir],
		})
	}

	return dbOrders, nil
}

// InvalidDatabaseSelectorTranslationErrorData is the data for the
// InvalidDatabaseSelectorTranslationError error.
type InvalidDatabaseSelectorTranslationErrorData struct {
	Field string `json:"field"`
}

// InvalidDatabaseSelectorTranslationError is returned when a field is not
// allowed.
var InvalidDatabaseSelectorTranslationError = core.NewAPIError("INVALID_DATABASE_SELECTOR_TRANSLATION")

// InvalidPredicateErrorData is the data for the InvalidPredicateError error.
type InvalidPredicateErrorData struct {
	Predicate Predicate `json:""`
}

// InvalidPredicateError is returned when a predicate is not allowed.
var InvalidPredicateError = core.NewAPIError("INVALID_PREDICATE")

// InvalidSelectorFieldErrorData is the data for the InvalidSelectorFieldError
// error.
type InvalidSelectorFieldErrorData struct {
	Field string `json:"field"`
}

// InvalidSelectorFieldError is returned when a field is not allowed.
var InvalidSelectorFieldError = core.NewAPIError("INVALID_SELECTOR_FIELD")

// PredicateNotAllowedErrorData is the data for the PredicateNotAllowedError
// error.
type PredicateNotAllowedErrorData struct {
	Predicate Predicate `json:"predicate"`
}

// PredicateNotAllowedError is returned when a predicate is not allowed.
var PredicateNotAllowedError = core.NewAPIError("PREDICATE_NOT_ALLOWED")

// Selector represents a data selector that specifies criteria for filtering
// data based on fields, predicates, and values.
type Selector struct {
	Predicate Predicate `json:"predicate"` // The predicate to use.
	Value     any       `json:"value"`     // The value to filter on.
}

// NewSelector creates a new selector with the provided predicate and value.
//
// Parameters:
//   - predicate: The predicate to use.
//   - value: The value to filter on.
//
// Returns:
//   - A new selector.
func NewSelector(predicate Predicate, value any) *Selector {
	return &Selector{
		Predicate: predicate,
		Value:     value,
	}
}

// String returns a string representation of the selector.
// It is useful for debugging and logging purposes.
//
// Returns:
//   - A formatted string showing the field, predicate, and value.
func (s Selector) String() string {
	return fmt.Sprintf("%s %v", s.Predicate, s.Value)
}

// Selectors represents a collection of selectors used for filtering data.
// It is a map where the key is the field name and the value is the selector.
type Selectors map[string]Selector

// AddSelector adds a new selector to the collection of selectors.
//
// Parameters:
//   - field: The field name.
//   - predicate: The predicate to use.
//   - value: The value to filter on.
//
// Returns:
//   - A new collection of selectors with the new selector added.
func (s Selectors) AddSelector(
	field string, predicate Predicate, value any,
) Selectors {
	s[field] = Selector{
		Predicate: predicate,
		Value:     value,
	}
	return s
}

// ToDBSelectors converts a slice of API-level selectors to database selectors.
//
// Selectors
// Parameters:
//   - apiToDBFieldMap: A map translating API field names to their corresponding
//     database field definitions.
//
// Returns:
//   - A slice of database.Selector, which represents the translated database
//     selectors.
//   - An error if any validation fails, such as invalid predicates or unknown
//     fields.
func (s Selectors) ToDBSelectors(
	apiToDBFieldMap map[string]DBField,
) ([]database.Selector, error) {
	var databaseSelectors []database.Selector

	for field := range s {
		selector := s[field]

		// Translate the predicate.
		lowerPredicate := Predicate(strings.ToLower(string(selector.Predicate)))
		dbPredicate, ok := ToDBPredicates[lowerPredicate]
		if !ok {
			return nil, InvalidPredicateError.
				WithData(
					InvalidPredicateErrorData{Predicate: selector.Predicate},
				).
				WithMessage(fmt.Sprintf(
					"cannot translate predicate: %s", selector.Predicate,
				))
		}

		// Translate the field.
		dbField, ok := apiToDBFieldMap[field]
		if !ok {
			return nil, InvalidDatabaseSelectorTranslationError.
				WithData(
					InvalidDatabaseSelectorTranslationErrorData{Field: field},
				).
				WithMessage(fmt.Sprintf(
					"cannot translate field: %s", field,
				))
		}

		databaseSelectors = append(databaseSelectors, database.Selector{
			Table:     dbField.Table,
			Column:    dbField.Column,
			Predicate: dbPredicate,
			Value:     selector.Value,
		})
	}

	return databaseSelectors, nil
}

// InvalidDatabaseUpdateTranslationErrorData is the data for the
// InvalidDatabaseUpdateTranslationError error.
type InvalidDatabaseUpdateTranslationErrorData struct {
	Field string `json:"field"`
}

// InvalidDatabaseUpdateTranslationError is used to indicate that the
// translation of a database update failed.
var InvalidDatabaseUpdateTranslationError = core.NewAPIError("INVALID_DATABASE_UPDATE_TRANSLATION")

// Updates represents a list of updates to apply to a database entity.
type Updates map[string]any

// ToDBUpdates translates a list of updates to a database update list
// and returns an error if the translation fails.
//
// Parameters:
//   - updates: The list of updates to translate.
//   - apiToDBFieldMap: The mapping of API field names to database field names.
//
// Returns:
//   - A list of database entity updates.
//   - An error if any field translation fails.
func (updates Updates) ToDBUpdates(
	apiToDBFieldMap map[string]DBField,
) ([]database.Update, error) {
	var dbUpdates []database.Update

	for field := range updates {
		value := updates[field]

		// Translate the field.
		dbField, ok := apiToDBFieldMap[field]
		if !ok {
			return nil, InvalidDatabaseUpdateTranslationError.
				WithData(
					InvalidDatabaseUpdateTranslationErrorData{Field: field},
				).
				WithMessage(fmt.Sprintf(
					"cannot translate field: %s", field,
				))
		}

		dbUpdates = append(dbUpdates, database.Update{
			Field: dbField.Column,
			Value: value,
		})
	}

	return dbUpdates, nil
}
