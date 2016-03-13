package parser

import (
	"testing"

	"github.com/kestred/philomath/code/ast"
	"github.com/stretchr/testify/assert"
)

func parseAsm(t *testing.T, input string) *ast.AsmBlock {
	asm := ast.Asm(input)
	parseAssembly(asm)
	return asm
}

func TestParseAsm(t *testing.T) {
	var asm *ast.AsmBlock

	// simple syscall example
	asm = parseAsm(t, `
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
	asm = parseAsm(t, `mov [retval], %rax`)
	assert.Equal(t, []ast.AsmBinding{{ast.Ident("retval"), 5}}, asm.Inputs)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Outputs)

	// non-mov instruction
	asm = parseAsm(t, `push retval`)
	assert.Equal(t, []ast.AsmBinding{{ast.Ident("retval"), 5}}, asm.Inputs)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Outputs)

	// other-mov instruction
	asm = parseAsm(t, `movnti retval, %rax`)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Inputs)
	assert.Equal(t, []ast.AsmBinding{{ast.Ident("retval"), 7}}, asm.Outputs)

	// ignore size specifiers
	asm = parseAsm(t, `
		mov %eax, dword ptr [%ecx]
		mov %eax, DWORD PTR [%ecx]
	`)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Inputs)
	assert.Equal(t, []ast.AsmBinding(nil), asm.Outputs)

	// ignore labels
	asm = parseAsm(t, `
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
