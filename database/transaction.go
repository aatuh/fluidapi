package database

import (
	"context"
	"fmt"
)

// TransactionalFunc is a function that takes in a transaction, and returns a
// result and an error.
type TransactionalFunc[Result any] func(
	ctx context.Context, tx Tx,
) (Result, error)

// Transaction executes a TransactionalFunc within a transaction.
// It recovers from panics, rolls back on errors, and commits if no error
// occurs.
//
// Parameters:
//   - ctx: The context for the transaction.
//   - tx: The transaction to use.
//   - transactionalFn: The function to execute in a transaction.
//
// Returns:
//   - Result: The result of the transactional function.
//   - error: An error if the transaction fails.
func Transaction[Result any](
	ctx context.Context,
	tx Tx,
	transactionalFn TransactionalFunc[Result],
) (result Result, txErr error) {
	defer func() {
		var recovered any
		panicOccurred := false
		// Recover from panics.
		if recovered = recover(); recovered != nil {
			panicOccurred = true
			txErr = fmt.Errorf("panic in transactional function: %v", recovered)
		}
		if err := finalizeTransaction(tx, txErr); err != nil {
			txErr = err
			var zero Result
			result = zero
		}
		// Propagate the panic if there was one.
		if panicOccurred {
			panic(recovered)
		}
	}()
	return transactionalFn(ctx, tx)
}

// finalizeTransaction commits or rollbacks a transaction.
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
