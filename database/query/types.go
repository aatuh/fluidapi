package query

// ErrorChecker is an interface for checking database errors for a specific
// driver.
type ErrorChecker interface {
	Check(err error) error
}
