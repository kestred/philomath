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

typedef void(*Fn0x0)();
void Call0x0(void* fn) { ((Fn0x0)fn)(); }

typedef unsigned long long u64;
typedef void(*Int1x0)(u64);
typedef void(*Int2x0)(u64, u64);
typedef void(*Int3x0)(u64, u64, u64);
typedef void(*Int4x0)(u64, u64, u64, u64);
typedef void(*Int5x0)(u64, u64, u64, u64, u64);
typedef void(*Int6x0)(u64, u64, u64, u64, u64, u64);
void CallInt1x0(void* fn, u64 x) { ((Int1x0)fn)(x); }
void CallInt2x0(void* fn, u64 x, u64 y) { ((Int2x0)fn)(x, y); }
void CallInt3x0(void* fn, u64 x, u64 y, u64 z) { ((Int3x0)fn)(x, y, z); }
void CallInt4x0(void* fn, u64 x, u64 y, u64 z, u64 w) { ((Int4x0)fn)(x, y, z, w); }
void CallInt5x0(void* fn, u64 x, u64 y, u64 z, u64 w, u64 v) { ((Int5x0)fn)(x, y, z, w, v); }
void CallInt6x0(void* fn, u64 x, u64 y, u64 z, u64 w, u64 v, u64 u) { ((Int6x0)fn)(x, y, z, w, v, u); }
typedef u64(*Int0x1)();
typedef u64(*Int1x1)(u64);
typedef u64(*Int2x1)(u64, u64);
typedef u64(*Int3x1)(u64, u64, u64);
typedef u64(*Int4x1)(u64, u64, u64, u64);
typedef u64(*Int5x1)(u64, u64, u64, u64, u64);
typedef u64(*Int6x1)(u64, u64, u64, u64, u64, u64);
u64 CallInt0x1(void* fn) { return ((Int0x1)fn)(); }
u64 CallInt1x1(void* fn, u64 x) { return ((Int1x1)fn)(x); }
u64 CallInt2x1(void* fn, u64 x, u64 y) { return ((Int2x1)fn)(x, y); }
u64 CallInt3x1(void* fn, u64 x, u64 y, u64 z) { return ((Int3x1)fn)(x, y, z); }
u64 CallInt4x1(void* fn, u64 x, u64 y, u64 z, u64 w) { return ((Int4x1)fn)(x, y, z, w); }
u64 CallInt5x1(void* fn, u64 x, u64 y, u64 z, u64 w, u64 v) { return ((Int5x1)fn)(x, y, z, w, v); }
u64 CallInt6x1(void* fn, u64 x, u64 y, u64 z, u64 w, u64 v, u64 u) { return ((Int6x1)fn)(x, y, z, w, v, u); }

