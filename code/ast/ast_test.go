package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// NOTE: This file clearly displays my general distrust of Go's type system

func isNode(v interface{}) bool {
	_, ok := v.(Node)
	return ok
}

func isDecl(v interface{}) bool {
	_, ok := v.(Decl)
	return ok
}

func isExpr(v interface{}) bool {
	_, ok := v.(Expr)
	return ok
}

func TestConstructors(t *testing.T) {
	decl := Mutable("foo", nil, Ident("bar"))
	assert.Equal(t, InferredType, decl.Type)
	assert.Equal(t, UnresolvedType, decl.Expr.GetType())
	assert.True(t, isNode(decl))
	assert.True(t, isDecl(decl))

	expr := ProcExp(nil, nil, Blok(nil))
	assert.Equal(t, InferredType, expr.Return)
	assert.Equal(t, UninferredType, expr.Type)
	assert.True(t, isNode(expr))
	assert.True(t, isExpr(expr))
}
