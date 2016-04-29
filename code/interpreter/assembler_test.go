package interpreter

import (
	"testing"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/bytecode"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	asm := &bytecode.AssemblyArgs{
		Source:        "  mov output, input",
		HasOutput:     true,
		OutputBinding: ast.AsmBinding{ast.Ident("output"), 6},
		InputBindings: []ast.AsmBinding{{ast.Ident("input"), 14}},
	}

	label, source := generateAssembly(asm)
	assert.Equal(t, "interpreter.wrapper1", label)
	assert.Equal(t, `
.intel_syntax
.global interpreter.wrapper1
.section .text

interpreter.wrapper1:
  mov %rax, %rdi
  ret
`, source)
}
