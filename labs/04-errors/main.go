package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type ValidationError struct {
	Field string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Reason)
}

func loadUser(id int) error {
	if id < 0 {
		return fmt.Errorf("load user %d: %w", id, ErrNotFound)
	}
	if id == 0 {
		return &ValidationError{Field: "id", Reason: "must be positive"}
	}
	return nil
}

func main() {
	err := loadUser(-1)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("errors.Is matched ErrNotFound:", err)
	}

	err = loadUser(0)
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Println("errors.As got field:", ve.Field)
	}

	// Unwrap chain
	wrapped := fmt.Errorf("outer: %w", ErrNotFound)
	fmt.Println("wrapped chain:", errors.Is(wrapped, ErrNotFound))
}