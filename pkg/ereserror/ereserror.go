package ereserror

import "github.com/pkg/errors"

type error interface {
	Error() string
}

type EresError struct {
	Message string `json:"message"` // Human readable message for clients
	Code    int    `json:"-"`       // HTTP Status code. We use `-` to skip json marshaling.
	Err     error  `json:"-"`       // The original error. Same reason as above.
}

func New(code int, err error, message string) error {
	return EresError{
		Message: message,
		Code:    code,
		Err:     errors.Wrap(err, message),
	}
}

// Returns Message if Err is nil. You can handle custom implementation of your own.
func (err EresError) Error() string {
	// guard against panics
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (err EresError) Unwrap() error {
	return err.Err // Returns inner error
}

// Returns the inner most CustomErrorWrapper
func (err EresError) Dig() EresError {
	var ew EresError
	if errors.As(err.Err, &ew) {
		// Recursively digs until wrapper error is not in which case it will stop
		return ew.Dig()
	}
	return err
}
