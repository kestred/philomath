package interpreter

import (
	"os/exec"
	"unsafe"
)

// LoadAsm compiles an assembly block, wrapping it in a function where
// variables are replaced with arguments and return values according to the
// Golang calling conventions
func LoadAsm() {
  /* TODO:

  Parse the ASM to find any variables

  Replace right-hand variables with arguments (eg. [%rsp+8], [%rsp+10])

  Replace left-hand variables with return values (eg. [%rsp+18])

  Format the assembly in a template

    .intel_syntax
    caller.asmN:

      <inline assembly>

      ret

  Using golang.org/pkg/os/exec, run the GNU Assembler

  Mmap a Read/Write memory block

  Read the assembled file into the block

  Mprotect the block with Read/Exec

  Find the "caller.asmN" function in the ELF object file

  Using golang.org/pkg/unsafe cast the pointer to the appropriate prototype

  Call the function as needed in Evaluate

  */
}

type OneToZero *func(a uintptr)
type OneToOne *func(a uintptr) uintptr
type TwoToZero *func(a, b uintptr)
type TwoToOne *func(a, b uintptr) uintptr
type ThreeToZero *func(a, b, c uintptr)
type ThreeToOne *func(a, b, c uintptr) uintptr
type FourToZero *func(a, b, c, d uintptr)
type FourToOne *func(a, b, c, d uintptr) uintptr
