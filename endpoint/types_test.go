package endpoint

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pakkasys/fluidapi/core"
	"github.com/pakkasys/fluidapi/database"
	"github.com/stretchr/testify/assert"
)

// TestValidateAndTranslateToDBOrders tests the
// ValidateAndTranslateToDBOrders function
func TestValidateAndTranslateToDBOrders_ValidInput(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
		{Field: "age", Direction: DIRECTION_DESC},
	}

	allowedFields := []string{"name", "age", "email"}
	fieldTranslations := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ValidateAndTranslateToDBOrders(
		orders,
		allowedFields,
		fieldTranslations,
	)

	assert.Nil(t, err, "Expected no error for valid orders")
	assert.Len(t, dbOrders, 2, "Expected two translated orders")
	assert.Equal(t, "users", dbOrders[0].Table, "Expected correct table for name order")
	assert.Equal(t, "user_name", dbOrders[0].Field, "Expected correct column for name order")
	assert.Equal(t, database.OrderAsc, dbOrders[0].Direction, "Expected ascending order for name order")
	assert.Equal(t, "users", dbOrders[1].Table, "Expected correct table for age order")
	assert.Equal(t, "user_age", dbOrders[1].Field, "Expected correct column for age order")
	assert.Equal(t, database.OrderDesc, dbOrders[1].Direction, "Expected descending order for age order")
}

// TestValidateAndTranslateToDBOrders tests the scenario where
// an invalid field is passed in
func TestValidateAndTranslateToDBOrders_InvalidField(t *testing.T) {
	orders := []Order{
		{Field: "invalid_field", Direction: DIRECTION_ASC},
	}

	allowedFields := []string{"name", "age", "email"}
	fieldTranslations := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ValidateAndTranslateToDBOrders(
		orders,
		allowedFields,
		fieldTranslations,
	)

	assert.Error(t, err, "Expected an error for an invalid field")
	assert.Nil(t, dbOrders, "Expected no database orders for an invalid field")
	apiErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Error should be of type *core.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "invalid_field", apiErr.Data().(Order).Field, "Error fields should match")
}

// TestValidateAndTranslateToDBOrders tests the scenario where
// an invalid direction is passed in
func TestValidateAndTranslateToDBOrders_InvalidDirection(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: "INVALID_DIRECTION"},
	}

	allowedFields := []string{"name", "age", "email"}
	fieldTranslations := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ValidateAndTranslateToDBOrders(
		orders,
		allowedFields,
		fieldTranslations,
	)

	assert.Error(t, err, "Expected an error for an invalid direction")
	assert.Nil(t, dbOrders, "Expected no database orders for an invalid direction")
	apiErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Error should be of type *core.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "name", apiErr.Data().(Order).Field, "Error fields should match")
}

// TestValidateAndTranslateToDBOrders tests the scenario where
// a field not in the translation map is passed in
func TestValidateAndTranslateToDBOrders_FieldNotInTranslationMap(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
	}

	allowedFields := []string{"name", "age", "email"}
	// Missing "name" field in the translation map to trigger the error
	fieldTranslations := map[string]DBField{
		"age": {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ValidateAndTranslateToDBOrders(
		orders,
		allowedFields,
		fieldTranslations,
	)

	assert.Error(t, err, "Expected an error for a field not present in the translation map")
	assert.Nil(t, dbOrders, "Expected no database orders for a field not present in the translation map")
	apiErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Error should be of type *core.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "name", apiErr.Data().(Order).Field, "Error fields should match")
}

// TestToDBOrders_ValidOrders tests the scenario where
// valid orders are passed in
func TestToDBOrders_ValidOrders(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
		{Field: "age", Direction: DIRECTION_DESC},
	}

	fieldTranslations := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ToDBOrders(orders, fieldTranslations)

	assert.Nil(t, err, "Expected no error for valid orders")
	assert.Len(t, dbOrders, 2, "Expected two translated orders")
	assert.Equal(t, "users", dbOrders[0].Table, "Expected correct table for name order")
	assert.Equal(t, "user_name", dbOrders[0].Field, "Expected correct column for name order")
	assert.Equal(t, database.OrderAsc, dbOrders[0].Direction, "Expected ascending order for name order")
	assert.Equal(t, "users", dbOrders[1].Table, "Expected correct table for age order")
	assert.Equal(t, "user_age", dbOrders[1].Field, "Expected correct column for age order")
	assert.Equal(t, database.OrderDesc, dbOrders[1].Direction, "Expected descending order for age order")
}

