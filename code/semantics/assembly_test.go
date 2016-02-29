package semantics

import (
	"testing"

	"github.com/kestred/philomath/code/ast"
	"github.com/stretchr/testify/assert"
)

func processAsm(t *testing.T, input string) *ast.AsmBlock {
	asm := ast.Asm(input)
	section := ast.FlattenTree(asm, nil)
	PreprocessAssembly(&section)
	return asm
}

func TestPreprocessAsm(t *testing.T) {
	var asm *ast.AsmBlock

	// simple syscall example
	asm = processAsm(t, `
		mov     %rax, unix_write
		mov     %rdi, unix_stdout
		mov     %rsi, message
		mov     %rdx, 13
		syscall
		mov     retval, %rax
	`)
	assert.Equal(t, []ast.AsmBinding{
		{ast.Ident("unix_write"), 17},
		{ast.Ident("unix_stdout"), 44},
		{ast.Ident("message"), 72},
	}, asm.Inputs)
	assert.Equal(t, []ast.AsmBinding{
		{ast.Ident("retval"), 119},
	}, asm.Outputs)

	// memory operands
	asm = processAsm(t, `mov [retval], %rax`)
	assert.Equal(t, []ast.AsmBinding{{ast.Ident("retval"), 5}}, asm.Inputs)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Outputs)

	// non-mov instruction
	asm = processAsm(t, `push retval`)
	assert.Equal(t, []ast.AsmBinding{{ast.Ident("retval"), 5}}, asm.Inputs)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Outputs)

	// other-mov instruction
	asm = processAsm(t, `movnti retval, %rax`)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Inputs)
	assert.Equal(t, []ast.AsmBinding{{ast.Ident("retval"), 7}}, asm.Outputs)

	// ignore size specifiers
	asm = processAsm(t, `
		mov %eax, dword ptr [%ecx]
		mov %eax, DWORD PTR [%ecx]
	`)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Inputs)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Outputs)

	// ignore labels
	asm = processAsm(t, `
		jump someplace
someplace:
		mov  %ebx, 1
		mov  %eax, input
		cmp  %eax, %ebx
		je someplace
	`)
	assert.Equal(t, []ast.AsmBinding{{ast.Ident("input"), 57}}, asm.Inputs)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Outputs)
}
