package tg_errors

import (
	"encoding/json"
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

// ToJSON returns descriptive representation of an error
func (e *Error) ToJSON() (string, error) {
	b, err := json.Marshal(struct {
		Handler, CausedBy, OriginalError string
	}{
		Handler:       e.context.handler,
		CausedBy:      e.context.causedBy,
		OriginalError: e.originalErr.Error(),
	})
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b), nil
}

func prependCausedBy(left, right string) string {
	return left + ":" + right
}
