package failure

import (
	"fmt"
)

// Builder helps build an error
type Builder interface {
	WithField(name string, value interface{}) Builder
	WithFields(fields Fields) Builder
	WithMessage(message string) Builder
	WithMessagef(format string, a ...interface{}) Builder
	ParentOf(inner error) Builder
	Done() error
}

// Impersonator gives the ability to work with external error types
// that implement this interface.
type Impersonator interface {
	error
	Impersonate(b Builder)
}

type builder struct {
	n *node
}

// Done returns the new built error
func (b *builder) Done() error {
	return b.n
}

// WithField adds a field named name with the value value to
// the error
func (b *builder) WithField(name string, value interface{}) Builder {

	if name == "" {
		return b
	}

	if name == MessageField {
		b.n.Message = fmt.Sprintf("%v", value)
		return b
	}

	if b.n.Fields == nil {
		b.n.Fields = make(Fields)
	}

	b.n.Fields[name] = value

	return b
}

// WithFields adds fields to the error
func (b *builder) WithFields(fields Fields) Builder {

	if fields == nil || len(fields) == 0 {
		return b
	}

	for k, v := range fields {

		if k == "" {
			continue
		}

		if k == MessageField {
			b.n.Message = fmt.Sprintf("%v", v)
		} else {

			if b.n.Fields == nil {
				b.n.Fields = make(Fields)
			}

			b.n.Fields[k] = v
		}
	}

	return b
}

// ParentOf sets the new error to be the parent of inner
func (b *builder) ParentOf(inner error) Builder {

	var n *node

	if inner != nil {
		n = impersonate(inner)
	}

	b.n.Inner = n

	return b
}

// WithMessage sets the error's message
func (b *builder) WithMessage(message string) Builder {
	b.n.Message = message
	return b
}

// WithMessagef formats the error's message according to
// a format specifier
func (b *builder) WithMessagef(format string, a ...interface{}) Builder {
	return b.WithMessage(fmt.Sprintf(format, a...))
}
