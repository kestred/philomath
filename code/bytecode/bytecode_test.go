package bytecode

import (
	"testing"

	// TODO: Maybe don't rely on parser et. al. when more code is stable
	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/parser"
	"github.com/kestred/philomath/code/semantics"
	"github.com/stretchr/testify/assert"
)

func generateBytecode(t *testing.T, input string) *Program {
	p := parser.Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	section := ast.FlattenTree(node, nil)
	semantics.ResolveNames(&section)
	semantics.InferTypes(&section)
	program := NewProgram()
	program.Extend(node)
	return program
}

func TestEncodeArithmetic(t *testing.T) {
	// constants
	constants := map[string][]byte{".LC1": Pack(int64(22))}
	expected := []Instruction{{LOAD, Constant(".LC1", Rg(0, Int64))}}
	program := generateBytecode(t, `22;`)
	assert.Equal(t, constants, program.Data)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	constants = map[string][]byte{".LC1": Pack(uint64(0755))}
	expected = []Instruction{{LOAD, Constant(".LC1", Rg(0, Uint64))}}
	program = generateBytecode(t, `0755;`)
	assert.Equal(t, constants, program.Data)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	constants = map[string][]byte{".LC1": Pack(2.0)}
	expected = []Instruction{{LOAD, Constant(".LC1", Rg(0, Float64))}}
	program = generateBytecode(t, `2.0;`)
	assert.Equal(t, constants, program.Data)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	// add, subtract, multiply, divide
	constants = map[string][]byte{
		".LC1": Pack(int64(2)),
		".LC2": Pack(int64(3)),
		".LC3": Pack(int64(4)),
		".LC4": Pack(int64(5)),
		".LC5": Pack(int64(6)),
	}
	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Int64))},
		{LOAD, Constant(".LC2", Rg(1, Int64))},
		{MULTIPLY, Binary(Rg(0, Int64), Rg(1, Int64), Rg(2, Int64))},
		{LOAD, Constant(".LC3", Rg(3, Int64))},
		{LOAD, Constant(".LC4", Rg(4, Int64))},
		{DIVIDE, Binary(Rg(3, Int64), Rg(4, Int64), Rg(5, Int64))},
		{ADD, Binary(Rg(2, Int64), Rg(5, Int64), Rg(6, Int64))},
		{LOAD, Constant(".LC5", Rg(7, Int64))},
		{SUBTRACT, Binary(Rg(6, Int64), Rg(7, Int64), Rg(8, Int64))},
	}
	program = generateBytecode(t, `2 * 3 + 4 / 5 - 6;`)
	assert.Equal(t, constants, program.Data)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	constants = map[string][]byte{
		".LC1": Pack(2.0),
		".LC2": Pack(4.0),
		".LC3": Pack(8.0),
		".LC4": Pack(16.0),
		".LC5": Pack(32.0),
	}
	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Float64))},
		{LOAD, Constant(".LC2", Rg(1, Float64))},
		{MULTIPLY, Binary(Rg(0, Float64), Rg(1, Float64), Rg(2, Float64))},
		{LOAD, Constant(".LC3", Rg(3, Float64))},
		{LOAD, Constant(".LC4", Rg(4, Float64))},
		{DIVIDE, Binary(Rg(3, Float64), Rg(4, Float64), Rg(5, Float64))},
		{ADD, Binary(Rg(2, Float64), Rg(5, Float64), Rg(6, Float64))},
		{LOAD, Constant(".LC5", Rg(7, Float64))},
		{SUBTRACT, Binary(Rg(6, Float64), Rg(7, Float64), Rg(8, Float64))},
	}
	program = generateBytecode(t, `2.0 * 4.0 + 8.0 / 16.0 - 32.0;`)
	assert.Equal(t, constants, program.Data)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	constants = map[string][]byte{
		".LC1": Pack(int64(2)),
		".LC2": Pack(int64(3)),
		".LC3": Pack(int64(4)),
		".LC4": Pack(int64(5)),
		".LC5": Pack(int64(6)),
	}
	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Uint64))},
		{LOAD, Constant(".LC2", Rg(1, Uint64))},
		{MULTIPLY, Binary(Rg(0, Uint64), Rg(1, Uint64), Rg(2, Uint64))},
		{LOAD, Constant(".LC3", Rg(3, Uint64))},
		{LOAD, Constant(".LC4", Rg(4, Uint64))},
		{DIVIDE, Binary(Rg(3, Uint64), Rg(4, Uint64), Rg(5, Uint64))},
		{ADD, Binary(Rg(2, Uint64), Rg(5, Uint64), Rg(6, Uint64))},
		{LOAD, Constant(".LC5", Rg(7, Uint64))},
		{SUBTRACT, Binary(Rg(6, Uint64), Rg(7, Uint64), Rg(8, Uint64))},
	}
	program = generateBytecode(t, `02 * 03 + 04 / 05 - 06;`)
	assert.Equal(t, constants, program.Data)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	// grouping
	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Int64))},
		{LOAD, Constant(".LC2", Rg(1, Int64))},
		{ADD, Binary(Rg(0, Int64), Rg(1, Int64), Rg(2, Int64))},
		{LOAD, Constant(".LC3", Rg(3, Int64))},
		{MULTIPLY, Binary(Rg(2, Int64), Rg(3, Int64), Rg(4, Int64))},
	}
	program = generateBytecode(t, `(2 + 3) * 4;`)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	// conversions
	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Int64))},
		{LOAD, Constant(".LC2", Rg(1, Int64))},
		{ADD, Binary(Rg(0, Int64), Rg(1, Int64), Rg(2, Int64))},
		{CAST_F64, Unary(Rg(2, Int64), Rg(3, Float64))},
		{LOAD, Constant(".LC3", Rg(4, Float64))},
		{ADD, Binary(Rg(3, Float64), Rg(4, Float64), Rg(5, Float64))},
	}
	program = generateBytecode(t, `(2 + 3) + 4.0;`)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Int64))},
		{CAST_F64, Unary(Rg(0, Int64), Rg(1, Float64))},
		{LOAD, Constant(".LC2", Rg(2, Float64))},
		{ADD, Binary(Rg(1, Float64), Rg(2, Float64), Rg(3, Float64))},
		{LOAD, Constant(".LC3", Rg(4, Int64))},
		{CAST_F64, Unary(Rg(4, Int64), Rg(5, Float64))},
		{ADD, Binary(Rg(3, Float64), Rg(5, Float64), Rg(6, Float64))},
	}
	program = generateBytecode(t, `(2 + 3.0) + 4;`)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Uint64))},
		{LOAD, Constant(".LC2", Rg(1, Uint64))},
		{ADD, Binary(Rg(0, Uint64), Rg(1, Uint64), Rg(2, Uint64))},
		{CAST_F64, Unary(Rg(2, Uint64), Rg(3, Float64))},
		{LOAD, Constant(".LC3", Rg(4, Float64))},
		{ADD, Binary(Rg(3, Float64), Rg(4, Float64), Rg(5, Float64))},
	}
	program = generateBytecode(t, `(02 + 03) + 4.0;`)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Uint64))},
		{CAST_F64, Unary(Rg(0, Uint64), Rg(1, Float64))},
		{LOAD, Constant(".LC2", Rg(2, Float64))},
		{ADD, Binary(Rg(1, Float64), Rg(2, Float64), Rg(3, Float64))},
		{LOAD, Constant(".LC3", Rg(4, Uint64))},
		{CAST_F64, Unary(Rg(4, Uint64), Rg(5, Float64))},
		{ADD, Binary(Rg(3, Float64), Rg(5, Float64), Rg(6, Float64))},
	}
	program = generateBytecode(t, `(02 + 3.0) + 04;`)
	assert.Equal(t, expected, program.Procedures[0].Instructions)
}

