package interpreter

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>

struct Loaded {
  void* Fn;
  char* Err;
};

struct Loaded
LoadCode(const char *file, const char* function) {
  char* err;

  void* lib = dlopen(file, RTLD_LAZY|RTLD_LOCAL);
  err = dlerror();
  if (err != 0) {
    struct Loaded result;
    result.Err = err;
    return result;
  }

  void* fn = dlsym(lib, function);
  err = dlerror();

  struct Loaded result;
  result.Fn = fn;
  result.Err = err;
  return result;
}

typedef unsigned long long u64;
typedef u64(*Fn4To1)(u64, u64, u64, u64);
u64 Call4To1(void* fn, u64 x, u64 y, u64 z, u64 w) {
  return ((Fn4To1)fn)(x, y, z, w);
}
*/
import "C"
import (
  "fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	unix "syscall"
  c "unsafe"
)

func init() {
	examplefn := "func1"
	exampleasm := `
.intel_syntax

.global func1
.section .text

func1:
  mov %rax, %rdi
  mov %rdi, %rsi
  mov %rsi, %rdx
  mov %rdx, %rcx
  syscall
  ret
`

	tmpdir, err := ioutil.TempDir("", "phi-")
	if err != nil {
		panic(err) //return nil, err
	}
  // TODO: Clean this up later, at not at function exit
	defer os.RemoveAll(tmpdir)

	objpath := tmpdir + "/" + examplefn + ".o"
	cmd := exec.Command("/usr/bin/as", "-o", objpath)
	cmd.Stdin = strings.NewReader(exampleasm)
  // TODO: Collect stderr and return it if err is not nil
  // cmd.Stdout = os.Stdout
  // cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err) //return nil, err
	}

	libpath := tmpdir + "/" + examplefn + ".so"
	cmd = exec.Command("/usr/bin/ld", "-o", libpath, objpath, "-shared", "--export-dynamic")
  // TODO: Collect stderr and return it if err is not nil
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err) //return nil, err
	}

  loaded := C.LoadCode(C.CString(libpath), C.CString(examplefn))
  fmt.Println(libpath, loaded.Fn, C.GoString(loaded.Err))

  if loaded.Err == nil {
    message := []byte("Hello world!\n")
    result := int(C.Call4To1(
      loaded.Fn,
      unix.SYS_WRITE,
      C.u64(unix.Stdout),
      C.u64(uintptr(c.Pointer(&message[0]))),
      C.u64(len(message)),
    ))

    if result >= 0 {
      fmt.Println("Bytes written:", result)
    } else {
      fmt.Println("Errno:", -result)
    }
  }
}

// NOTE: Below is my attempt to compile inline assembly and dynamically
// load it from pure Go code without resorting to cgo and linking to C-source.
//
// I've left it here in case someone smarter/more-persistent is able to figure
// out what I've done wrong.
//
// It seems like Go's default linker uses the stdcall calling convention
// even when compiling for x86 64-bit, but I couldn't confirm it.

/*

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	unix "syscall"
	c "unsafe"
)

*/

