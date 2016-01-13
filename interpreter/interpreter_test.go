package interpreter

import (
	"testing"

	// TODO: Stop relying on parser when more code is stable
	"github.com/kestred/philomath/bytecode"
	"github.com/kestred/philomath/parser"
	"github.com/stretchr/testify/assert"
)

func evalExpression(input string) bytecode.Data {
	p := parser.Parser{}
	p.Init("test", false, []byte(input))
	expr := p.ParseExpression()
	scope := &bytecode.Scope{}
	scope.Init()
	insts := bytecode.FromExpr(expr, scope)
	return Evaluate(insts, scope.Constants, scope.NextRegister)
}

func TestEvaluateNoop(t *testing.T) {
	insts := []bytecode.Instruction{{Code: bytecode.NOOP}}
	consts := []bytecode.Data{0}
	result := Evaluate(insts, consts, 1)
	assert.Equal(t, 0, int(result))

	insts = []bytecode.Instruction{
		{Code: bytecode.NOOP},
		{Code: bytecode.NOOP},
		{Code: bytecode.NOOP},
		{Code: bytecode.LOAD_CONST, Out: 1, Left: 1},
		{Code: bytecode.NOOP},
		{Code: bytecode.LOAD_CONST, Out: 2, Left: 2},
		{Code: bytecode.NOOP},
		{Code: bytecode.INT64_ADD, Left: 1, Right: 2, Out: 3},
	}
	consts = []bytecode.Data{0, 1, 2}
	result = Evaluate(insts, consts, 3)
	assert.Equal(t, 3, int(result))
}

func TestEvaluateArithmetic(t *testing.T) {
	// constant
	result := evalExpression(`22`)
	assert.Equal(t, 22, int(result))

	// add, subtract, multiply, divide
	result = evalExpression(`2 * 3 + 27 / 9 - 15`)
	assert.Equal(t, 2*3+27/9-15, int(result))
	assert.Equal(t, -6, int(result))
}