typedef double f64;
typedef void(*Float1x0)(f64);
typedef void(*Float2x0)(f64, f64);
typedef void(*Float3x0)(f64, f64, f64);
typedef void(*Float4x0)(f64, f64, f64, f64);
typedef void(*Float5x0)(f64, f64, f64, f64, f64);
typedef void(*Float6x0)(f64, f64, f64, f64, f64, f64);
void CallFloat1x0(void* fn, f64 x) { ((Float1x0)fn)(x); }
void CallFloat2x0(void* fn, f64 x, f64 y) { ((Float2x0)fn)(x, y); }
void CallFloat3x0(void* fn, f64 x, f64 y, f64 z) { ((Float3x0)fn)(x, y, z); }
void CallFloat4x0(void* fn, f64 x, f64 y, f64 z, f64 w) { ((Float4x0)fn)(x, y, z, w); }
void CallFloat5x0(void* fn, f64 x, f64 y, f64 z, f64 w, f64 v) { ((Float5x0)fn)(x, y, z, w, v); }
void CallFloat6x0(void* fn, f64 x, f64 y, f64 z, f64 w, f64 v, f64 u) { ((Float6x0)fn)(x, y, z, w, v, u); }
typedef f64(*Float0x1)();
typedef f64(*Float1x1)(f64);
typedef f64(*Float2x1)(f64, f64);
typedef f64(*Float3x1)(f64, f64, f64);
typedef f64(*Float4x1)(f64, f64, f64, f64);
typedef f64(*Float5x1)(f64, f64, f64, f64, f64);
typedef f64(*Float6x1)(f64, f64, f64, f64, f64, f64);
f64 CallFloat0x1(void* fn) { return ((Float0x1)fn)(); }
f64 CallFloat1x1(void* fn, f64 x) { return ((Float1x1)fn)(x); }
f64 CallFloat2x1(void* fn, f64 x, f64 y) { return ((Float2x1)fn)(x, y); }
f64 CallFloat3x1(void* fn, f64 x, f64 y, f64 z) { return ((Float3x1)fn)(x, y, z); }
f64 CallFloat4x1(void* fn, f64 x, f64 y, f64 z, f64 w) { return ((Float4x1)fn)(x, y, z, w); }
f64 CallFloat5x1(void* fn, f64 x, f64 y, f64 z, f64 w, f64 v) { return ((Float5x1)fn)(x, y, z, w, v); }
f64 CallFloat6x1(void* fn, f64 x, f64 y, f64 z, f64 w, f64 v, f64 u) { return ((Float6x1)fn)(x, y, z, w, v, u); }
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	c "unsafe"

	// "github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/bytecode"
	"github.com/kestred/philomath/code/utils"
)

type Function c.Pointer

func Assemble(asm *bytecode.AssemblyArgs) {
	label, source := generateAssembly(asm)
	tmpdir, err := ioutil.TempDir("", "phi-")
	utils.Assert(err == nil, "unable to create tmpdir for building inline assembly")
	defer os.RemoveAll(tmpdir)

	objpath := tmpdir + "/" + label + ".o"
	cmd := exec.Command("/usr/bin/as", "-o", objpath)
	cmd.Stdin = strings.NewReader(source)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	utils.Assert(err == nil, `/usr/bin/as: %s`, err)

	libpath := tmpdir + "/" + label + ".so"
	cmd = exec.Command("/usr/bin/ld", "-o", libpath, objpath, "-shared", "--export-dynamic")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	utils.Assert(err == nil, `/usr/bin/ld: %s`, err)

	loaded := C.LoadCode(C.CString(libpath), C.CString(label))
	if loaded.Err != nil {
		panic(C.GoString(loaded.Err)) //return nil, err
	}

	asm.Wrapper = loaded.Fn
}

var nextLabel uint

// TODO: more intelligent asm generation
func generateAssembly(asm *bytecode.AssemblyArgs) (label string, source string) {
	var offset int
	var parts []string
	for i, binding := range asm.InputBindings {
		// insert the output binding if it appears before this input
		if asm.HasOutput && asm.OutputBinding.Offset < binding.Offset {
			binding := asm.OutputBinding
			register := "%rax"

			parts = append(parts, asm.Source[offset:binding.Offset], register)
			offset = binding.Offset + len(binding.Name.Literal)
		}

		var register string
		switch i {
		case 0:
			register = "%rdi"
		case 1:
			register = "%rsi"
		case 2:
			register = "%rdx"
		case 3:
			register = "%rcx"
		case 4:
			register = "%r8"
		case 5:
			register = "%r9"
		}

		parts = append(parts, asm.Source[offset:binding.Offset], register)
		offset = binding.Offset + len(binding.Name.Literal)
	}

	// insert the output binding if it has not yet been inserted
	if asm.HasOutput && asm.OutputBinding.Offset >= offset {
		binding := asm.OutputBinding
		register := "%rax"

		parts = append(parts, asm.Source[offset:binding.Offset], register)
		offset = binding.Offset + len(binding.Name.Literal)
	}

	// combine the parts into an updated source
	parts = append(parts, asm.Source[offset:len(asm.Source)])
	source = strings.Join(parts, "")

	nextLabel += 1
	label = fmt.Sprintf("interpreter.wrapper%d", nextLabel)
	source = formatAssembly(label, source)
	return
}

