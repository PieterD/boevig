package deebee

import (
	"errors"
	"fmt"
)

type unwrappableErrorInterface interface {
	Error() string
	Unwrap() error
}

type retryAllowedError struct {
	cause error
}

func (e retryAllowedError) Error() string {
	return fmt.Sprintf("retry allowed: %v", e.cause)
}

func (e retryAllowedError) Unwrap() error {
	return e.cause
}

var _ unwrappableErrorInterface = retryAllowedError{}

func isRetryAllowedError(err error) bool {
	return errors.As(err, &retryAllowedError{})
}

func newRetryAllowedError(cause error) error {
	return retryAllowedError{cause: cause}
}
