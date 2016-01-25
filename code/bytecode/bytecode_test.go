package bytecode

import (
	"testing"

	// TODO: Maybe don't rely on parser et. al. when more code is stable
	"github.com/kestred/philomath/code/code"
	"github.com/kestred/philomath/code/parser"
	"github.com/kestred/philomath/code/semantics"
	"github.com/stretchr/testify/assert"
)

func encodeExpression(t *testing.T, input string) (*Scope, []Instruction) {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	expr := p.ParseExpression()
	if len(p.Errors) > 0 {
		t.Fatalf("Unexpected parse error\n\n%v", p.Errors[0].Error())
	}

	section := code.PrepareTree(expr, nil)
	semantics.InferTypes(&section)
	scope := &Scope{}
	scope.Init()
	return scope, FromExpr(expr, scope)
}

func encodeBlock(t *testing.T, input string) (*Scope, []Instruction) {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	block := p.ParseBlock()
	if len(p.Errors) > 0 {
		t.Fatalf("Unexpected parse error\n\n%v", p.Errors[0].Error())
	}

	section := code.PrepareTree(block, nil)
	semantics.ResolveNames(&section)
	semantics.InferTypes(&section)
	scope := &Scope{}
	scope.Init()
	return scope, FromBlock(block, scope)
}

func TestUnsafeSafety(t *testing.T) {
	// from Number (to Data)
	ldata := FromU64(2)
	idata := FromI64(-2)
	rdata := FromF64(2.0)
	assert.NotEqual(t, ldata, idata)
	assert.NotEqual(t, idata, rdata)
	assert.NotEqual(t, ldata, rdata)
	assert.Equal(t, Data(2), ldata)
	assert.Equal(t, Data(^uint(0)-1), idata)
	assert.Equal(t, Data(0x4000000000000000), rdata)

	// (from Data) to Number
	lhs := ToI64(Data(2))
	rhs := ToF64(Data(2))
	assert.NotEqual(t, float64(lhs), float64(rhs))
	assert.Equal(t, int64(2), lhs)
	assert.Equal(t, float64(1e-323), rhs)

	// round-trip
	assert.Equal(t, float64(-37.84), ToF64(FromF64(-37.84)))
	assert.Equal(t, int64(-488), ToI64(FromI64(-488)))
	assert.Equal(t, uint64(1099511627776), ToU64(FromU64(1099511627776)))
}

func TestEncodeArithmetic(t *testing.T) {
	// constants
	constants := []Data{0, 22}
	expected := []Instruction{{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)}}
	scope, insts := encodeExpression(t, `22`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	constants = []Data{0, FromU64(0755)}
	expected = []Instruction{{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)}}
	scope, insts = encodeExpression(t, `0755`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	constants = []Data{0, FromF64(2.0)}
	expected = []Instruction{{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)}}
	scope, insts = encodeExpression(t, `2.0`)
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
	scope, insts = encodeExpression(t, `2 * 3 + 4 / 5 - 6`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	constants = []Data{0, FromF64(2.0), FromF64(4.0), FromF64(8.0), FromF64(16.0), FromF64(32.0)}
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
	scope, insts = encodeExpression(t, `2.0 * 4.0 + 8.0 / 16.0 - 32.0`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	constants = []Data{0, FromU64(2), FromU64(3), FromU64(4), FromU64(5), FromU64(6)}
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
	scope, insts = encodeExpression(t, `02 * 03 + 04 / 05 - 06`)
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
	_, insts = encodeExpression(t, `(2 + 3) * 4`)
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
	_, insts = encodeExpression(t, `(2 + 3) + 4.0`)
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
	_, insts = encodeExpression(t, `(2 + 3.0) + 4`)
	assert.Equal(t, expected, insts)

	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: U64_ADD, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Code: CONVERT_U64_TO_F64, Out: Register(4), Left: Register(3)},
		{Code: LOAD_CONST, Out: Register(5), Left: Constant(3)},
		{Code: F64_ADD, Out: Register(6), Left: Register(4), Right: Register(5)},
	}
	_, insts = encodeExpression(t, `(02 + 03) + 4.0`)
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
	_, insts = encodeExpression(t, `(02 + 3.0) + 04`)
	assert.Equal(t, expected, insts)
}

