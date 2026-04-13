package cliutil

import (
	"errors"
	"fmt"
	"os"

	"github.com/walnut1024/efi-cli/internal/output"
)

type ErrorKind string

const (
	ErrCodeNotFound  ErrorKind = "code_not_found"
	ErrAmbiguousCode ErrorKind = "ambiguous_code"
	ErrUpstream      ErrorKind = "upstream_error"
	ErrDecode        ErrorKind = "decode_error"
	ErrInvalidArg    ErrorKind = "invalid_argument"
)

type CLIError struct {
	Kind    ErrorKind      `json:"kind"`
	Message string         `json:"message"`
	Input   string         `json:"input,omitempty"`
	Op      string         `json:"op,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
	Cause   error          `json:"-"`
}

func (e *CLIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause == nil {
		if e.Input == "" {
			return e.Message
		}
		return fmt.Sprintf("%s: %s", e.Message, e.Input)
	}
	if e.Input == "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s: %v", e.Message, e.Input, e.Cause)
}

func (e *CLIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func NewError(kind ErrorKind, message, input, op string, cause error, meta map[string]any) *CLIError {
	return &CLIError{
		Kind:    kind,
		Message: message,
		Input:   input,
		Op:      op,
		Meta:    meta,
		Cause:   cause,
	}
}

func AsCLIError(err error) (*CLIError, bool) {
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr, true
	}
	return nil, false
}

func PrintErrorAndExit(err error, cfg OutputConfig) {
	if cfg.useJSONOutput() {
		payload := map[string]any{
			"error": map[string]any{
				"kind":    "unknown_error",
				"message": err.Error(),
			},
		}
		if cliErr, ok := AsCLIError(err); ok {
			payload["error"] = map[string]any{
				"kind":    cliErr.Kind,
				"message": cliErr.Message,
				"input":   cliErr.Input,
				"op":      cliErr.Op,
				"meta":    cliErr.Meta,
			}
		}
		if writeErr := output.WriteJSON(os.Stderr, payload, cfg.Pretty || !cfg.Compact); writeErr != nil {
			fmt.Fprintln(os.Stderr, writeErr)
		}
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