// TestToDBOrders_InvalidField tests the scenario where
// an invalid field is passed in
func TestToDBOrders_InvalidField(t *testing.T) {
	orders := []Order{
		{Field: "invalid_field", Direction: DIRECTION_ASC},
	}

	fieldTranslations := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbOrders, err := ToDBOrders(orders, fieldTranslations)

	assert.Error(t, err, "Expected an error for an invalid field")
	assert.Nil(t, dbOrders, "Expected no database orders for an invalid field")
	apiErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Error should be of type *core.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "invalid_field", apiErr.Data().(Order).Field, "Error fields should match")
}

// TestValidateAndDeduplicateOrders tests the ValidateAndDeduplicateOrders
// function with valid input.
func TestValidateAndDeduplicateOrders_ValidInput(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
		{Field: "age", Direction: DIRECTION_DESC},
	}

	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, orders, result, "Orders should match the input")
}

// TestValidateAndDeduplicateOrders_DuplicateFields tests the case where
// duplicate fields are provided.
func TestValidateAndDeduplicateOrders_DuplicateFields(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: DIRECTION_ASC},
		{Field: "name", Direction: DIRECTION_DESC},
		{Field: "age", Direction: DIRECTION_ASC},
	}

	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.Nil(t, err)
	assert.Len(t, result, 2, "Duplicate fields should be removed")
	assert.Equal(t, result[0], orders[0], "First order should be preserved")
	assert.Equal(t, result[1], orders[2], "Order should be preserved")
}

// TestValidateAndDeduplicateOrders_InvalidDirection tests the case where an
// invalid direction is provided.
func TestValidateAndDeduplicateOrders_InvalidDirection(t *testing.T) {
	orders := []Order{
		{Field: "name", Direction: "INVALID_DIRECTION"},
	}

	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.Error(t, err)
	assert.Nil(t, result, "Result should be nil")
	apiErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Error should be of type *core.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "name", apiErr.Data().(Order).Field, "Error fields should match")
}

// TestValidateAndDeduplicateOrders_InvalidField tests the case where an
// invalid field is provided.
func TestValidateAndDeduplicateOrders_InvalidField(t *testing.T) {
	orders := []Order{
		{Field: "invalid_field", Direction: DIRECTION_ASC},
	}

	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.Error(t, err)
	assert.Nil(t, result, "Result should be nil")
	apiErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Error should be of type *core.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "invalid_field", apiErr.Data().(Order).Field, "Error fields should match")
}

// TestValidateAndDeduplicateOrders_EmptyOrders tests the case where an
// empty list of orders is provided.
func TestValidateAndDeduplicateOrders_EmptyOrders(t *testing.T) {
	orders := []Order{}
	allowedFields := []string{"name", "age", "email"}

	result, err := ValidateAndDeduplicateOrders(orders, allowedFields)

	assert.Nil(t, err)
	assert.Empty(t, result, "Result should be empty")
}

// TestValidate_ValidOrder tests the case where a valid order is provided.
func TestValidate_ValidOrder(t *testing.T) {
	order := Order{
		Field:     "name",
		Direction: DIRECTION_ASC,
	}

	allowedFields := []string{"name", "age", "email"}

	err := validate(order, allowedFields)

	assert.Nil(t, err, "Expected no error for a valid order")
}

// TestValidate_InvalidDirection tests the case where an invalid direction
// is provided.
func TestValidate_InvalidDirection(t *testing.T) {
	order := Order{
		Field:     "name",
		Direction: "INVALID_DIRECTION",
	}

	allowedFields := []string{"name", "age", "email"}

	err := validate(order, allowedFields)

	assert.Error(t, err, "Expected an error for an invalid direction")
	apiErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Error should be of type *core.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "name", apiErr.Data().(Order).Field, "Error fields should match")
}