func u64Pack(i C.u64) []byte {
	return bytecode.Pack(uint64(i))
}

func u64Unpack(b []byte) C.u64 {
	var ui uint64
	err := bytecode.Unpack(b, &ui)
	utils.Assert(err == nil, `%v (at inline assembly)`, err)
	return C.u64(ui)
}

func CallAsm(asm *bytecode.AssemblyArgs, registers [][]byte) {
	if asm.Wrapper == nil {
		Assemble(asm)
	}

	params := make([][]byte, len(asm.InputRegisters))
	for i, register := range asm.InputRegisters {
		params[i] = registers[register.Loc]
	}

	if asm.HasOutput {
		// if ast.IsFloat(asm.OutputBinding.Name.Type) {
		returns := registers[asm.OutputRegister.Loc : asm.OutputRegister.Loc+1]
		CallAsmInt(Function(asm.Wrapper), params, returns)
		// }
	} else if len(asm.InputBindings) > 0 {
		// if ast.IsFloat(asm.InputBindings[0].Name.Type) {
		CallAsmInt(Function(asm.Wrapper), params, nil)
		// }
	} else {
		C.Call0x0(asm.Wrapper)
	}
}

// FIXME: outs shouldn't be an array anymore
func CallAsmInt(fn Function, ins [][]byte, outs [][]byte) {
	index := callIndex(len(ins), len(outs))
	switch index {
	case index1x0:
		C.CallInt1x0(c.Pointer(fn), u64Unpack(ins[0]))
	case index2x0:
		C.CallInt2x0(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1]))
	case index3x0:
		C.CallInt3x0(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1]), u64Unpack(ins[2]))
	case index4x0:
		C.CallInt4x0(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1]), u64Unpack(ins[2]), u64Unpack(ins[3]))
	case index5x0:
		C.CallInt4x0(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1]), u64Unpack(ins[2]), u64Unpack(ins[3]))
	case index6x0:
		C.CallInt4x0(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1]), u64Unpack(ins[2]), u64Unpack(ins[3]))
	case index0x1:
		outs[0] = u64Pack(C.CallInt0x1(c.Pointer(fn)))
	case index1x1:
		outs[0] = u64Pack(C.CallInt1x1(c.Pointer(fn), u64Unpack(ins[0])))
	case index2x1:
		outs[0] = u64Pack(C.CallInt2x1(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1])))
	case index3x1:
		outs[0] = u64Pack(C.CallInt3x1(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1]), u64Unpack(ins[2])))
	case index4x1:
		outs[0] = u64Pack(C.CallInt4x1(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1]), u64Unpack(ins[2]), u64Unpack(ins[3])))
	case index5x1:
		outs[0] = u64Pack(C.CallInt5x1(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1]), u64Unpack(ins[2]), u64Unpack(ins[3]), u64Unpack(ins[4])))
	case index6x1:
		outs[0] = u64Pack(C.CallInt6x1(c.Pointer(fn), u64Unpack(ins[0]), u64Unpack(ins[1]), u64Unpack(ins[2]), u64Unpack(ins[3]), u64Unpack(ins[4]), u64Unpack(ins[5])))
	default:
		utils.InvalidCodePath()
	}
}

func callIndex(ins, outs int) uint8 {
	return uint8((ins << 1) | outs)
}

var (
	index1x0 = callIndex(1, 0)
	index2x0 = callIndex(2, 0)
	index3x0 = callIndex(3, 0)
	index4x0 = callIndex(4, 0)
	index5x0 = callIndex(5, 0)
	index6x0 = callIndex(6, 0)
	index0x1 = callIndex(1, 1)
	index1x1 = callIndex(1, 1)
	index2x1 = callIndex(2, 1)
	index3x1 = callIndex(3, 1)
	index4x1 = callIndex(4, 1)
	index5x1 = callIndex(5, 1)
	index6x1 = callIndex(6, 1)
)
