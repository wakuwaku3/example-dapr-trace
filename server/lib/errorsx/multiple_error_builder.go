package errorsx

import "errors"

type MultipleErrorBuilder struct {
	Errors []error `json:"errors,omitempty"`
}

func NewMultipleErrorBuilder() *MultipleErrorBuilder {
	return &MultipleErrorBuilder{
		Errors: []error{},
	}
}
func (e *MultipleErrorBuilder) Append(err error) *MultipleErrorBuilder {
	e.Errors = append(e.Errors, err)
	return e
}
func (e *MultipleErrorBuilder) Build() error {
	if len(e.Errors) == 0 {
		return nil
	}
	if len(e.Errors) == 1 {
		return e.Errors[0]
	}

	return Wrap(errors.Join(e.Errors...), e.Errors)
}
