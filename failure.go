package failure

import (
	"fmt"
)

// New returns an error that formats as the given message.
func New(message string) error {
	return Build(message).Done()
}

// Newf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Newf(format string, a ...interface{}) error {
	return Buildf(format, a...).Done()
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error (same as calling Newf).
func Errorf(format string, a ...interface{}) error {
	return Newf(format, a...)
}

// Build returns an error builder for building an error that
// formats as the given message.
func Build(message string) Builder {
	return &builder{n: &node{Message: message}}
}

// Buildf returns an error builder for building an error that
// formats according to a format specifier.
func Buildf(format string, a ...interface{}) Builder {
	return Build(fmt.Sprintf(format, a...))
}

// Buildc returns an error builder for building an error which is
// based on copy of err.
func Buildc(err error) Builder {

	if err == nil {
		panic("err is nil")
	}

	return &builder{n: impersonate(err).copy()}
}

// Field returns the value (or error if not found) for a field of
// err given its name.
func Field(err error, name string) (interface{}, error) {

	if err == nil {
		return nil, newNoSuchFieldError(name)
	}

	return impersonate(err).field(name)
}

// FieldOrDefault returns the value for a field of err given its
// name or def if it doesn't exist.
func FieldOrDefault(err error, name string, def interface{}) interface{} {

	v, e := Field(err, name)

	if e != nil {
		return def
	}

	return v
}

// TestField returns true only if a field named name with the value
// value exists in err.
func TestField(err error, name string, value interface{}) bool {

	v, e := Field(err, name)

	if e != nil {
		return false
	}

	return v == value
}

// TestFieldRecursively returns true only if a field named name with the value
// value exists in err or one of its descendants.
func TestFieldRecursively(err error, name string, value interface{}) bool {

	n := impersonate(err)

	for n != nil {
		if v, e := n.field(name); e == nil && v == value {
			return true
		}

		n = n.Inner
	}

	return false
}

// Inner returns the inner error of err.
func Inner(err error) error {

	n := impersonate(err)

	if n.Inner == nil {
		return nil
	}

	return n.Inner
}

// Origin returns the deepest inner error of err.
func Origin(err error) error {
	return impersonate(err).origin()
}

// IsParentOf returns true only if inner is a descendant of err.
func IsParentOf(err, inner error) bool {
	return impersonate(err).isParentOf(impersonate(inner))
}

// Same returns true only if err0 and err1 have the same message, same fields
// and their values and same descendants.
func Same(err0, err1 error) bool {

	if err0 == nil && err1 == nil {
		return true
	}

	if err0 == nil || err1 == nil {
		return false
	}

	return impersonate(err0).same(impersonate(err1))
}

// Like returns true only if err0 and err1 have the same message
// or both nil.
func Like(err0, err1 error) bool {

	if err0 == nil && err1 == nil {
		return true
	}

	if err0 == nil || err1 == nil {
		return false
	}

	return impersonate(err0).like(impersonate(err1))
}

// Message returns err's message field.
func Message(err error) string {

	if err == nil {
		panic("err is nil")
	}

	return impersonate(err).Message
}

// Depth returns the count for err and its descendants
func Depth(err error) int {
	return impersonate(err).depth()
}
