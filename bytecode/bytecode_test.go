package bytecode

import (
	"testing"

	// TODO: Stop relying on parser when more code is stable
	"github.com/kestred/philomath/parser"
	"github.com/stretchr/testify/assert"
)

func encodeExpression(input string) (*Scope, []Instruction) {
	var p parser.Parser
	p.Init("test", false, []byte(input))
	expr := p.ParseExpression()

	scope := &Scope{}
	scope.Init()
	return scope, FromExpr(expr, scope)
}

func TestEncodeArithmetic(t *testing.T) {
	// constant
	expected := []Instruction{{Code: STORE_CONST, Out: 1}}
	scope, insts := encodeExpression(`22`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, Data(22), scope.Constants[1])

	// add, subtract, multiply, divide
	expected = []Instruction{
		{Code: INT64_MULTIPLY, Left: 1, Right: 2, Out: 3},
		{Code: INT64_DIVIDE, Left: 4, Right: 5, Out: 6},
		{Code: INT64_ADD, Left: 3, Right: 6, Out: 7},
		{Code: INT64_SUBTRACT, Left: 7, Right: 8, Out: 9},
	}
	scope, insts = encodeExpression(`2 * 3 + 4 / 5 - 6`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, Data(2), scope.Constants[1])
	assert.Equal(t, Data(3), scope.Constants[2])
	assert.Equal(t, Data(4), scope.Constants[4])
	assert.Equal(t, Data(5), scope.Constants[5])
	assert.Equal(t, Data(6), scope.Constants[8])

	// grouping
	expected = []Instruction{
		{Code: INT64_ADD, Left: 1, Right: 2, Out: 3},
		{Code: INT64_MULTIPLY, Left: 3, Right: 4, Out: 5},
	}
	scope, insts = encodeExpression(`(2 + 3) * 4`)
	assert.Equal(t, expected, insts)
}
