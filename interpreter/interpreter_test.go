package interpreter

import (
	"testing"

	// TODO: Maybe stop relying on parser et. al. when more code is stable?
	"github.com/kestred/philomath/bytecode"
	"github.com/kestred/philomath/parser"
	"github.com/kestred/philomath/semantics"
	"github.com/stretchr/testify/assert"
)

func evalExpression(t *testing.T, input string) bytecode.Data {
	p := parser.Parser{}
	p.Init("example", false, []byte(input))
	expr := p.ParseExpression()
	if len(p.Errors) > 0 {
		t.Fatalf("Unexpected parse error\n\n%v", p.Errors[0].Error())
	}

	semantics.InferTypes(expr)
	scope := &bytecode.Scope{}
	scope.Init()
	insts := bytecode.FromExpr(expr, scope)
	return Evaluate(insts, scope.Constants, scope.NextRegister)
}

func evalBlock(t *testing.T, input string) bytecode.Data {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	block := p.ParseBlock()
	if len(p.Errors) > 0 {
		t.Fatalf("Unexpected parse error\n\n%v", p.Errors[0].Error())
	}

	semantics.InferTypes(block)
	scope := &bytecode.Scope{}
	scope.Init()
	insts := bytecode.FromBlock(block, scope)
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
		{Code: bytecode.I64_ADD, Left: 1, Right: 2, Out: 3},
	}
	consts = []bytecode.Data{0, 1, 2}
	result = Evaluate(insts, consts, 3)
	assert.Equal(t, 3, int(result))
}

func TestEvaluateArithmetic(t *testing.T) {
	// constant
	result := evalExpression(t, `22`)
	assert.Equal(t, int64(22), bytecode.ToI64(result))

	// add, subtract, multiply, divide
	result = evalExpression(t, `2 * 3 + 27 / 9 - 15`)
	assert.Equal(t, int64(2*3+27/9-15), bytecode.ToI64(result))
	assert.Equal(t, int64(-6), bytecode.ToI64(result))

	result = evalExpression(t, `2.0 * 4.0 + 8.0 / 16.0 - 32.0`)
	assert.Equal(t, float64(2.0*4.0+8.0/16.0-32.0), bytecode.ToF64(result))
	assert.Equal(t, float64(-23.5), bytecode.ToF64(result))

	result = evalExpression(t, `02 * 03 + 04 / 05 - 01`)
	assert.Equal(t, uint64(02*03+04/05-01), bytecode.ToU64(result))
	assert.Equal(t, uint64(5), bytecode.ToU64(result))

	result = evalExpression(t, `(2 + 3) + 4.0`)
	assert.Equal(t, float64((2+3)+4.0), bytecode.ToF64(result))
	assert.Equal(t, float64(9.0), bytecode.ToF64(result))

	result = evalExpression(t, `(2 + 3.0) + 4`)
	assert.Equal(t, float64((2+3.0)+4), bytecode.ToF64(result))
	assert.Equal(t, float64(9.0), bytecode.ToF64(result))

	result = evalExpression(t, `(02 + 03) + 4.0`)
	assert.Equal(t, float64((02+03)+4.0), bytecode.ToF64(result))
	assert.Equal(t, float64(9.0), bytecode.ToF64(result))

	result = evalExpression(t, `(02 + 3.0) + 04`)
	assert.Equal(t, float64((02+3.0)+04), bytecode.ToF64(result))
	assert.Equal(t, float64(9.0), bytecode.ToF64(result))
}

func TestEncodeBlock(t *testing.T) {
	result := evalBlock(t, `{
		hoge :: 3;          # constant decl
		hoge + 2;           # one ident in expr, result ignored

		piyo := 0.5 * hoge; # mutable decl
		piyo / hoge;        # two ident in expr; TODO: return statement
	}`)

	const hoge = 3
	var piyo = 0.5 * float64(hoge)
	assert.Equal(t, float64(piyo/float64(hoge)), bytecode.ToF64(result))
	assert.Equal(t, float64(0.5), bytecode.ToF64(result))
}
