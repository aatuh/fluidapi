package test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/pakkasys/fluidapi/database"
// 	"github.com/pakkasys/fluidapi/database/mock"
// 	"github.com/pakkasys/fluidapi/util"
// 	"github.com/stretchr/testify/assert"
// )

// // TestExecuteTransaction_Success tests the case where a transaction is
// // successfully executed.
// func TestExecuteTransaction_Success(t *testing.T) {
// 	mockTx := new(mock.MockTx)

// 	// Mock the transactional function to return a successful result
// 	txFn := func(ctx context.Context, tx database.Tx) (string, error) {
// 		return "success", nil
// 	}

// 	// Setup the mock transaction expectations
// 	mockTx.On("Commit").Return(nil).Once()

// 	ctx := util.NewContext(context.Background())
// 	result, err := database.Transaction(ctx, mockTx, txFn)

// 	assert.Nil(t, err)
// 	assert.Equal(t, "success", result)
// 	mockTx.AssertExpectations(t)
// }

// // TestExecuteTransaction_TransactionError tests the case where the
// // transactional function returns an error.
// func TestExecuteTransaction_TransactionalFnError(t *testing.T) {
// 	mockTx := new(mock.MockTx)

// 	// Mock the transactional function to return an error
// 	txFn := func(ctx context.Context, tx database.Tx) (string, error) {
// 		return "", errors.New("application error")
// 	}

// 	// Setup the mock transaction expectations
// 	mockTx.On("Rollback").Return(nil).Once()

// 	ctx := util.NewContext(context.Background())
// 	result, err := database.Transaction(ctx, mockTx, txFn)

// 	assert.Equal(t, "", result)
// 	assert.EqualError(t, err, "application error")
// 	mockTx.AssertExpectations(t)
// }

// // TestExecuteTransaction_TransactionalFnError tests the case where the
// // transactional function returns an error.
// func TestExecuteTransaction_FinalizeError(t *testing.T) {
// 	mockTx := new(mock.MockTx)

// 	// Mock the transactional function to return an error
// 	txFn := func(ctx context.Context, tx database.Tx) (string, error) {
// 		return "", nil
// 	}

// 	// Setup the mock transaction expectations
// 	mockTx.On("Commit").Return(errors.New("commit error")).Once()

// 	ctx := util.NewContext(context.Background())
// 	result, err := database.Transaction(ctx, mockTx, txFn)

// 	assert.Equal(t, "", result)
// 	assert.EqualError(t, err, "failed to commit transaction: commit error")
// 	mockTx.AssertExpectations(t)
// }

// // TestFinalizeTransaction_SuccessfulCommit tests the successful commit case.
// func TestFinalizeTransaction_SuccessfulCommit(t *testing.T) {
// 	mockTx := new(mock.MockTx)

// 	// Setup the mock transaction expectations
// 	mockTx.On("Commit").Return(nil).Once()

// 	err := finalizeTransaction(mockTx, nil)

// 	assert.Nil(t, err)
// 	mockTx.AssertExpectations(t)
// }

// // TestFinalizeTransaction_CommitError tests the case where commit fails.
// func TestFinalizeTransaction_CommitError(t *testing.T) {
// 	mockTx := new(mock.MockTx)

// 	// Setup the mock transaction expectations
// 	mockTx.On("Commit").Return(errors.New("commit error")).Once()

// 	err := finalizeTransaction(mockTx, nil)

// 	assert.EqualError(t, err, "failed to commit transaction: commit error")
// 	mockTx.AssertExpectations(t)
// }

// // TestFinalizeTransaction_RollbackError tests the case where rollback fails.
// func TestFinalizeTransaction_RollbackError(t *testing.T) {
// 	mockTx := new(mock.MockTx)

// 	// Setup the mock transaction expectations
// 	mockTx.On("Rollback").Return(errors.New("rollback error")).Once()

// 	err := finalizeTransaction(mockTx, errors.New("transaction error"))

// 	assert.EqualError(t, err, "failed to rollback transaction: rollback error")
// 	mockTx.AssertExpectations(t)
// }
