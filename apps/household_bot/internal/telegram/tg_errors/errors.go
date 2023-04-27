package tg_errors

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"
)

type Error struct {
	originalErr error
	context     struct {
		// Handler name
		handler string
		// Caused by method that returned an error
		// e.g.
		// err := h.repo.GetAll()
		// if err != nil {
		//		>> repo.GetAll is causedBy
		// }
		//
		//
		causedBy string
	}
}

type Config struct {
	OriginalErr error
	Handler     string
	CausedBy    string
}

func New(c Config) *Error {
	var tgError *Error
	if errors.As(c.OriginalErr, &tgError) {
		return &Error{
			originalErr: multierr.Append(tgError.originalErr, c.OriginalErr),
			context: struct {
				handler  string
				causedBy string
			}{
				handler:  c.Handler,
				causedBy: prependCausedBy(tgError.context.causedBy, c.CausedBy),
			},
		}
	}
	return &Error{
		originalErr: c.OriginalErr,
		context: struct {
			handler  string
			causedBy string
		}{
			handler:  c.Handler,
			causedBy: c.CausedBy,
		},
	}
}

func (e *Error) Error() string {
	return e.originalErr.Error()
}

func (e *Error) Is(target error) bool {
	return e.originalErr == target
}

func (e *Error) String() string {
	return fmt.Sprintf("handler [%s] causedBy [%s] msg [%s]", e.context.handler, e.context.causedBy, e.originalErr.Error())
}

func prependCausedBy(left, right string) string {
	return left + ":" + right
}