// TestValidate_InvalidField tests the case where an invalid field is provided.
func TestValidate_InvalidField(t *testing.T) {
	order := Order{
		Field:     "invalid_field",
		Direction: DIRECTION_ASC,
	}

	allowedFields := []string{"name", "age", "email"}

	err := validate(order, allowedFields)

	assert.Error(t, err, "Expected an error for an invalid field")
	apiErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Error should be of type *core.Error")
	assert.Equal(t, "INVALID_ORDER_FIELD", apiErr.ID, "Error ID should match")
	assert.Equal(t, "invalid_field", apiErr.Data().(Order).Field, "Errror fields should match")
}

// Validate_ValidLimit tests the Validate function for a valid limit.
func TestValidate_ValidLimit(t *testing.T) {
	page := &Page{
		Offset: 0,
		Limit:  5,
	}
	maxLimit := 10

	err := page.Validate(maxLimit)

	assert.Nil(t, err, "Expected no error when limit is within maxLimit")
}

// Validate_LimitExceeded tests the Validate function for a limit that exceeds
// the max limit.
func TestValidate_LimitExceeded(t *testing.T) {
	page := &Page{
		Offset: 0,
		Limit:  15,
	}
	maxLimit := 10

	err := page.Validate(maxLimit)

	assert.Error(t, err, "Expected an error when limit exceeds maxLimit")
	apiErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Error should be of type *core.Error")
	assert.Equal(t, "MAX_PAGE_LIMIT_EXCEEDED", apiErr.ID, "Error ID should")
	assert.Equal(t, maxLimit, apiErr.Data().(MaxPageLimitExceededErrorData).MaxLimit, "Max limit should match")
}

// Validate_ZeroLimit tests the Validate function for a limit of zero.
func TestValidate_ZeroLimit(t *testing.T) {
	page := &Page{
		Offset: 0,
		Limit:  0,
	}
	maxLimit := 10

	err := page.Validate(maxLimit)

	assert.Nil(t, err, "Expected no error when limit is zero")
}

// TestToDBUpdates_ValidInput tests the successful translation of updates to
// database updates.
func TestToDBUpdates_ValidInput(t *testing.T) {
	updates := []Update{
		{Field: "name", Value: "Alice"},
		{Field: "age", Value: 30},
	}

	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbUpdates, err := ToDBUpdates(updates, apiToDBFieldMap)

	assert.Nil(t, err, "Expected no error for valid updates")
	assert.Len(t, dbUpdates, 2, "Expected two translated updates")
	assert.Equal(t, "user_name", dbUpdates[0].Field, "Expected correct column for name update")
	assert.Equal(t, "Alice", dbUpdates[0].Value, "Expected correct value for name update")
	assert.Equal(t, "user_age", dbUpdates[1].Field, "Expected correct column for age update")
	assert.Equal(t, 30, dbUpdates[1].Value, "Expected correct value for age update")
}

// TestToDBUpdates_InvalidField tests the case when an update field cannot be
// translated.
func TestToDBUpdates_InvalidField(t *testing.T) {
	updates := []Update{
		{Field: "unknown_field", Value: "value"},
	}

	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbUpdates, err := ToDBUpdates(updates, apiToDBFieldMap)

	assert.Error(t, err, "Expected an error for an unknown field")
	assert.Nil(t, dbUpdates, "Expected no database updates for an unknown field")
	updateErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Expected error to be INVALID_DATABASE_UPDATE_TRANSLATION")
	assert.Equal(t, "unknown_field", updateErr.Data().(InvalidDatabaseUpdateTranslationErrorData).Field, "Expected error field to match the unknown field")
}

// TestToDBUpdates_EmptyUpdates tests the case when the input updates list is
// empty.
func TestToDBUpdates_EmptyUpdates(t *testing.T) {
	updates := []Update{}

	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbUpdates, err := ToDBUpdates(updates, apiToDBFieldMap)

	assert.Nil(t, err, "Expected no error for empty updates")
	assert.Empty(t, dbUpdates, "Expected no database updates for empty input")
}

