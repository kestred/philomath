package interpreter

import (
	"testing"

	"github.com/kestred/philomath/code/bytecode"
	"github.com/kestred/philomath/code/code"
	"github.com/kestred/philomath/code/parser"
	"github.com/kestred/philomath/code/semantics"
	"github.com/stretchr/testify/assert"
)

func evalExample(t *testing.T, input string) bytecode.Data {
	p := parser.Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	section := code.PrepareTree(node, nil)
	semantics.ResolveNames(&section)
	semantics.InferTypes(&section)
	program := bytecode.NewProgram()
	program.Extend(node)

	t.Log(program.Procedures[0].Instructions)
	return Evaluate(program.Procedures[0])
}

func TestEvaluateNoop(t *testing.T) {
	var program *bytecode.Program
	var result bytecode.Data

	// just a noop
	program = bytecode.NewProgram()
	program.Procedures[0].NextRegister = 1
	program.Procedures[0].Instructions = []bytecode.Instruction{{Op: bytecode.NOOP}}
	result = Evaluate(program.Procedures[0])
	assert.Equal(t, 0, int(result))

	// interleaved noops
	program = bytecode.NewProgram()
	program.Constants = []bytecode.Data{0, 1, 2}
	program.Procedures[0].ExprRegister = 3
	program.Procedures[0].NextRegister = 4
	program.Procedures[0].Instructions = []bytecode.Instruction{
		{Op: bytecode.NOOP},
		{Op: bytecode.NOOP},
		{Op: bytecode.NOOP},
		{Op: bytecode.LOAD_CONST, Out: 1, Left: 1},
		{Op: bytecode.NOOP},
		{Op: bytecode.LOAD_CONST, Out: 2, Left: 2},
		{Op: bytecode.NOOP},
		{Op: bytecode.I64_ADD, Out: 3, Left: 1, Right: 2},
	}
	result = Evaluate(program.Procedures[0])
	assert.Equal(t, 3, int(result))
}

func TestEvaluateArithmetic(t *testing.T) {
	// constant
	result := evalExample(t, `22;`)
	assert.Equal(t, int64(22), bytecode.ToI64(result))

	// add, subtract, multiply, divide
	result = evalExample(t, `2 * 3 + 27 / 9 - 15;`)
	assert.Equal(t, int64(2*3+27/9-15), bytecode.ToI64(result))
	assert.Equal(t, int64(-6), bytecode.ToI64(result))

	result = evalExample(t, `2.0 * 4.0 + 8.0 / 16.0 - 32.0;`)
	assert.Equal(t, float64(2.0*4.0+8.0/16.0-32.0), bytecode.ToF64(result))
	assert.Equal(t, float64(-23.5), bytecode.ToF64(result))

	result = evalExample(t, `02 * 03 + 04 / 05 - 01;`)
	assert.Equal(t, uint64(02*03+04/05-01), bytecode.ToU64(result))
	assert.Equal(t, uint64(5), bytecode.ToU64(result))

	result = evalExample(t, `(2 + 3) + 4.0;`)
	assert.Equal(t, float64((2+3)+4.0), bytecode.ToF64(result))
	assert.Equal(t, float64(9.0), bytecode.ToF64(result))

	result = evalExample(t, `(2 + 3.0) + 4;`)
	assert.Equal(t, float64((2+3.0)+4), bytecode.ToF64(result))
	assert.Equal(t, float64(9.0), bytecode.ToF64(result))

	result = evalExample(t, `(02 + 03) + 4.0;`)
	assert.Equal(t, float64((02+03)+4.0), bytecode.ToF64(result))
	assert.Equal(t, float64(9.0), bytecode.ToF64(result))

	result = evalExample(t, `(02 + 3.0) + 04;`)
	assert.Equal(t, float64((02+3.0)+04), bytecode.ToF64(result))
	assert.Equal(t, float64(9.0), bytecode.ToF64(result))
}

func TestEncodeBlock(t *testing.T) {
	// declarations
	result := evalExample(t, `{
		hoge :: 3;          // constant decl
		hoge + 2;           // one ident in expr, result ignored

		piyo := 0.5 * hoge; // mutable decl
		piyo / hoge;        // two ident in expr; TODO: return statement
	}`)
	const hoge = 3
	var piyo = 0.5 * float64(hoge)
	assert.Equal(t, float64(piyo/float64(hoge)), bytecode.ToF64(result))
	assert.Equal(t, float64(0.5), bytecode.ToF64(result))

	// assignment Statement
	result = evalExample(t, `{
		xyzzy := 012;
		xyzzy = 0700;
		xyzzy + 0;
	}`)
	var xyzzy = uint64(012)
	xyzzy = uint64(0700)
	assert.Equal(t, uint64(xyzzy), bytecode.ToU64(result))
	assert.Equal(t, uint64(0700), bytecode.ToU64(result))

	// assignment with cast
	result = evalExample(t, `{
		plugh := 1 - 37;
		plugh = 0.25 * plugh;
		plugh + 0;
	}`)
	var plugh = int64(1) - int64(37)
	plugh = int64(0.25 * float64(plugh))
	assert.Equal(t, int64(plugh), bytecode.ToI64(result))
	assert.Equal(t, int64(-9), bytecode.ToI64(result))

	// parallel assignment (with and without casts)
	result = evalExample(t, `{
		plugh := 1 - 37;
		xyzzy := 012;
		nerrf := 14;

		xyzzy = 0700;
		plugh = 0.25 * plugh;

		// parallel assignment (with and without casts)
		xyzzy, nerrf, plugh = plugh, (xyzzy / 5.0), nerrf;
		nerrf + 0;
	}`)

	var nerrf = int64(14)
	tmp1 := uint64(plugh)
	tmp2 := int64(float64(xyzzy) / 5.0)
	tmp3 := int64(nerrf)
	xyzzy = tmp1
	nerrf = tmp2
	plugh = tmp3
	assert.Equal(t, int64(nerrf), bytecode.ToI64(result))
	assert.Equal(t, int64(0700/5), bytecode.ToI64(result))
}
