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
	constants := []Data{0, 22}
	expected := []Instruction{{Code: LOAD_CONST, Out: 1, Left: 1}}
	scope, insts := encodeExpression(`22`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	// add, subtract, multiply, divide
	constants = []Data{0, 2, 3, 4, 5, 6}
	expected = []Instruction{
		{Code: LOAD_CONST, Out: 1, Left: 1},
		{Code: LOAD_CONST, Out: 2, Left: 2},
		{Code: INT64_MULTIPLY, Left: 1, Right: 2, Out: 3},
		{Code: LOAD_CONST, Out: 4, Left: 3},
		{Code: LOAD_CONST, Out: 5, Left: 4},
		{Code: INT64_DIVIDE, Left: 4, Right: 5, Out: 6},
		{Code: INT64_ADD, Left: 3, Right: 6, Out: 7},
		{Code: LOAD_CONST, Out: 8, Left: 5},
		{Code: INT64_SUBTRACT, Left: 7, Right: 8, Out: 9},
	}
	scope, insts = encodeExpression(`2 * 3 + 4 / 5 - 6`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	// grouping
	expected = []Instruction{
		{Code: LOAD_CONST, Out: 1, Left: 1},
		{Code: LOAD_CONST, Out: 2, Left: 2},
		{Code: INT64_ADD, Left: 1, Right: 2, Out: 3},
		{Code: LOAD_CONST, Out: 4, Left: 3},
		{Code: INT64_MULTIPLY, Left: 3, Right: 4, Out: 5},
	}
	scope, insts = encodeExpression(`(2 + 3) * 4`)
	assert.Equal(t, expected, insts)
}
