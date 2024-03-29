package go_playground_pool

const (
	errCancelled = "ERROR: Work Unit Cancelled"
	errRecovery  = "ERROR: Work Unit failed due to algorithm recoverable error: '%v'\n, Stack Trace:\n %s"
	errClosed    = "ERROR: Work Unit added/run after the pool had been closed or cancelled"
)

// ErrRecovery contains the error when algorithm consumer goroutine needed to be recovers
type ErrRecovery struct {
	s string
}

// Error prints recovery error
func (e *ErrRecovery) Error() string {
	return e.s
}

// ErrPoolClosed is the error returned to all work units that may have been in or added to the pool after it's closing.
type ErrPoolClosed struct {
	s string
}

// Error prints Work Unit Close error
func (e *ErrPoolClosed) Error() string {
	return e.s
}

// ErrCancelled is the error returned to algorithm Work Unit when it has been cancelled.
type ErrCancelled struct {
	s string
}

// Error prints Work Unit Cancellation error
func (e *ErrCancelled) Error() string {
	return e.s
}
