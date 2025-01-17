package errors

import (
	apierror "github.com/pakkasys/fluidapi/core/api/error"
)

// TODO: Move elsewhere (nearer where it's used)
var (
	DuplicateEntryError    = apierror.New[error]("DUPLICATE_ENTRY")
	ForeignConstraintError = apierror.New[error]("FOREIGN_CONSTRAINT_ERROR")
)
