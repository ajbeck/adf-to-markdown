//go:build goexperiment.jsonv2

package adfmarkdown

import "fmt"

type Error struct {
	Path   string
	Kind   ErrorKind
	Detail string
	Cause  error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s at %s: %s: %v", e.Kind, e.Path, e.Detail, e.Cause)
	}
	return fmt.Sprintf("%s at %s: %s", e.Kind, e.Path, e.Detail)
}

func (e *Error) Unwrap() error { return e.Cause }

func newDecodeError(path string, kind ErrorKind, detail string, cause ...error) error {
	var err error
	if len(cause) > 0 {
		err = cause[0]
	}
	return &Error{
		Path:   path,
		Kind:   kind,
		Detail: detail,
		Cause:  err,
	}
}
