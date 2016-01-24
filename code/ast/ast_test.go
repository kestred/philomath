package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// NOTE: This file clearly displays my general distrust of Go's type system

func IsNode(v interface{}) bool {
	_, ok := v.(Node)
	return ok
}

func IsDecl(v interface{}) bool {
	_, ok := v.(Decl)
	return ok
}

func TestConstructors(t *testing.T) {
	decl := Mutable("foo", nil, Ident("bar"))
	assert.Equal(t, UninferredType, decl.Type)
	assert.True(t, IsNode(decl))
	assert.True(t, IsDecl(decl))
}
