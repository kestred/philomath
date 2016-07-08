package interpreter

import (
	"testing"

	"github.com/kestred/philomath/code/parser"
	"github.com/kestred/philomath/code/semantics"
	"github.com/stretchr/testify/assert"

	bc "github.com/kestred/philomath/code/bytecode"
)

func evalExample(t *testing.T, input string) []byte {
	p := parser.Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	section := semantics.FlattenTree(node, nil)
	semantics.ResolveNames(&section)
	semantics.InferTypes(&section)
	program := bc.NewProgram()
	program.Extend(node)

	t.Log(program.Procedures[0].Instructions)
	return Evaluate(program.Procedures[0])
}

func TestEvaluateNoop(t *testing.T) {
	var program *bc.Program
	var result []byte

	// just a noop
	program = bc.NewProgram()
	program.Procedures[0].NextFree = 1
	program.Procedures[0].Instructions = []bc.Instruction{{Op: bc.NOOP}}
	result = Evaluate(program.Procedures[0])
	assert.Equal(t, []byte(nil), result)

	// interleaved noops
	program = bc.NewProgram()
	program.Data = map[string][]byte{".LC1": bc.Pack(uint64(1)), ".LC2": bc.Pack(uint64(2))}
	program.Procedures[0].PrevResult = bc.Rg(3, bc.Int64)
	program.Procedures[0].NextFree = 4
	program.Procedures[0].Instructions = []bc.Instruction{
		{bc.NOOP, nil},
		{bc.NOOP, nil},
		{bc.NOOP, nil},
		{bc.LOAD, bc.Constant(".LC1", bc.Rg(1, bc.Int64))},
		{bc.NOOP, nil},
		{bc.LOAD, bc.Constant(".LC2", bc.Rg(2, bc.Int64))},
		{bc.NOOP, nil},
		{bc.ADD, bc.Binary(bc.Rg(1, bc.Int64), bc.Rg(2, bc.Int64), bc.Rg(3, bc.Int64))},
	}
	result = Evaluate(program.Procedures[0])
	assert.Equal(t, bc.Pack(int64(3)), result)
}

func TestEvaluateArithmetic(t *testing.T) {
	// constant
	result := evalExample(t, `22;`)
	assert.Equal(t, bc.Pack(int64(22)), result)

	// add, subtract, multiply, divide
	result = evalExample(t, `2 * 3 + 27 / 9 - 15;`)
	assert.Equal(t, bc.Pack(int64(2*3+27/9-15)), result)
	assert.Equal(t, bc.Pack(int64(-6)), result)

	result = evalExample(t, `2.0 * 4.0 + 8.0 / 16.0 - 32.0;`)
	assert.Equal(t, bc.Pack(float64(2.0*4.0+8.0/16.0-32.0)), result)
	assert.Equal(t, bc.Pack(float64(-23.5)), result)

	result = evalExample(t, `02 * 03 + 04 / 05 - 01;`)
	assert.Equal(t, bc.Pack(uint64(02*03+04/05-01)), result)
	assert.Equal(t, bc.Pack(uint64(5)), result)

	result = evalExample(t, `(2 + 3) + 4.0;`)
	assert.Equal(t, bc.Pack(float64((2+3)+4.0)), result)
	assert.Equal(t, bc.Pack(float64(9.0)), result)

	result = evalExample(t, `(2 + 3.0) + 4;`)
	assert.Equal(t, bc.Pack(float64((2+3.0)+4)), result)
	assert.Equal(t, bc.Pack(float64(9.0)), result)

	result = evalExample(t, `(02 + 03) + 4.0;`)
	assert.Equal(t, bc.Pack(float64((02+03)+4.0)), result)
	assert.Equal(t, bc.Pack(float64(9.0)), result)

	result = evalExample(t, `(02 + 3.0) + 04;`)
	assert.Equal(t, bc.Pack(float64((02+3.0)+04)), result)
	assert.Equal(t, bc.Pack(float64(9.0)), result)
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
	assert.Equal(t, bc.Pack(float64(piyo/float64(hoge))), result)
	assert.Equal(t, bc.Pack(float64(0.5)), result)

	// assignment Statement
	result = evalExample(t, `{
		xyzzy := 012;
		xyzzy = 0700;
		xyzzy + 0;
	}`)
	var xyzzy = uint64(012)
	xyzzy = uint64(0700)
	assert.Equal(t, bc.Pack(uint64(xyzzy)), result)
	assert.Equal(t, bc.Pack(uint64(0700)), result)

	// assignment with cast
	result = evalExample(t, `{
		plugh := 1 - 37;
		plugh = 0.25 * plugh;
		plugh + 0;
	}`)
	var plugh = int64(1) - int64(37)
	plugh = int64(0.25 * float64(plugh))
	assert.Equal(t, bc.Pack(int64(plugh)), result)
	assert.Equal(t, bc.Pack(int64(-9)), result)

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
	assert.Equal(t, bc.Pack(int64(nerrf)), result)
	assert.Equal(t, bc.Pack(int64(0700/5)), result)
}
