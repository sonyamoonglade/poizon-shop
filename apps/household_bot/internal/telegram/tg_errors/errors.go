package tg_errors

import (
	"encoding/json"
	"fmt"
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
		handler, causedBy, originalError string
	}{
		handler:       e.context.handler,
		causedBy:      e.context.causedBy,
		originalError: e.originalErr.Error(),
	})
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b), nil
}
