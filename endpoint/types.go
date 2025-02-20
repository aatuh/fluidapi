package endpoint

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/core"
	"github.com/pakkasys/fluidapi/database"
)

type Predicate string

func (p Predicate) String() string {
	return string(p)
}

type Predicates []Predicate

func (p Predicates) String() string {
	str := make([]string, len(p))
	for i, predicate := range p {
		str[i] = predicate.String()
	}
	return strings.Join(str, ",")
}

const (
	GREATER                Predicate = ">"
	GREATER_SHORT          Predicate = "gt"
	GREATER_OR_EQUAL       Predicate = ">="
	GREATER_OR_EQUAL_SHORT Predicate = "ge"
	EQUAL                  Predicate = "="
	EQUAL_SHORT            Predicate = "eq"
	NOT_EQUAL              Predicate = "!="
	NOT_EQUAL_SHORT        Predicate = "ne"
	LESS                   Predicate = "<"
	LESS_SHORT             Predicate = "LT"
	LESS_OR_EQUAL          Predicate = "<="
	LESS_OR_EQUAL_SHORT    Predicate = "le"
	IN                     Predicate = "in"
	NOT_IN                 Predicate = "not_in"
)

var ToDBPredicates = map[Predicate]database.Predicate{
	GREATER:                database.GREATER,
	GREATER_SHORT:          database.GREATER,
	GREATER_OR_EQUAL:       database.GREATER_OR_EQUAL,
	GREATER_OR_EQUAL_SHORT: database.GREATER_OR_EQUAL,
	EQUAL:                  database.EQUAL,
	EQUAL_SHORT:            database.EQUAL,
	NOT_EQUAL:              database.NOT_EQUAL,
	NOT_EQUAL_SHORT:        database.NOT_EQUAL,
	LESS:                   database.LESS,
	LESS_SHORT:             database.LESS,
	LESS_OR_EQUAL:          database.LESS_OR_EQUAL,
	LESS_OR_EQUAL_SHORT:    database.LESS_OR_EQUAL,
	IN:                     database.IN,
	NOT_IN:                 database.NOT_IN,
}

var AllPredicates = []Predicate{
	GREATER,
	GREATER_SHORT,
	GREATER_OR_EQUAL,
	GREATER_OR_EQUAL_SHORT,
	EQUAL,
	EQUAL_SHORT,
	NOT_EQUAL,
	NOT_EQUAL_SHORT,
	LESS,
	LESS_SHORT,
	LESS_OR_EQUAL,
	LESS_OR_EQUAL_SHORT,
	IN,
	NOT_IN,
}

var OnlyEqualPredicates = []Predicate{
	EQUAL,
	EQUAL_SHORT,
}

var EqualAndNotEqualPredicates = []Predicate{
	EQUAL,
	EQUAL_SHORT,
	NOT_EQUAL,
	NOT_EQUAL_SHORT,
}

var OnlyGreaterPredicates = []Predicate{
	GREATER_OR_EQUAL,
	GREATER_OR_EQUAL_SHORT,
	GREATER,
	GREATER_SHORT,
}

var OnlyLessPredicates = []Predicate{
	LESS_OR_EQUAL,
	LESS_OR_EQUAL_SHORT,
	LESS,
	LESS_SHORT,
}

var OnlyInAndNotInPredicates = []Predicate{
	IN,
	NOT_IN,
}

// Definition represents an endpoint definition.
type Definition struct {
	URL    string
	Method string
	Stack  Stack
}

// Option is a function that modifies a definition when it is cloned
type Option func(*Definition)

// Clone clones an endpoint definition with options
func (d *Definition) Clone(options ...Option) *Definition {
	cloned := *d
	for _, option := range options {
		option(&cloned)
	}
	return &cloned
}

// WithURL returns an option that sets the URL of the endpoint
func WithURL(url string) Option {
	return func(e *Definition) {
		e.URL = url
	}
}

// WithMethod returns an option that sets the method of the endpoint
func WithMethod(method string) Option {
	return func(e *Definition) {
		e.Method = method
	}
}

// WithMiddlewareStack return an option that sets the middleware stack
func WithMiddlewareStack(stack Stack) Option {
	return func(e *Definition) {
		e.Stack = stack
	}
}

// WithMiddlewareWrappersFunc returns an option that sets the middleware stack
func WithMiddlewareWrappersFunc(
	middlewareWrappersFunc func(definition *Definition) Stack,
) Option {
	return func(e *Definition) {
		e.Stack = middlewareWrappersFunc(e)
	}
}

type Definitions []Definition

// ToEndpoints converts a list of endpoint definitions to a list of API
// endpoints.
func (d Definitions) ToEndpoints() []core.Endpoint {
	endpoints := []core.Endpoint{}

	for _, definition := range d {
		middlewares := []core.Middleware{}
		for _, mw := range definition.Stack {
			middlewares = append(middlewares, mw.Middleware)
		}

		endpoints = append(endpoints, core.Endpoint{
			URL:         definition.URL,
			Method:      definition.Method,
			Middlewares: middlewares,
		})
	}

	return endpoints
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

// Wrapper wraps a middleware function with additional metadata.
type Wrapper struct {
	Middleware core.Middleware
	ID         string
	Data       any
}

// Stack represents a list of middleware wrappers.
type Stack []Wrapper

// Middlewares returns the middlewares in the stack.
//
// Parameters:
//   - s: The middleware stack.
//
// Returns:
//   - The middlewares in the stack.
func (s Stack) Middlewares() []core.Middleware {
	middlewares := []core.Middleware{}
	for _, mw := range s {
		middlewares = append(middlewares, mw.Middleware)
	}
	return middlewares
}

// InsertAfterID inserts a middleware wrapper after the given ID.
//
// Parameters:
//   - id: The ID of the middleware to insert after.
//   - wrapper: The middleware wrapper to insert.
//
// Returns:
//   - True if the middleware was inserted, false otherwise.
func (s *Stack) InsertAfterID(id string, wrapper Wrapper) bool {
	for i, mw := range *s {
		if mw.ID == id {
			if i == len(*s)-1 {
				*s = append(*s, wrapper)
			} else {
				*s = append(
					(*s)[:i+1],
					append([]Wrapper{wrapper}, (*s)[i+1:]...)...,
				)
			}
			return true
		}
	}
	return false
}
