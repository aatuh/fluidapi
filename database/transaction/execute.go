package transaction

import (
	"context"
	"fmt"

	"github.com/pakkasys/fluidapi/database"
)

// TransactionalFunc is a function that takes a transaction and returns a
// result.
type TransactionalFunc[Result any] func(
	ctx context.Context,
	tx database.Tx,
) (Result, error)

// Execute executes a TransactionalFunc in a transaction.
//
//   - tx: The transaction to use.
//   - transactionalFn: The function to execute in a transaction.
func Execute[Result any](
	ctx context.Context,
	tx database.Tx,
	transactionalFn TransactionalFunc[Result],
) (result Result, txErr error) {
	defer func() {
		var r any
		hasPanic := false
		if r = recover(); r != nil {
			hasPanic = true
			txErr = fmt.Errorf("panic in transactional function: %v", r)
		}
		if err := finalize(tx, txErr); err != nil {
			txErr = err
			var zero Result
			result = zero
		}

		// Propagate the panic if there was one
		if hasPanic {
			panic(r)
		}
	}()
	return transactionalFn(ctx, tx)
}

func finalize(tx database.Tx, txErr error) error {
	if txErr != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %v", err)
		}
		return nil
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