// TestToDBUpdates_MultipleInvalidFields tests the case when multiple update
// fields cannot be translated.
func TestToDBUpdates_MultipleInvalidFields(t *testing.T) {
	updates := []Update{
		{Field: "unknown_field_1", Value: "value1"},
		{Field: "unknown_field_2", Value: "value2"},
	}

	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbUpdates, err := ToDBUpdates(updates, apiToDBFieldMap)

	assert.Error(t, err, "Expected an error for unknown fields")
	assert.Nil(t, dbUpdates, "Expected no database updates for unknown fields")
	updateErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Expected error to be of type InvalidDatabaseUpdateTranslationError")
	assert.Contains(t, []string{"unknown_field_1", "unknown_field_2"}, updateErr.Data().(InvalidDatabaseUpdateTranslationErrorData).Field, "Expected error field to match one of the unknown fields")
}

// TestSelectorString tests the String() method of the Selector struct.
func TestSelectorString(t *testing.T) {
	sel := Selector{
		Field:     "name",
		Predicate: EQUAL,
		Value:     "Alice",
	}

	expected := "name = Alice"
	assert.Equal(t, expected, sel.String(), "String() method should return the correct string representation")
}

// TestSelectors_GetByFields_SingleField tests the GetByFields method of
// Selectors when searching for a single field.
func TestSelectors_GetByFields_SingleField(t *testing.T) {
	selectors := Selectors{
		{Field: "name", Predicate: EQUAL, Value: "Alice"},
		{Field: "age", Predicate: GREATER, Value: 25},
		{Field: "email", Predicate: EQUAL, Value: "alic@example.com"},
	}

	result := selectors.GetByFields("name")

	assert.Len(t, result, 1, "Expected 1 selector for field 'name'")
	assert.Equal(t, "name", result[0].Field, "Field should match")
	assert.Equal(t, EQUAL, result[0].Predicate, "Predicate should match")
	assert.Equal(t, "Alice", result[0].Value, "Value should match")
}

// TestSelectors_GetByFields_MultipleFields tests the GetByFields method of
// Selectors when searching for multiple fields.
func TestSelectors_GetByFields_MultipleFields(t *testing.T) {
	selectors := Selectors{
		{Field: "name", Predicate: EQUAL, Value: "Alice"},
		{Field: "age", Predicate: GREATER, Value: 25},
		{Field: "email", Predicate: EQUAL, Value: "alic@example.com"},
	}

	result := selectors.GetByFields("name", "email")

	assert.Len(t, result, 2, "Expected 2 selectors for fields 'name' and 'email'")
	assert.Equal(t, "name", result[0].Field, "First result field should be 'name'")
	assert.Equal(t, "email", result[1].Field, "Second result field should be 'email'")
}

// TestSelectors_GetByFields_NoMatch tests the GetByFields method of Selectors
// when there are no matching fields.
func TestSelectors_GetByFields_NoMatch(t *testing.T) {
	selectors := Selectors{
		{Field: "name", Predicate: EQUAL, Value: "Alice"},
		{Field: "age", Predicate: GREATER, Value: 25},
	}

	result := selectors.GetByFields("email")

	assert.Len(t, result, 0, "Expected no selectors for field 'email'")
}

// TestToDBSelectors_ValidSelectors tests the successful translation of API
// selectors to DB selectors.
func TestToDBSelectors_ValidSelectors(t *testing.T) {
	apiSelectors := []Selector{
		{
			Field:             "name",
			Predicate:         EQUAL,
			Value:             "Alice",
			AllowedPredicates: []Predicate{EQUAL, NOT_EQUAL},
		},
		{
			Field:             "age",
			Predicate:         GREATER,
			Value:             25,
			AllowedPredicates: []Predicate{GREATER, LESS},
		},
	}

	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.Nil(t, err, "Expected no error for valid selectors")
	assert.Len(t, dbSelectors, 2, "Expected two translated DB selectors")

	assert.Equal(t, "users", dbSelectors[0].Table, "Expected correct table for 'name' selector")
	assert.Equal(t, "user_name", dbSelectors[0].Field, "Expected correct column for 'name' selector")
	assert.Equal(t, database.Predicate("="), dbSelectors[0].Predicate, "Expected correct predicate for 'name' selector")
	assert.Equal(t, "Alice", dbSelectors[0].Value, "Expected correct value for 'name' selector")

	assert.Equal(t, "users", dbSelectors[1].Table, "Expected correct table for 'age' selector")
	assert.Equal(t, "user_age", dbSelectors[1].Field, "Expected correct column for 'age' selector")
	assert.Equal(t, database.Predicate(">"), dbSelectors[1].Predicate, "Expected correct predicate for 'age' selector")
	assert.Equal(t, 25, dbSelectors[1].Value, "Expected correct value for 'age' selector")
}

