package runner

import (
	"net/http"

	apierror "github.com/pakkasys/fluidapi/core/api/error"
	"github.com/pakkasys/fluidapi/endpoint/middleware/inputlogic"
	"github.com/pakkasys/fluidapi/endpoint/order"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/pakkasys/fluidapi/endpoint/selector"
)

var (
	DuplicateEntryError    = apierror.New[error]("DUPLICATE_ENTRY")
	ForeignConstraintError = apierror.New[error]("FOREIGN_CONSTRAINT_ERROR")
)

var CreateErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ID:         DuplicateEntryError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: false,
	},
	{
		ID:         ForeignConstraintError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: false,
	},
}

var GetErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ID:         selector.InvalidPredicateError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         selector.PredicateNotAllowedError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         selector.InvalidSelectorFieldError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         order.InvalidOrderFieldError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         page.MaxPageLimitExceededError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
}

var UpdateErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ID:         selector.InvalidPredicateError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         selector.InvalidSelectorFieldError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         selector.PredicateNotAllowedError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         NeedAtLeastOneSelectorError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         NeedAtLeastOneUpdateError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         order.InvalidOrderFieldError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         DuplicateEntryError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: false,
	},
	{
		ID:         ForeignConstraintError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: false,
	},
}

var DeleteErrors []inputlogic.ExpectedError = []inputlogic.ExpectedError{
	{
		ID:         selector.InvalidPredicateError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         selector.InvalidSelectorFieldError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         selector.PredicateNotAllowedError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
	{
		ID:         NeedAtLeastOneSelectorError.ID(),
		Status:     http.StatusBadRequest,
		PublicData: true,
	},
}
