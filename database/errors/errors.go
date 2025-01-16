package errors

import (
	apierror "github.com/pakkasys/fluidapi/core/api/error"
)

var (
	DuplicateEntryError    = apierror.New[error]("DUPLICATE_ENTRY")
	ForeignConstraintError = apierror.New[error]("FOREIGN_CONSTRAINT_ERROR")
)