// TestToDBSelectors_InvalidPredicate tests the case when a selector has an
// invalid
func TestToDBSelectors_InvalidPredicate(t *testing.T) {
	apiSelectors := []Selector{
		{
			Field:             "name",
			Predicate:         EQUAL,
			Value:             "Alice",
			AllowedPredicates: []Predicate{NOT_EQUAL},
		},
	}

	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.Error(t, err, "Expected error for invalid predicate")
	assert.Nil(t, dbSelectors, "Expected no database selectors when predicate is not allowed")
	predicateErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Expected error to be PREDICATE_NOT_ALLOWED")
	assert.Equal(t, EQUAL, predicateErr.Data().(PredicateNotAllowedErrorData).Predicate, "Expected error predicate to match the disallowed predicate")
}

// TestToDBSelectors_InvalidField tests the case when a selector field cannot be
// translated.
func TestToDBSelectors_InvalidField(t *testing.T) {
	apiSelectors := []Selector{
		{
			Field:             "unknown_field",
			Predicate:         EQUAL,
			Value:             "Alice",
			AllowedPredicates: []Predicate{EQUAL},
		},
	}

	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.Error(t, err, "Expected an error for an unknown field")
	assert.Nil(t, dbSelectors, "Expected no database selectors for an unknown field")
	fieldErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Expected error to be INVALID_DATABASE_SELECTOR_TRANSLATION")
	assert.Equal(t, "unknown_field", fieldErr.Data().(InvalidDatabaseSelectorTranslationErrorData).Field, "Expected error field to match the unknown field")
}

// TestToDBSelectors_InvalidDBPredicate tests the case when a predicate cannot
// be translated to a DB
func TestToDBSelectors_InvalidDBPredicate(t *testing.T) {
	apiSelectors := []Selector{
		{
			Field:             "name",
			Predicate:         "NONE",
			Value:             "Alice",
			AllowedPredicates: []Predicate{"NONE"},
		},
	}

	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.Error(t, err, "Expected an error for an invalid DB predicate")
	assert.Nil(t, dbSelectors, "Expected no database selectors for an invalid DB predicate")
	dbPredicateErr, ok := err.(*core.APIError)
	assert.True(t, ok, "Expected error to be INVALID_PREDICATE")
	assert.Equal(t, Predicate("NONE"), dbPredicateErr.Data().(InvalidPredicateErrorData).Predicate, "Expected error predicate to match the invalid predicate")
}

// TestToDBSelectors_EmptySelectors tests the case when the input selectors list
// is empty.
func TestToDBSelectors_EmptySelectors(t *testing.T) {
	apiSelectors := []Selector{}

	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "user_name"},
		"age":  {Table: "users", Column: "user_age"},
	}

	dbSelectors, err := ToDBSelectors(apiSelectors, apiToDBFieldMap)

	assert.Nil(t, err, "Expected no error for empty selectors")
	assert.Empty(t, dbSelectors, "Expected no database selectors for empty input")
}

// MockMiddleware is a simple middleware for testing.
func NewMockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Middleware", "Mocked")
		next.ServeHTTP(w, r)
	})
}

// TestMiddlewares_EmptyStack tests the Middlewares function when the stack is
// empty.
func TestMiddlewares_EmptyStack(t *testing.T) {
	mwStack := Stack{}

	middlewares := mwStack.Middlewares()

	assert.Empty(t, middlewares, "Expected middlewares to be empty for an empty stack")
}

