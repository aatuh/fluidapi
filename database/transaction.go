package database

import (
	"context"
	"fmt"
)

// TransactionalFunc is a function that takes a transaction and returns a
// result.
type TransactionalFunc[Result any] func(
	ctx context.Context,
	tx Tx,
) (Result, error)

// Transaction executes a TransactionalFunc in a transaction.
//
//   - tx: The transaction to use.
//   - transactionalFn: The function to execute in a transaction.
func Transaction[Result any](
	ctx context.Context,
	tx Tx,
	transactionalFn TransactionalFunc[Result],
) (result Result, txErr error) {
	defer func() {
		var r any
		hasPanic := false
		if r = recover(); r != nil {
			hasPanic = true
			txErr = fmt.Errorf("panic in transactional function: %v", r)
		}
		if err := finalizeTransaction(tx, txErr); err != nil {
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

func finalizeTransaction(tx Tx, txErr error) error {
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
