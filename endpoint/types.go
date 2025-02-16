package endpoint

import (
	"fmt"
	"slices"

	"github.com/pakkasys/fluidapi/core"
	"github.com/pakkasys/fluidapi/database"
)

type Predicate string
type Predicates []Predicate

const (
	GREATER                Predicate = ">"
	GREATER_SHORT          Predicate = "GT"
	GREATER_OR_EQUAL       Predicate = ">="
	GREATER_OR_EQUAL_SHORT Predicate = "GE"
	EQUAL                  Predicate = "="
	EQUAL_SHORT            Predicate = "EQ"
	NOT_EQUAL              Predicate = "!="
	NOT_EQUAL_SHORT        Predicate = "NE"
	LESS                   Predicate = "<"
	LESS_SHORT             Predicate = "LT"
	LESS_OR_EQUAL          Predicate = "<="
	LESS_OR_EQUAL_SHORT    Predicate = "LE"
	IN                     Predicate = "IN"
	NOT_IN                 Predicate = "NOT_IN"
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

		endpoints = append(
			endpoints,
			core.Endpoint{
				URL:         definition.URL,
				Method:      definition.Method,
				Middlewares: middlewares,
			},
		)
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
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=0"`
}

// Validate validates the input page.
func (p *Page) Validate(maxLimit int) error {
	if p.Limit > maxLimit {
		return MaxPageLimitExceededError.WithData(
			MaxPageLimitExceededErrorData{
				MaxLimit: maxLimit,
			},
		)
	}
	return nil
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

const (
	DIRECTION_ASC        OrderDirection = "ASC"
	DIRECTION_ASCENDING  OrderDirection = "ASCENDING"
	DIRECTION_DESC       OrderDirection = "DESC"
	DIRECTION_DESCENDING OrderDirection = "DESCENDING"
)

// Directions is a list of all possible order directions.
var Directions []OrderDirection = []OrderDirection{
	DIRECTION_ASC,
	DIRECTION_ASCENDING,
	DIRECTION_DESC,
	DIRECTION_DESCENDING,
}

// DirectionDatabaseTranslations is a map of order directions to database
// order directions.
var DirectionDatabaseTranslations = map[OrderDirection]database.OrderDirection{
	DIRECTION_ASC:        database.OrderAsc,
	DIRECTION_ASCENDING:  database.OrderAsc,
	DIRECTION_DESC:       database.OrderDesc,
	DIRECTION_DESCENDING: database.OrderDesc,
}

// Order is used to specify the order of the result set.
type Order struct {
	Field     string         `json:"field"`
	Direction OrderDirection `json:"direction"`
}

type Orders []Order

// ValidateAndTranslateToDBOrders validates and translates the provided
// orders into database orders.
// It also returns an error if any of the orders are invalid.
//
//   - orders: The list of orders to validate and translate.
//   - allowedOrderFields: The list of allowed order fields.
//   - apiToDBFieldMap: The mapping of API field names to database field names.
func (o Orders) ValidateAndTranslateToDBOrders(
	allowedOrderFields []string,
	apiToDBFieldMap map[string]DBField,
) ([]database.Order, error) {
	newOrders, err := o.ValidateAndDeduplicateOrders(allowedOrderFields)
	if err != nil {
		return nil, err
	}

	dbOrders, err := newOrders.ToDBOrders(apiToDBFieldMap)
	if err != nil {
		return nil, err
	}

	return dbOrders, nil
}

// ValidateAndDeduplicateOrders validates and deduplicates the provided orders.
// It returns a new list of orders with no duplicates.
// It also returns an error if any of the orders are invalid.
//
//   - allowedOrderFields: The list of allowed order fields.
func (o Orders) ValidateAndDeduplicateOrders(
	allowedOrderFields []string,
) (Orders, error) {
	newOrders := []Order{}
	addedFields := make(map[string]bool)

	for i := range o {
		order := o[i]

		if err := order.validate(allowedOrderFields); err != nil {
			return nil, err
		}

		if !addedFields[order.Field] {
			newOrders = append(newOrders, order)
			addedFields[order.Field] = true
		}
	}

	return newOrders, nil
}

// ToDBOrders translates the provided orders into database orders.
// It returns an error if any of the orders are invalid.
//
//   - apiToDBFieldMap: The mapping of API field names to database field names.
func (o Orders) ToDBOrders(
	apiToDBFieldMap map[string]DBField,
) ([]database.Order, error) {
	newOrders := []database.Order{}

	for i := range o {
		order := o[i]

		translatedField := apiToDBFieldMap[order.Field]

		// Translate column
		dbColumn := translatedField.Column
		if dbColumn == "" {
			return nil, InvalidOrderFieldError.WithData(
				InvalidOrderFieldErrorData{
					Field: order.Field,
				},
			)
		}
		order.Field = dbColumn

		newOrders = append(
			newOrders,
			database.Order{
				Table:     translatedField.Table,
				Field:     order.Field,
				Direction: DirectionDatabaseTranslations[order.Direction],
			},
		)
	}

	return newOrders, nil
}

func (o *Order) validate(allowedOrderFields []string) error {
	// Check that the order direction is valid
	if !slices.Contains(Directions, o.Direction) {
		return InvalidOrderFieldError.WithData(
			InvalidOrderFieldErrorData{
				Field: o.Field,
			},
		)
	}

	// Check that the order field is allowed
	if !slices.Contains(allowedOrderFields, o.Field) {
		return InvalidOrderFieldError.WithData(
			InvalidOrderFieldErrorData{
				Field: o.Field,
			},
		)
	}

	return nil
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
	// Predicates allowed for this selector
	AllowedPredicates []Predicate
	// The name of the field being filtered
	Field string
	// The predicate for filtering
	Predicate Predicate
	// The value to filter by
	Value any
}

// String returns a string representation of the selector.
// It is useful for debugging and logging purposes.
//
// Returns:
// - A formatted string showing the field, predicate, and value.
func (s Selector) String() string {
	return fmt.Sprintf("%s %s %v", s.Field, s.Predicate, s.Value)
}

// Selectors represents a collection of selectors used for filtering data.
type Selectors []Selector

// GetByFields returns all selectors that match the given fields.
//
// Parameters:
// - fields: The fields to search for in the selectors.
//
// Returns:
// - A slice of selectors that match the provided field names.
func (s Selectors) GetByFields(fields ...string) []Selector {
	selectors := Selectors{}
	for f := range fields {
		for j := range s {
			if s[j].Field == fields[f] {
				selectors = append(selectors, s[j])
			}
		}
	}
	return selectors
}

// ToDBSelectors converts a slice of API-level selectors to database selectors.
// It validates predicates and translates the fields and predicates for use with
// the database.
//
// Parameters:
//   - selectors: A slice of API-level selectors that specify the criteria for
//     selecting data.
//   - apiToDBFieldMap: A map translating API field names to their corresponding
//     database field definitions.
//
// Returns:
//   - A slice of database.Selector, which represents the translated database
//     selectors.
//   - An error if any validation fails, such as invalid predicates or unknown
//     fields.
func (selectors Selectors) ToDBSelectors(
	apiToDBFieldMap map[string]DBField,
) ([]database.Selector, error) {
	var databaseSelectors []database.Selector

	for i := range selectors {
		selector := selectors[i]

		// Validate the input predicate
		if !slices.Contains(
			selector.AllowedPredicates,
			selector.Predicate,
		) {
			return nil, PredicateNotAllowedError.WithData(
				PredicateNotAllowedErrorData{Predicate: selector.Predicate},
			)
		}

		// Translate the predicate
		dbPredicate, ok := ToDBPredicates[selector.Predicate]
		if !ok {
			return nil, InvalidPredicateError.WithData(
				InvalidPredicateErrorData{Predicate: selector.Predicate},
			)
		}

		// Translate the field
		dbField, ok := apiToDBFieldMap[selector.Field]
		if !ok {
			return nil, InvalidDatabaseSelectorTranslationError.WithData(
				InvalidDatabaseSelectorTranslationErrorData{
					Field: selector.Field,
				},
			)
		}

		databaseSelectors = append(
			databaseSelectors,
			database.Selector{
				Table:     dbField.Table,
				Field:     dbField.Column,
				Predicate: dbPredicate,
				Value:     selector.Value,
			},
		)
	}

	return databaseSelectors, nil
}

type InvalidDatabaseUpdateTranslationErrorData struct {
	Field string `json:"field"`
}

var InvalidDatabaseUpdateTranslationError = core.NewAPIError("INVALID_DATABASE_UPDATE_TRANSLATION")

// Update represents a data update with a field and a value.
type Update struct {
	Field string // The field to be updated
	Value any    // The new value for the field
}

type Updates []Update

// GetByField returns update with the given field.
//
//   - field: the field to search for
func (u Updates) GetByField(field string) *Update {
	for j := range u {
		if u[j].Field == field {
			return &u[j]
		}
	}
	return nil
}

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

	for i := range updates {
		matchedUpdate := updates[i]

		// Translate the field
		dbField, ok := apiToDBFieldMap[matchedUpdate.Field]
		if !ok {
			return nil, InvalidDatabaseUpdateTranslationError.WithData(
				InvalidDatabaseUpdateTranslationErrorData{
					Field: matchedUpdate.Field,
				},
			)
		}

		dbUpdates = append(
			dbUpdates,
			database.UpdateField{
				Field: dbField.Column,
				Value: matchedUpdate.Value,
			},
		)
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
					append(
						[]Wrapper{wrapper},
						(*s)[i+1:]...,
					)...,
				)
			}
			return true
		}
	}
	return false
}
