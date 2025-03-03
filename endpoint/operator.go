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

const (
	GREATER          Predicate = ">"
	GT               Predicate = "gt"
	GREATER_OR_EQUAL Predicate = ">="
	GE               Predicate = "ge"
	EQUAL            Predicate = "="
	EQ               Predicate = "eq"
	NOT_EQUAL        Predicate = "!="
	NE               Predicate = "ne"
	LESS             Predicate = "<"
	LT               Predicate = "LT"
	LESS_OR_EQUAL    Predicate = "<="
	LE               Predicate = "le"
	IN               Predicate = "in"
	NOT_IN           Predicate = "not_in"
)

// ToDBPredicates maps API-level predicates to database predicates.
var ToDBPredicates = map[Predicate]database.Predicate{
	GREATER:          database.GREATER,
	GT:               database.GREATER,
	GREATER_OR_EQUAL: database.GREATER_OR_EQUAL,
	GE:               database.GREATER_OR_EQUAL,
	EQUAL:            database.EQUAL,
	EQ:               database.EQUAL,
	NOT_EQUAL:        database.NOT_EQUAL,
	NE:               database.NOT_EQUAL,
	LESS:             database.LESS,
	LT:               database.LESS,
	LESS_OR_EQUAL:    database.LESS_OR_EQUAL,
	LE:               database.LESS_OR_EQUAL,
	IN:               database.IN,
	NOT_IN:           database.NOT_IN,
}

// AllPredicates is a slice of all available predicates.
var AllPredicates = []Predicate{
	GREATER,
	GT,
	GREATER_OR_EQUAL,
	GE,
	EQUAL,
	EQ,
	NOT_EQUAL,
	NE,
	LESS,
	LT,
	LESS_OR_EQUAL,
	LE,
	IN,
	NOT_IN,
}

// OnlyEqualPredicates is a slice of predicates that only allow equality.
var OnlyEqualPredicates = []Predicate{
	EQUAL,
	EQ,
}

// EqualAndNotEqualPredicates is a slice of predicates that allow both equality
// and inequality.
var EqualAndNotEqualPredicates = []Predicate{
	EQUAL,
	EQ,
	NOT_EQUAL,
	NE,
}

// OnlyGreaterPredicates is a slice of predicates that only allow greater
// values.
var OnlyGreaterPredicates = []Predicate{
	GREATER_OR_EQUAL,
	GE,
	GREATER,
	GT,
}

// OnlyLessPredicates is a slice of predicates that only allow less values.
var OnlyLessPredicates = []Predicate{
	LESS_OR_EQUAL,
	LE,
	LESS,
	LT,
}

// OnlyInAndNotInPredicates is a slice of predicates that only allow
// IN and NOT_IN
var OnlyInAndNotInPredicates = []Predicate{
	IN,
	NOT_IN,
}

// DBField is used to translate between API field and database field.
type DBField struct {
	Table  string
	Column string
}

type MaxPageLimitExceededErrorData struct {
	MaxLimit int `json:"max_limit"`
}

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

type InvalidOrderFieldErrorData struct {
	Field string `json:"field"`
}

var InvalidOrderFieldError = core.NewAPIError("INVALID_ORDER_FIELD")

// OrderDirection is used to specify the direction of the order.
type OrderDirection string

// String returns the string representation of the order direction.
func (o OrderDirection) String() string {
	return string(o)
}

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
	newOrders, err := o.Dedup()
	if err != nil {
		return nil, err
	}

	dbOrders, err := newOrders.ToDBOrders(apiToDBFieldMap)
	if err != nil {
		return nil, err
	}

	return dbOrders, nil
}

// Dedup deduplicates the provided orders.
// It returns a new list of orders with no duplicates.
// It also returns an error if any of the orders are invalid.
func (o Orders) Dedup() (Orders, error) {
	dedup := map[string]OrderDirection{}
	existing := make(map[string]bool)

	for field := range o {
		order := o[field]
		if !existing[field] {
			dedup[field] = order
			existing[field] = true
		}
	}

	return dedup, nil
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

type InvalidDatabaseSelectorTranslationErrorData struct {
	Field string `json:"field"`
}

var InvalidDatabaseSelectorTranslationError = core.NewAPIError("INVALID_DATABASE_SELECTOR_TRANSLATION")

type InvalidPredicateErrorData struct {
	Predicate Predicate `json:""`
}

var InvalidPredicateError = core.NewAPIError("INVALID_PREDICATE")

type InvalidSelectorFieldErrorData struct {
	Field string `json:"field"`
}

var InvalidSelectorFieldError = core.NewAPIError("INVALID_SELECTOR_FIELD")

type PredicateNotAllowedErrorData struct {
	Predicate Predicate `json:"predicate"`
}

var PredicateNotAllowedError = core.NewAPIError("PREDICATE_NOT_ALLOWED")

// Selector represents a data selector that specifies criteria for filtering
// data based on fields, predicates, and values.
type Selector struct {
	Predicate Predicate `json:"predicate"` // The predicate to use.
	Value     any       `json:"value"`     // The value to filter on.
}

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
// - A formatted string showing the field, predicate, and value.
func (s Selector) String() string {
	return fmt.Sprintf("%s %v", s.Predicate, s.Value)
}

// Selectors represents a collection of selectors used for filtering data.
// It is a map where the key is the field name and the value is the selector.
type Selectors map[string]Selector

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
			Field:     dbField.Column,
			Predicate: dbPredicate,
			Value:     selector.Value,
		})
	}

	return databaseSelectors, nil
}

type InvalidDatabaseUpdateTranslationErrorData struct {
	Field string `json:"field"`
}

var InvalidDatabaseUpdateTranslationError = core.NewAPIError("INVALID_DATABASE_UPDATE_TRANSLATION")

// Updates represents a list of updates to apply to a database entity.
type Updates map[string]any

// ToDBUpdates translates a list of updates to a database update list
// and returns an error if the translation fails.
//
// Parameters:
// - updates: The list of updates to translate.
// - apiToDBFieldMap: The mapping of API field names to database field names.
//
// Returns:
// - A list of database entity updates.
// - An error if any field translation fails.
func (updates Updates) ToDBUpdates(
	apiToDBFieldMap map[string]DBField,
) ([]database.UpdateField, error) {
	var dbUpdates []database.UpdateField

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

		dbUpdates = append(dbUpdates, database.UpdateField{
			Field: dbField.Column,
			Value: value,
		})
	}

	return dbUpdates, nil
}