func TestEncodeBlock(t *testing.T) {
	// Declarations
	constants := map[string][]byte{
		".LC1": Pack(int64(3)),
		".LC2": Pack(int64(2)),
		".LC3": Pack(0.5),
	}
	expected := []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Int64))},
		{LOAD, Constant(".LC2", Rg(1, Int64))},
		{ADD, Binary(Rg(0, Int64), Rg(1, Int64), Rg(2, Int64))},
		{LOAD, Constant(".LC3", Rg(3, Float64))},
		{CAST_F64, Unary(Rg(0, Int64), Rg(4, Float64))},
		{MULTIPLY, Binary(Rg(3, Float64), Rg(4, Float64), Rg(5, Float64))},
		{CAST_F64, Unary(Rg(0, Int64), Rg(6, Float64))},
		{DIVIDE, Binary(Rg(5, Float64), Rg(6, Float64), Rg(7, Float64))},
	}
	program := generateBytecode(t, `{
		hoge :: 3;          // constant decl
		hoge + 2;           // one ident in expr, result ignored

		piyo := 0.5 * hoge; // mutable decl
		piyo / hoge;        // two ident in expr; TODO: return statement
	}`)
	assert.Equal(t, constants, program.Data)
	assert.Equal(t, expected, program.Procedures[0].Instructions)
	t.Log(Pretty(program.Procedures[0].Instructions))

	// Simple and Parallel Assignment
	constants = map[string][]byte{
		".LC1": Pack(int64(1)),
		".LC2": Pack(int64(4)),
		".LC3": Pack(uint64(012)),
		".LC4": Pack(int64(14)),
		".LC5": Pack(uint64(0700)),
		".LC6": Pack(0.25),
		".LC7": Pack(5.0),
		".LC8": Pack(int64(10000)),
		".LC9": Pack(int64(100)),
	}
	expected = []Instruction{
		// plugh := 1 - 3;
		{LOAD, Constant(".LC1", Rg(0, Int64))},
		{LOAD, Constant(".LC2", Rg(1, Int64))},
		{SUBTRACT, Binary(Rg(0, Int64), Rg(1, Int64), Rg(2, Int64))},
		// xyzzy := 012;
		{LOAD, Constant(".LC3", Rg(3, Uint64))},
		// nerrf := 14;
		{LOAD, Constant(".LC4", Rg(4, Int64))},
		// xyzzy = 0700;
		{LOAD, Constant(".LC5", Rg(5, Uint64))},
		{COPY, Unary(Rg(5, Uint64), Rg(3, Uint64))},
		// plugh = 0.25 * plugh;
		{LOAD, Constant(".LC6", Rg(6, Float64))},
		{CAST_F64, Unary(Rg(2, Int64), Rg(7, Float64))},
		{MULTIPLY, Binary(Rg(6, Float64), Rg(7, Float64), Rg(8, Float64))},
		{CAST_I64, Unary(Rg(8, Float64), Rg(9, Int64))},
		{COPY, Unary(Rg(9, Int64), Rg(2, Int64))},
		// xyzzy, nerrf, plugh = plugh, (xyzzy / 5.0), nerrf;
		{COPY, Unary(Rg(2, Int64), Rg(10, Int64))},
		{CAST_F64, Unary(Rg(3, Uint64), Rg(11, Float64))},
		{LOAD, Constant(".LC7", Rg(12, Float64))},
		{DIVIDE, Binary(Rg(11, Float64), Rg(12, Float64), Rg(13, Float64))},
		{COPY, Unary(Rg(13, Float64), Rg(14, Float64))},
		{COPY, Unary(Rg(4, Int64), Rg(15, Int64))},
		{CAST_U64, Unary(Rg(10, Int64), Rg(16, Uint64))},
		{COPY, Unary(Rg(16, Uint64), Rg(3, Uint64))},
		{CAST_I64, Unary(Rg(14, Float64), Rg(17, Int64))},
		{COPY, Unary(Rg(17, Int64), Rg(4, Int64))},
		{COPY, Unary(Rg(15, Int64), Rg(2, Int64))},
		// barrf := xyzzy * 10000 + nerrf * 100;
		{LOAD, Constant(".LC8", Rg(18, Int64))},
		{CAST_U64, Unary(Rg(18, Int64), Rg(19, Uint64))},
		{MULTIPLY, Binary(Rg(3, Uint64), Rg(19, Uint64), Rg(20, Uint64))},
		{LOAD, Constant(".LC9", Rg(21, Int64))},
		{MULTIPLY, Binary(Rg(4, Int64), Rg(21, Int64), Rg(22, Int64))},
		{CAST_U64, Unary(Rg(22, Int64), Rg(23, Uint64))},
		{ADD, Binary(Rg(20, Uint64), Rg(23, Uint64), Rg(24, Uint64))},
	}
	// FIXME: This example will break with more strict with implict casts
	program = generateBytecode(t, `{
		plugh := 1 - 4;
		xyzzy := 012;
		nerrf := 14;

		xyzzy = 0700;                      // assignment
		plugh = 0.25 * plugh;              // assignment with cast

		// parallel assignment (with and without casts)
		xyzzy, nerrf, plugh = plugh, (xyzzy / 5.0), nerrf;
		barrf := xyzzy * 10000 + nerrf * 100;
	}`)
	assert.Equal(t, constants, program.Data)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	// Nested block
	constants = map[string][]byte{
		".LC1": Pack(uint64(0600)),
		".LC2": Pack(6.29),
		".LC3": Pack(int64(2)),
		".LC4": Pack(int64(3)),
	}
	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Uint64))},
		{LOAD, Constant(".LC2", Rg(1, Float64))},
		{LOAD, Constant(".LC3", Rg(2, Int64))},
		{CAST_F64, Unary(Rg(2, Int64), Rg(3, Float64))},
		{DIVIDE, Binary(Rg(1, Float64), Rg(3, Float64), Rg(4, Float64))},
		{CAST_F64, Unary(Rg(0, Uint64), Rg(5, Float64))},
		{SUBTRACT, Binary(Rg(4, Float64), Rg(5, Float64), Rg(6, Float64))},
		{LOAD, Constant(".LC4", Rg(7, Int64))},
		{CAST_U64, Unary(Rg(7, Int64), Rg(8, Uint64))},
		{COPY, Unary(Rg(8, Uint64), Rg(0, Uint64))},
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
	assert.Equal(t, constants, program.Data)
	assert.Equal(t, expected, program.Procedures[0].Instructions)

	// Inline assembly
	constants = map[string][]byte{".LC1": Pack(int64(3)), ".LC2": Pack(int64(0))}
	expected = []Instruction{
		{LOAD, Constant(".LC1", Rg(0, Int64))},
		{LOAD, Constant(".LC2", Rg(1, Uint64))},
		{CALL_ASM, nil},
	}
	program = generateBytecode(t, `{
	  input := 3;
		output := 0;

	  #asm { mov output, input }

		output;
	}`)
	assert.Equal(t, constants, program.Data)
	insts := program.Procedures[0].Instructions
	for i, inst := range insts {
		assert.Equal(t, expected[i].Op, inst.Op)
	}
	if assert.Equal(t, 3, len(insts)) {
		asm := insts[2].Args.(*AssemblyArgs)
		assert.Equal(t, " mov output, input ", asm.Source)
		assert.Equal(t, []Register{Rg(0, Int64)}, asm.InputRegisters)
		assert.Equal(t, Rg(1, Uint64), asm.OutputRegister)
		assert.True(t, asm.HasOutput)
	}
}
