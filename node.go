package failure

import (
	"encoding/json"
)

const (
	// MessageField is the name for error's message field
	MessageField = "message"
)

var (
	// ErrNoSuchField returned if the field does not exist
	ErrNoSuchField = New("failure: no such field")

	// ErrInvalidFieldName returned if the field name is empty
	ErrInvalidFieldName = New("failure: invalid field name")
)

// Fields is a map of error key value data pairs
type Fields map[string]interface{}

type node struct {
	Message string `json:"message,omitempty"`
	Fields  Fields `json:"fields,omitempty"`
	Inner   *node  `json:"inner,omitempty"`
}

func (n *node) Error() string {

	j, err := json.Marshal(*n)

	if err != nil {
		panic(err)
	}

	return string(j)
}

func (n *node) isParentOf(inner *node) bool {

	if n == nil {
		return false
	}

	if inner == nil {
		return true
	}

	var c bool

	n = n.Inner

	for c = n == inner; n != nil && !c; c = n == inner {
		n = n.Inner
	}

	return c
}

func (n *node) origin() *node {

	o := n

	for o != nil {
		n = o
		o = o.Inner
	}

	return n
}

func (n *node) field(name string) (interface{}, error) {

	if name == "" {
		return nil, ErrInvalidFieldName
	}

	if n == nil {
		return nil, newNoSuchFieldError(name)
	}

	if name == MessageField {
		return n.Message, nil
	}

	if n.Fields == nil {
		return nil, newNoSuchFieldError(name)
	}

	v, found := n.Fields[name]

	if !found {
		return nil, newNoSuchFieldError(name)
	}

	return v, nil
}

func (n *node) same(other *node) bool {

	for n != nil && other != nil {

		if n == other {
			return true
		}

		if n.Message != other.Message {
			return false
		}

		if len(n.Fields) != len(other.Fields) {
			return false
		}

		for k, v0 := range n.Fields {

			if v1, found := other.Fields[k]; found {

				if v0 != v1 {
					return false
				}

			} else {
				return false
			}
		}

		n = n.Inner
		other = other.Inner
	}

	return (n != nil || other == nil) && (n == nil || other != nil)
}

func (n *node) like(other *node) bool {

	if n == other {
		return true
	}

	if n == nil || other == nil {
		return false
	}

	return n.Message == other.Message
}

func (n *node) copy() *node {

	if n == nil {
		return nil
	}

	var p *node

	r := &node{Message: n.Message}
	c := r

	for {

		if n.Fields != nil {

			c.Fields = make(Fields)

			for k, v := range n.Fields {
				c.Fields[k] = v
			}
		}

		if p != nil {
			p.Inner = c
		}

		p = c
		n = n.Inner

		if n != nil {
			c = &node{Message: n.Message}
		} else {
			break
		}
	}

	return r
}

func (n *node) depth() int {

	var d int

	for n != nil {
		n = n.Inner
		d++
	}

	return d
}

func impersonate(err error) *node {

	if err == nil {
		return nil
	}

	n, ok := err.(*node)

	if ok {
		return n
	}

	imp, ok := err.(Impersonator)

	if !ok {
		return &node{Message: err.Error()}
	}

	b := &builder{n: &node{Message: err.Error()}}
	imp.Impersonate(b)

	return b.n
}

func newNoSuchFieldError(name string) error {
	return Buildc(ErrNoSuchField).WithField("name", name).Done()
}