// TestMiddlewares_StackWithMiddlewares tests the Middlewares function when the
// stack has middlewares.
func TestMiddlewares_StackWithMiddlewares(t *testing.T) {
	// Create some mock middleware wrappers.
	mw1 := Wrapper{
		ID:         "auth",
		Middleware: NewMockMiddleware,
	}
	mw2 := Wrapper{
		ID:         "logging",
		Middleware: NewMockMiddleware,
	}

	// Create a middleware stack with these wrappers.
	mwStack := Stack{mw1, mw2}

	middlewares := mwStack.Middlewares()

	assert.Equal(t, 2, len(middlewares), "Middleware stack should have 2 middlewares")
	assert.NotNil(t, middlewares[0], "First middleware should not be nil")
	assert.NotNil(t, middlewares[1], "Second middleware should not be nil")
}

// TestMiddlewares_Order tests the order of the middlewares returned by
// Middlewares function.
func TestMiddlewares_Order(t *testing.T) {
	callOrder := []string{}

	// Middleware 1: Add "first" to a shared slice
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "first")
			next.ServeHTTP(w, r)
		})
	}

	// Middleware 2: Add "second" to the shared slice
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "second")
			next.ServeHTTP(w, r)
		})
	}

	// Create the middleware stack
	mwStack := Stack{
		Wrapper{ID: "first", Middleware: mw1},
		Wrapper{ID: "second", Middleware: mw2},
	}

	// Get the middleware functions from the stack
	middlewares := mwStack.Middlewares()

	// Define a final handler
	finalHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {},
	)

	// Apply middlewares to the final handler
	wrappedHandler := core.ApplyMiddlewares(finalHandler, middlewares...)

	// Create a slice to track the middleware execution order
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req = req.WithContext(context.Background())

	// Perform the request
	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	// Verify the expected order of middleware execution
	assert.Equal(t, []string{"first", "second"}, callOrder, "Middlewares should be executed in the correct order")
}

// TestInsertAfterID_Success tests the InsertAfterID function when the
// middleware is inserted successfully.
func TestInsertAfterID_Success(t *testing.T) {
	mw1 := Wrapper{ID: "auth"}
	mw2 := Wrapper{ID: "logging"}
	mwStack := Stack{mw1, mw2}

	newMiddleware := Wrapper{ID: "metrics"}

	inserted := mwStack.InsertAfterID("auth", newMiddleware)

	assert.True(t, inserted, "Middleware should be inserted")
	assert.Equal(t, 3, len(mwStack), "Middleware stack should have 3 elements")
	assert.Equal(t, "auth", mwStack[0].ID)
	assert.Equal(t, "metrics", mwStack[1].ID, "Middleware not in 2nd position")
	assert.Equal(t, "logging", mwStack[2].ID)
}

// TestInsertAfterID_AppendToEnd tests the InsertAfterID function when the
// middleware is appended to the end.
func TestInsertAfterID_AppendToEnd(t *testing.T) {
	mw1 := Wrapper{ID: "auth"}
	mw2 := Wrapper{ID: "logging"}
	mwStack := Stack{mw1, mw2}

	newMiddleware := Wrapper{ID: "metrics"}

	inserted := mwStack.InsertAfterID("logging", newMiddleware)

	assert.True(t, inserted, "Middleware should be inserted")
	assert.Equal(t, 3, len(mwStack), "Middleware stack should have 3 elements")
	assert.Equal(t, "auth", mwStack[0].ID)
	assert.Equal(t, "logging", mwStack[1].ID)
	assert.Equal(t, "metrics", mwStack[2].ID, "New middleware not in the end")
}

// TestInsertAfterID_IDNotFound tests the InsertAfterID function when the
// middleware ID is not found.
func TestInsertAfterID_IDNotFound(t *testing.T) {
	mw1 := Wrapper{ID: "auth"}
	mw2 := Wrapper{ID: "logging"}
	mwStack := Stack{mw1, mw2}

	newMiddleware := Wrapper{ID: "metrics"}

	inserted := mwStack.InsertAfterID("non-existent-id", newMiddleware)

	assert.False(t, inserted, "Middleware should not be inserted")
	assert.Equal(t, 2, len(mwStack), "Middleware stack should have 2 elements")
	assert.Equal(t, "auth", mwStack[0].ID)
	assert.Equal(t, "logging", mwStack[1].ID)
}
