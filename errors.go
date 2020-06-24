package fuego

import "errors"

var (
	// ErrBatchWriteNotStarted indicates that the a commit can't happen if the batch write hasn't been started.
	ErrBatchWriteNotStarted = errors.New("fuego: no batch write started")
)
