package internal

import (
	"math"
)

const (
	// MaxOperationsPerBatchedWrite is the max. number of op that a batched write can hold and process.
	MaxOperationsPerBatchedWrite = 500
)

// CalculateRequiredBatches returns the number of batched writes necessary
// given a amount of operations.
func CalculateRequiredBatches(operationCount int) int {
	return int(math.Ceil(float64(operationCount) / float64(MaxOperationsPerBatchedWrite)))
}
