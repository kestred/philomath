package bytecode

import (
	"testing"

	// TODO: Maybe don't rely on parser et. al. when more code is stable
	"github.com/kestred/philomath/code/code"
	"github.com/kestred/philomath/code/parser"
	"github.com/kestred/philomath/code/semantics"
	"github.com/stretchr/testify/assert"
)

func generateBytecode(t *testing.T, input string) Program {
	p := parser.Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	section := code.PrepareTree(node, nil)
	semantics.ResolveNames(&section)
	semantics.InferTypes(&section)
	program, _ := Generate(node)
	return program
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
	expected := []Instruction{{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)}}
	program := generateBytecode(t, `22;`)
	assert.Equal(t, expected, program.Instructions)
	assert.Equal(t, constants, program.Constants)

	constants = []Data{0, FromU64(0755)}
	expected = []Instruction{{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)}}
	program = generateBytecode(t, `0755;`)
	assert.Equal(t, expected, program.Instructions)
	assert.Equal(t, constants, program.Constants)

	constants = []Data{0, FromF64(2.0)}
	expected = []Instruction{{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)}}
	program = generateBytecode(t, `2.0;`)
	assert.Equal(t, expected, program.Instructions)
	assert.Equal(t, constants, program.Constants)

	// add, subtract, multiply, divide
	constants = []Data{0, 2, 3, 4, 5, 6}
	expected = []Instruction{
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Op: I64_MULTIPLY, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Op: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Op: LOAD_CONST, Out: Register(5), Left: Constant(4)},
		{Op: I64_DIVIDE, Out: Register(6), Left: Register(4), Right: Register(5)},
		{Op: I64_ADD, Out: Register(7), Left: Register(3), Right: Register(6)},
		{Op: LOAD_CONST, Out: Register(8), Left: Constant(5)},
		{Op: I64_SUBTRACT, Out: Register(9), Left: Register(7), Right: Register(8)},
	}
	program = generateBytecode(t, `2 * 3 + 4 / 5 - 6;`)
	assert.Equal(t, expected, program.Instructions)
	assert.Equal(t, constants, program.Constants)

	constants = []Data{0, FromF64(2.0), FromF64(4.0), FromF64(8.0), FromF64(16.0), FromF64(32.0)}
	expected = []Instruction{
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Op: F64_MULTIPLY, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Op: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Op: LOAD_CONST, Out: Register(5), Left: Constant(4)},
		{Op: F64_DIVIDE, Out: Register(6), Left: Register(4), Right: Register(5)},
		{Op: F64_ADD, Out: Register(7), Left: Register(3), Right: Register(6)},
		{Op: LOAD_CONST, Out: Register(8), Left: Constant(5)},
		{Op: F64_SUBTRACT, Out: Register(9), Left: Register(7), Right: Register(8)},
	}
	program = generateBytecode(t, `2.0 * 4.0 + 8.0 / 16.0 - 32.0;`)
	assert.Equal(t, expected, program.Instructions)
	assert.Equal(t, constants, program.Constants)

	constants = []Data{0, FromU64(2), FromU64(3), FromU64(4), FromU64(5), FromU64(6)}
	expected = []Instruction{
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Op: U64_MULTIPLY, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Op: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Op: LOAD_CONST, Out: Register(5), Left: Constant(4)},
		{Op: U64_DIVIDE, Out: Register(6), Left: Register(4), Right: Register(5)},
		{Op: U64_ADD, Out: Register(7), Left: Register(3), Right: Register(6)},
		{Op: LOAD_CONST, Out: Register(8), Left: Constant(5)},
		{Op: U64_SUBTRACT, Out: Register(9), Left: Register(7), Right: Register(8)},
	}
	program = generateBytecode(t, `02 * 03 + 04 / 05 - 06;`)
	assert.Equal(t, expected, program.Instructions)
	assert.Equal(t, constants, program.Constants)

	// grouping
	expected = []Instruction{
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Op: I64_ADD, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Op: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Op: I64_MULTIPLY, Out: Register(5), Left: Register(3), Right: Register(4)},
	}
	program = generateBytecode(t, `(2 + 3) * 4;`)
	assert.Equal(t, expected, program.Instructions)

	// conversions
	expected = []Instruction{
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Op: I64_ADD, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Op: CONVERT_I64_TO_F64, Out: Register(4), Left: Register(3)},
		{Op: LOAD_CONST, Out: Register(5), Left: Constant(3)},
		{Op: F64_ADD, Out: Register(6), Left: Register(4), Right: Register(5)},
	}
	program = generateBytecode(t, `(2 + 3) + 4.0;`)
	assert.Equal(t, expected, program.Instructions)

	expected = []Instruction{
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: CONVERT_I64_TO_F64, Out: Register(2), Left: Register(1)},
		{Op: LOAD_CONST, Out: Register(3), Left: Constant(2)},
		{Op: F64_ADD, Out: Register(4), Left: Register(2), Right: Register(3)},
		{Op: LOAD_CONST, Out: Register(5), Left: Constant(3)},
		{Op: CONVERT_I64_TO_F64, Out: Register(6), Left: Register(5)},
		{Op: F64_ADD, Out: Register(7), Left: Register(4), Right: Register(6)},
	}
	program = generateBytecode(t, `(2 + 3.0) + 4;`)
	assert.Equal(t, expected, program.Instructions)

	expected = []Instruction{
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Op: U64_ADD, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Op: CONVERT_U64_TO_F64, Out: Register(4), Left: Register(3)},
		{Op: LOAD_CONST, Out: Register(5), Left: Constant(3)},
		{Op: F64_ADD, Out: Register(6), Left: Register(4), Right: Register(5)},
	}
	program = generateBytecode(t, `(02 + 03) + 4.0;`)
	assert.Equal(t, expected, program.Instructions)

	expected = []Instruction{
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: CONVERT_U64_TO_F64, Out: Register(2), Left: Register(1)},
		{Op: LOAD_CONST, Out: Register(3), Left: Constant(2)},
		{Op: F64_ADD, Out: Register(4), Left: Register(2), Right: Register(3)},
		{Op: LOAD_CONST, Out: Register(5), Left: Constant(3)},
		{Op: CONVERT_U64_TO_F64, Out: Register(6), Left: Register(5)},
		{Op: F64_ADD, Out: Register(7), Left: Register(4), Right: Register(6)},
	}
	program = generateBytecode(t, `(02 + 3.0) + 04;`)
	assert.Equal(t, expected, program.Instructions)
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
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: NOOP, Out: Register(1)},
		{Op: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Op: I64_ADD, Out: Register(3), Left: Register(1), Right: Register(2)},
		{Op: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		{Op: NOOP, Out: Register(1)},
		{Op: CONVERT_I64_TO_F64, Out: Register(5), Left: Register(1)},
		{Op: F64_MULTIPLY, Out: Register(6), Left: Register(4), Right: Register(5)},
		{Op: NOOP, Out: Register(6)},
		{Op: NOOP, Out: Register(1)},
		{Op: CONVERT_I64_TO_F64, Out: Register(7), Left: Register(1)},
		{Op: F64_DIVIDE, Out: Register(8), Left: Register(6), Right: Register(7)},
	}
	program := generateBytecode(t, `{
		hoge :: 3;          # constant decl
		hoge + 2;           # one ident in expr, result ignored

		piyo := 0.5 * hoge; # mutable decl
		piyo / hoge;        # two ident in expr; TODO: return statement
	}`)
	assert.Equal(t, expected, program.Instructions)
	assert.Equal(t, constants, program.Constants)

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
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Op: I64_SUBTRACT, Out: Register(3), Left: Register(1), Right: Register(2)},
		// xyzzy := 012;
		{Op: LOAD_CONST, Out: Register(4), Left: Constant(3)},
		// nerrf := 14;
		{Op: LOAD_CONST, Out: Register(5), Left: Constant(4)},
		// xyzzy = 0700;
		{Op: LOAD_CONST, Out: Register(6), Left: Constant(5)},
		{Op: COPY_VALUE, Out: Register(4), Left: Register(6)},
		// plugh = 0.25 * plugh;
		{Op: LOAD_CONST, Out: Register(7), Left: Constant(6)},
		{Op: NOOP, Out: Register(3)},
		{Op: CONVERT_I64_TO_F64, Out: Register(8), Left: Register(3)},
		{Op: F64_MULTIPLY, Out: Register(9), Left: Register(7), Right: Register(8)},
		{Op: CONVERT_F64_TO_I64, Out: Register(10), Left: Register(9)},
		{Op: COPY_VALUE, Out: Register(3), Left: Register(10)},
		// xyzzy, nerrf, plugh = plugh, (xyzzy / 5.0), nerrf;
		{Op: NOOP, Out: Register(3)},
		{Op: COPY_VALUE, Out: Register(11), Left: Register(3)},
		{Op: NOOP, Out: Register(4)},
		{Op: CONVERT_U64_TO_F64, Out: Register(12), Left: Register(4)},
		{Op: LOAD_CONST, Out: Register(13), Left: Constant(7)},
		{Op: F64_DIVIDE, Out: Register(14), Left: Register(12), Right: Register(13)},
		{Op: COPY_VALUE, Out: Register(15), Left: Register(14)},
		{Op: NOOP, Out: Register(5)},
		{Op: COPY_VALUE, Out: Register(16), Left: Register(5)},
		{Op: COPY_VALUE, Out: Register(4), Left: Register(11)},
		{Op: CONVERT_F64_TO_I64, Out: Register(17), Left: Register(15)},
		{Op: COPY_VALUE, Out: Register(5), Left: Register(17)},
		{Op: COPY_VALUE, Out: Register(3), Left: Register(16)},
		// barrf := xyzzy * 10000 + nerrf * 100;
		{Op: NOOP, Out: Register(4)},
		{Op: LOAD_CONST, Out: Register(18), Left: Constant(8)},
		{Op: U64_MULTIPLY, Out: Register(19), Left: Register(4), Right: Register(18)},
		{Op: NOOP, Out: Register(5)},
		{Op: LOAD_CONST, Out: Register(20), Left: Constant(9)},
		{Op: I64_MULTIPLY, Out: Register(21), Left: Register(5), Right: Register(20)},
		{Op: U64_ADD, Out: Register(22), Left: Register(19), Right: Register(21)},
	}
	// FIXME: This example will break with more strict with implict casts
	program = generateBytecode(t, `{
		plugh := 1 - 4;
		xyzzy := 012;
		nerrf := 14;

		xyzzy = 0700;                      # assignment
		plugh = 0.25 * plugh;              # assignment with cast

		# parallel assignment (with and without casts)
		xyzzy, nerrf, plugh = plugh, (xyzzy / 5.0), nerrf;
		barrf := xyzzy * 10000 + nerrf * 100;
	}`)
	assert.Equal(t, expected, program.Instructions)
	assert.Equal(t, constants, program.Constants)

	// Nested block
	constants = []Data{
		0: 0,
		1: FromU64(0600),
		2: FromF64(6.29),
		3: 2,
		4: 3,
	}
	expected = []Instruction{
		{Op: LOAD_CONST, Out: Register(1), Left: Constant(1)},
		{Op: LOAD_CONST, Out: Register(2), Left: Constant(2)},
		{Op: NOOP, Out: Register(2)},
		{Op: LOAD_CONST, Out: Register(3), Left: Constant(3)},
		{Op: CONVERT_I64_TO_F64, Out: Register(4), Left: Register(3)},
		{Op: F64_DIVIDE, Out: Register(5), Left: Register(2), Right: Register(4)},
		{Op: NOOP, Out: Register(5)},
		{Op: NOOP, Out: Register(1)},
		{Op: CONVERT_U64_TO_F64, Out: Register(6), Left: Register(1)},
		{Op: F64_SUBTRACT, Out: Register(7), Left: Register(5), Right: Register(6)},
		{Op: LOAD_CONST, Out: Register(8), Left: Constant(4)},
		{Op: COPY_VALUE, Out: Register(1), Left: Register(8)},
		{Op: NOOP, Out: Register(2)},
	}
	program = generateBytecode(t, `{
		ham  := 0600;
		eggs :: 6.29;

		{
			spam := eggs / 2;
			spam - ham;
			ham = 3;
		}

		eggs;
	}`)
	assert.Equal(t, expected, program.Instructions)
	assert.Equal(t, constants, program.Constants)
}
