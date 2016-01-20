package bytecode

import (
	"testing"

	// TODO: Maybe don't rely on parser et. al. when more code is stable
	"github.com/kestred/philomath/parser"
	"github.com/kestred/philomath/semantics"
	"github.com/stretchr/testify/assert"
)

func encodeExpression(input string) (*Scope, []Instruction) {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	expr := p.ParseExpression()
	semantics.InferTypes(expr)

	scope := &Scope{}
	scope.Init()
	return scope, FromExpr(expr, scope)
}

func TestEncodeArithmetic(t *testing.T) {
	// constants
	constants := []Data{0, 22}
	expected := []Instruction{{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)}}
	scope, insts := encodeExpression(`22`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	constants = []Data{0, 0755}
	expected = []Instruction{{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)}}
	scope, insts = encodeExpression(`0755`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	constants = []Data{0, 0x4000000000000000}
	expected = []Instruction{{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)}}
	scope, insts = encodeExpression(`2.0`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	// add, subtract, multiply, divide
	constants = []Data{0, 2, 3, 4, 5, 6}
	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: I64_MULTIPLY, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Code: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Code: LOAD_CONST, Out: Register(5), Left: Constant(4)},
		{Code: I64_DIVIDE, Out: Register(6), Left: Register(4), Right: Register(5)},
		{Code: I64_ADD, Out: Register(7), Left: Register(3), Right: Register(6)},
		{Code: LOAD_CONST, Out: Register(8), Left: Constant(5)},
		{Code: I64_SUBTRACT, Out: Register(9), Left: Register(7), Right: Register(8)},
	}
	scope, insts = encodeExpression(`2 * 3 + 4 / 5 - 6`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	constants = []Data{0,
		0x4000000000000000,
		0x4010000000000000,
		0x4020000000000000,
		0x4030000000000000,
		0x4040000000000000,
	}
	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: F64_MULTIPLY, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Code: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Code: LOAD_CONST, Out: Register(5), Left: Constant(4)},
		{Code: F64_DIVIDE, Out: Register(6), Left: Register(4), Right: Register(5)},
		{Code: F64_ADD, Out: Register(7), Left: Register(3), Right: Register(6)},
		{Code: LOAD_CONST, Out: Register(8), Left: Constant(5)},
		{Code: F64_SUBTRACT, Out: Register(9), Left: Register(7), Right: Register(8)},
	}
	scope, insts = encodeExpression(`2.0 * 4.0 + 8.0 / 16.0 - 32.0`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	constants = []Data{0, 2, 3, 4, 5, 6}
	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: U64_MULTIPLY, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Code: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Code: LOAD_CONST, Out: Register(5), Left: Constant(4)},
		{Code: U64_DIVIDE, Out: Register(6), Left: Register(4), Right: Register(5)},
		{Code: U64_ADD, Out: Register(7), Left: Register(3), Right: Register(6)},
		{Code: LOAD_CONST, Out: Register(8), Left: Constant(5)},
		{Code: U64_SUBTRACT, Out: Register(9), Left: Register(7), Right: Register(8)},
	}
	scope, insts = encodeExpression(`02 * 03 + 04 / 05 - 06`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	// grouping
	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: I64_ADD, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Code: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Code: I64_MULTIPLY, Out: Register(5), Left: Register(3), Right: Register(4)},
	}
	scope, insts = encodeExpression(`(2 + 3) * 4`)
	assert.Equal(t, expected, insts)

	// conversions
	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: I64_ADD, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Code: CONVERT_I64_TO_F64, Out: Register(4), Left: Register(3)},
		{Code: LOAD_CONST, Out: Register(5), Left: Constant(3)},
		{Code: F64_ADD, Out: Register(6), Left: Register(4), Right: Register(5)},
	}
	scope, insts = encodeExpression(`(2 + 3) + 4.0`)
	assert.Equal(t, expected, insts)

	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: CONVERT_I64_TO_F64, Out: Register(2), Left: Register(1)},
		{Code: LOAD_CONST, Out: Register(3), Left: Constant(2)},
		{Code: F64_ADD, Out: Register(4), Left: Register(2), Right: Register(3)},
		{Code: LOAD_CONST, Out: Register(5), Left: Constant(3)},
		{Code: CONVERT_I64_TO_F64, Out: Register(6), Left: Register(5)},
		{Code: F64_ADD, Out: Register(7), Left: Register(4), Right: Register(6)},
	}
	scope, insts = encodeExpression(`(2 + 3.0) + 4`)
	assert.Equal(t, expected, insts)

	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: U64_ADD, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Code: CONVERT_U64_TO_F64, Out: Register(4), Left: Register(3)},
		{Code: LOAD_CONST, Out: Register(5), Left: Constant(3)},
		{Code: F64_ADD, Out: Register(6), Left: Register(4), Right: Register(5)},
	}
	scope, insts = encodeExpression(`(02 + 03) + 4.0`)
	assert.Equal(t, expected, insts)

	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: CONVERT_U64_TO_F64, Out: Register(2), Left: Register(1)},
		{Code: LOAD_CONST, Out: Register(3), Left: Constant(2)},
		{Code: F64_ADD, Out: Register(4), Left: Register(2), Right: Register(3)},
		{Code: LOAD_CONST, Out: Register(5), Left: Constant(3)},
		{Code: CONVERT_U64_TO_F64, Out: Register(6), Left: Register(5)},
		{Code: F64_ADD, Out: Register(7), Left: Register(4), Right: Register(6)},
	}
	scope, insts = encodeExpression(`(02 + 3.0) + 04`)
	assert.Equal(t, expected, insts)
}
