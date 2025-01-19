package page

import (
	"testing"

	apierror "github.com/pakkasys/fluidapi/core/api/error"
	"github.com/stretchr/testify/assert"
)

// Validate_ValidLimit tests the Validate function for a valid limit.
func TestValidate_ValidLimit(t *testing.T) {
	page := &Page{
		Offset: 0,
		Limit:  5,
	}
	maxLimit := 10

	err := page.Validate(maxLimit)

	assert.NoError(t, err, "Expected no error when limit is within maxLimit")
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
	apiErr, ok := err.(*apierror.Error[MaxPageLimitExceededErrorData])
	assert.True(t, ok, "Error should be of type *api.Error")
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

	assert.NoError(t, err, "Expected no error when limit is zero")
}

// TestGetLimitOffsetClauseFromPage_NoPage tests the case where no page is
// provided.
func TestGetLimitOffsetClauseFromPage_NoPage(t *testing.T) {
	var p *Page = nil
	limitOffsetClause := GetLimitOffsetClauseFromPage(p)
	assert.Equal(t, "", limitOffsetClause)
}

// TestGetLimitOffsetClauseFromPage_WithPage tests the case where a page with
// limit and offset is provided.
func TestGetLimitOffsetClauseFromPage_WithPage(t *testing.T) {
	p := &Page{Limit: 10, Offset: 20}

	limitOffsetClause := GetLimitOffsetClauseFromPage(p)

	expected := "LIMIT 10 OFFSET 20"
	assert.Equal(t, expected, limitOffsetClause)
}

// TestGetLimitOffsetClauseFromPage_ZeroLimit tests the case where limit is 0.
func TestGetLimitOffsetClauseFromPage_ZeroLimit(t *testing.T) {
	p := &Page{Limit: 0, Offset: 20}

	limitOffsetClause := GetLimitOffsetClauseFromPage(p)

	expected := "LIMIT 0 OFFSET 20"
	assert.Equal(t, expected, limitOffsetClause)
}

// TestGetLimitOffsetClauseFromPage_ZeroOffset tests the case where offset is 0.
func TestGetLimitOffsetClauseFromPage_ZeroOffset(t *testing.T) {
	p := &Page{Limit: 10, Offset: 0}

	limitOffsetClause := GetLimitOffsetClauseFromPage(p)

	expected := "LIMIT 10 OFFSET 0"
	assert.Equal(t, expected, limitOffsetClause)
}

// TestGetLimitOffsetClauseFromPage_ZeroLimitAndOffset tests the case where both
// limit and offset are 0.
func TestGetLimitOffsetClauseFromPage_ZeroLimitAndOffset(t *testing.T) {
	p := &Page{Limit: 0, Offset: 0}

	limitOffsetClause := GetLimitOffsetClauseFromPage(p)

	expected := "LIMIT 0 OFFSET 0"
	assert.Equal(t, expected, limitOffsetClause)
}