/*

// -T-O-D-O-: OSX, Windows, and OpenBSD support

// LoadAsm compiles an assembly block, wrapping it in a function where
// variables are replaced with arguments and return values according to the
// Golang calling conventions
func LoadAsm() (Function, error) {
	// -T-O-D-O-:
	//
	// Parse the ASM to find any variables
	//
	// Replace constants with their constant value
	//
	// Replace right-hand variables with arguments (eg. [%rsp+8], [%rsp+10])
	//
	// Replace left-hand variables with return values (eg. [%rsp+18])
	//
	// Format the assembly in a template
	//
	//   .intel_syntax
	//   caller.asmN:
	//
	//     <inline assembly>
	//
	//     ret


	examplefn := "main.asm1"
	exampleasm := `
.intel_syntax

main.asm1:
  mov %rax, [%rsp+0x08]
  mov %rdi, [%rsp+0x10]
  mov %rsi, [%rsp+0x18]
  mov %rdx, [%rsp+0x20]
  syscall
  mov [%rsp+0x28], %rax
  ret
`

	tmpdir, err := ioutil.TempDir("", "phi-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpdir)

	outpath := tmpdir + "/" + examplefn + ".o"
	cmd := exec.Command("/usr/bin/as", "-o", outpath)
	cmd.Stdin = strings.NewReader(exampleasm)
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(outpath)
	if err != nil {
		return nil, err
	}
	size := int(info.Size())

	// -T-O-D-O-: Pack multiple small functions into a single page
	write_prot := unix.PROT_READ | unix.PROT_WRITE
	anon_flag := unix.MAP_PRIVATE | unix.MAP_ANONYMOUS
	mem, err := unix.Mmap(0, 0, size, write_prot, anon_flag)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(outpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	read, err := file.Read(mem)
	if err != nil {
		return nil, err
	} else if read != size {
		if read < size {
			return nil, errors.New("object file size was smaller than expected")
		} else {
			return nil, errors.New("object file size was larger than expected")
		}
	}

	exec_prot := unix.PROT_READ | unix.PROT_EXEC
	err = unix.Mprotect(mem, exec_prot)
	if err != nil {
		return nil, err
	}

	if len(mem) < 0x40 {
		return nil, errors.New("invalid ELF object file; too short")
	} else if mem[0] != 0x7F || mem[1] != 'E' || mem[2] != 'L' || mem[3] != 'F' {
		return nil, errors.New("invalid ELF object file; bad masgic")
	} else if mem[0x4] != 2 {
		return nil, errors.New("unsupported ELF object file; must be 64-bit")
	} else if mem[0x5] != 1 {
		return nil, errors.New("unsupported ELF object file; must be little-endian")
	} else if mem[0x6] != 1 {
		return nil, errors.New("unsupported ELF object file; must be version 1")
	} else if mem[0x10] != 1 {
		return nil, errors.New("unsupported ELF object file; must be relocatable")
	} // -T-O-D-O-: Check instruction set matches local machine

	// phoff := *(*uint64)(c.Pointer(&mem[0x20]))
	// phsize := *(*uint16)(c.Pointer(&mem[0x36]))
	phnum := *(*uint16)(c.Pointer(&mem[0x38]))
	// phend := phoff + uint64(phsize*phnum)
	if phnum != 0 {
		return nil, errors.New("unsupported ELF object file; not currently expecting program headers")
	}

	ehsize := *(*uint16)(c.Pointer(&mem[0x34]))

	// -T-O-D-O-: Find the function locations "correctly"
	fnptr := c.Pointer(&mem[ehsize])
	fnwrap := &FuncWrapper{fnptr}
	// -T-O-D-O-: Cast to correct type based on parsed assembly
	fn := FourToOne(c.Pointer(&fnwrap.Pointer))
	return fn, nil
}

type FuncWrapper struct {
	Pointer c.Pointer
}

type Function interface{}
type OneToZero *func(x uintptr)
type OneToOne *func(x uintptr) uintptr
type TwoToZero *func(x, y uintptr)
type TwoToOne *func(x, y uintptr) uintptr
type ThreeToZero *func(x, y, z uintptr)
type ThreeToOne *func(x, y, z uintptr) uintptr
type FourToZero *func(x, y, z, w uintptr)
type FourToOne *func(x, y, z, w uintptr) uintptr

*/

/*  File: main.go
package main

import code "github.com/kestred/philomath/code/interpreter"
import unix "syscall"
import c "unsafe"

var msg = []byte("Hello world!\n")

func main() {
	fn, err := code.LoadAsm()
	if err != nil {
		return
	}

  call := fn.(code.FourToOne)
	(*call)(
    unix.SYS_WRITE,
    uintptr(unix.Stdout),
    uintptr(c.Pointer(&msg[0])),
    uintptr(len(msg)),
  )
}

*/