func TestEncodeBlock(t *testing.T) {
	// Declarations
	constants := []Data{
		0: 0,
		1: 3,
		2: 2,
		3: FromF64(0.5),
	}
	expected := []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: NOOP, Out: Register(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: I64_ADD, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Code: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Code: NOOP, Out: Register(1)},
		{Code: CONVERT_I64_TO_F64, Out: Register(5), Left: Register(1)},
		{Code: F64_MULTIPLY, Out: Register(6), Left: Register(4), Right: Register(5)},
		{Code: NOOP, Out: Register(6)},
		{Code: NOOP, Out: Register(1)},
		{Code: CONVERT_I64_TO_F64, Out: Register(7), Left: Register(1)},
		{Code: F64_DIVIDE, Out: Register(8), Left: Register(6), Right: Register(7)},
	}
	scope, insts := encodeBlock(t, `{
		hoge :: 3;          # constant decl
		hoge + 2;           # one ident in expr, result ignored

		piyo := 0.5 * hoge; # mutable decl
		piyo / hoge;        # two ident in expr; TODO: return statement
	}`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	// Simple and Parallel Assignment
	constants = []Data{
		0: 0,
		1: 1,
		2: 4,
		3: FromU64(012),
		4: 14,
		5: FromU64(0700),
		6: FromF64(0.25),
		7: FromF64(5.0),
		8: 10000,
		9: 100,
	}
	expected = []Instruction{
		// plugh := 1 - 3;
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: I64_SUBTRACT, Out: Register(3), Left: Register(1), Right: Register(2)},
		// xyzzy := 012;
		{Code: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		// nerrf := 14;
		{Code: LOAD_CONST, Out: Register(5), Left: Constant(4)},
		// xyzzy = 0700;
		{Code: LOAD_CONST, Out: Register(6), Left: Constant(5)},
		{Code: COPY_VALUE, Out: Register(4), Left: Register(6)},
		// plugh = 0.25 * plugh;
		{Code: LOAD_CONST, Out: Register(7), Left: Constant(6)},
		{Code: NOOP, Out: Register(3)},
		{Code: CONVERT_I64_TO_F64, Out: Register(8), Left: Register(3)},
		{Code: F64_MULTIPLY, Out: Register(9), Left: Register(7), Right: Register(8)},
		{Code: CONVERT_F64_TO_I64, Out: Register(10), Left: Register(9)},
		{Code: COPY_VALUE, Out: Register(3), Left: Register(10)},
		// xyzzy, nerrf, plugh = plugh, (xyzzy / 5.0), nerrf;
		{Code: NOOP, Out: Register(3)},
		{Code: COPY_VALUE, Out: Register(11), Left: Register(3)},
		{Code: NOOP, Out: Register(4)},
		{Code: CONVERT_U64_TO_F64, Out: Register(12), Left: Register(4)},
		{Code: LOAD_CONST, Out: Register(13), Left: Constant(7)},
		{Code: F64_DIVIDE, Out: Register(14), Left: Register(12), Right: Register(13)},
		{Code: COPY_VALUE, Out: Register(15), Left: Register(14)},
		{Code: NOOP, Out: Register(5)},
		{Code: COPY_VALUE, Out: Register(16), Left: Register(5)},
		{Code: COPY_VALUE, Out: Register(4), Left: Register(11)},
		{Code: CONVERT_F64_TO_I64, Out: Register(17), Left: Register(15)},
		{Code: COPY_VALUE, Out: Register(5), Left: Register(17)},
		{Code: COPY_VALUE, Out: Register(3), Left: Register(16)},
		// barrf := xyzzy * 10000 + nerrf * 100;
		{Code: NOOP, Out: Register(4)},
		{Code: LOAD_CONST, Out: Register(18), Left: Constant(8)},
		{Code: U64_MULTIPLY, Out: Register(19), Left: Register(4), Right: Register(18)},
		{Code: NOOP, Out: Register(5)},
		{Code: LOAD_CONST, Out: Register(20), Left: Constant(9)},
		{Code: I64_MULTIPLY, Out: Register(21), Left: Register(5), Right: Register(20)},
		{Code: U64_ADD, Out: Register(22), Left: Register(19), Right: Register(21)},
	}
	// FIXME: This example will break with more strict with implict casts
	scope, insts = encodeBlock(t, `{
		plugh := 1 - 4;
		xyzzy := 012;
		nerrf := 14;

		xyzzy = 0700;                      # assignment
		plugh = 0.25 * plugh;              # assignment with cast

		# parallel assignment (with and without casts)
		xyzzy, nerrf, plugh = plugh, (xyzzy / 5.0), nerrf;
		barrf := xyzzy * 10000 + nerrf * 100;
	}`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)

	// Nested block
	constants = []Data{
		0: 0,
		1: FromU64(0600),
		2: FromF64(6.29),
		3: 2,
		4: 3,
	}
	expected = []Instruction{
		{Code: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Code: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Code: NOOP, Out: Register(2)},
		{Code: LOAD_CONST, Out: Register(3), Left: Constant(3)},
		{Code: CONVERT_I64_TO_F64, Out: Register(4), Left: Register(3)},
		{Code: F64_DIVIDE, Out: Register(5), Left: Register(2), Right: Register(4)},
		{Code: NOOP, Out: Register(5)},
		{Code: NOOP, Out: Register(1)},
		{Code: CONVERT_U64_TO_F64, Out: Register(6), Left: Register(1)},
		{Code: F64_SUBTRACT, Out: Register(7), Left: Register(5), Right: Register(6)},
		{Code: LOAD_CONST, Out: Register(8), Left: Constant(4)},
		{Code: COPY_VALUE, Out: Register(1), Left: Register(8)},
		{Code: NOOP, Out: Register(2)},
	}
	scope, insts = encodeBlock(t, `{
		ham  := 0600;
		eggs :: 6.29;

		{
			spam := eggs / 2;
			spam - ham;
			ham = 3;
		}

		eggs;
	}`)
	assert.Equal(t, expected, insts)
	assert.Equal(t, constants, scope.Constants)
}
