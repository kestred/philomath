package bytecode

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/utils"
)

type Opcode int16
type Type int16
type Location int16
type Register struct {
	Loc Location
	Typ Type
}

func Rg(loc Location, typ Type) Register {
	return Register{Loc: loc, Typ: typ}
}

const (
	None Type = iota
	Uint8
	Uint16
	Uint32
	Uint64
	// Uint128
	// Uint256
	Int8
	Int16
	Int32
	Int64
	// Int128
	// Int256
	Float32
	Float64
	// Float128
	// Float256
	Pointer
)

func (t Type) String() string {
	switch t {
	case None:
		return "None"
	case Uint8:
		return "Uint8"
	case Uint16:
		return "Uint16"
	case Uint32:
		return "Uint32"
	case Uint64:
		return "Uint64"
	case Int8:
		return "Int8"
	case Int16:
		return "Int16"
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case Float32:
		return "Float32"
	case Float64:
		return "Float64"
	case Pointer:
		return "Pointer"
	default:
		return fmt.Sprintf("Type(%d)", t)
	}
}

const OutOfRegisters = 32767
const (
	NOOP Opcode = iota

	PUSH  // make room on the stack
	COPY  // move from register to register
	LOAD  // move from pointer to register
	STORE // move from register to pointer

	CALL
	CALL_ASM
	RETURN

	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE

	CAST_I64
	CAST_U64
	CAST_F64
)

var opcodes = [...]string{
	NOOP: "No operation",

	PUSH:  "Push",
	COPY:  "Copy",
	LOAD:  "Load",
	STORE: "Store",

	CALL:     "Call procedure",
	CALL_ASM: "Call assembly",
	RETURN:   "Return",

	ADD:      "Addition",
	SUBTRACT: "Subtraction",
	MULTIPLY: "Multiplication",
	DIVIDE:   "Division",

	CAST_I64: "Cast to signed",
	CAST_U64: "Cast to unsigned",
	CAST_F64: "Cast to float",
}

func (op Opcode) String() string {
	s := ""
	if 0 <= op && op < Opcode(len(opcodes)) {
		s = opcodes[op]
	}
	if s == "" {
		s = "Opcode(" + strconv.Itoa(int(op)) + ")"
	}
	return s
}

func Pack(v interface{}) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, v) // FIXME: Cross-compilation endianness
	return buf.Bytes()
}

func Unpack(b []byte, v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("can't unpack register into non-pointer type")
	}

	return binary.Read(bytes.NewReader(b), binary.LittleEndian, v) // FIXME: Cross-compilation endianness
}

type Instruction struct {
	Op   Opcode
	Args interface{}
}

func Inst(op Opcode, args interface{}) Instruction {
	return Instruction{Op: op, Args: args}
}

func Pretty(insts []Instruction) string {
	var buf bytes.Buffer
	for _, inst := range insts {
		fmt.Fprintln(&buf, inst)
	}
	return buf.String()
}

type NullaryArgs struct {
	Rg Register
}

func Nullary(rg Register) NullaryArgs {
	return NullaryArgs{Rg: rg}
}

type UnaryArgs struct {
	In  Register
	Out Register
}

func Unary(in, out Register) UnaryArgs {
	return UnaryArgs{In: in, Out: out}
}

type BinaryArgs struct {
	Left  Register
	Right Register
	Out   Register
}

func Binary(left, right, out Register) BinaryArgs {
	return BinaryArgs{Left: left, Right: right, Out: out}
}

type ConstantArgs struct {
	Name string
	Out  Register
	Ptr  bool
}

func Constant(name string, out Register) ConstantArgs {
	return ConstantArgs{Name: name, Out: out}
}

func ConstPtr(name string, out Register) ConstantArgs {
	return ConstantArgs{Name: name, Out: out, Ptr: true}
}

type AssemblyArgs struct {
	Source  string
	Wrapper unsafe.Pointer // for interpreter
	Tempdir string

	HasOutput      bool
	OutputBinding  ast.AsmBinding
	OutputRegister Register
	InputBindings  []ast.AsmBinding
	InputRegisters []Register
}

type ProcedureArgs struct {
	Proc Procedure
	Out  Register
	In   []Register
}

func Proc(proc Procedure, out Register, in []Register) ProcedureArgs {
	return ProcedureArgs{Proc: proc, Out: out, In: in}
}

type Program struct {
	Bss        map[string]int // map string to reserved size
	Data       map[string][]byte
	Text       map[string]int // map to Procedure index
	Procedures []Procedure

	nextConstantId int
}

func NewProgram() *Program {
	prog := &Program{}
	prog.Bss = map[string]int{}
	prog.Data = map[string][]byte{}
	prog.Text = map[string]int{"start_": 0}
	prog.NewProcedure() // start_
	return prog
}

func (p *Program) Extend(node ast.Node) {
	p.Procedures[0].Extend(node)
}

func (p *Program) DefineBss(name string, size int) {
	p.Bss[name] = size
}

func (p *Program) DefineData(name string, data []byte) {
	p.Data[name] = data
}

func (p *Program) NextConstantName() string {
	p.nextConstantId += 1
	return ".LC" + strconv.Itoa(p.nextConstantId)
}

func (p *Program) NewProcedure() *Procedure {
	next := len(p.Procedures)
	p.Procedures = append(p.Procedures, Procedure{
		Index:        next,
		Program:      p,
		Instructions: []Instruction{},
		Registers:    map[ast.Decl]Register{},
		PrevResult:   Rg(-1, None),
		NextFree:     +0,
	})
	return &p.Procedures[next]
}

type Procedure struct {
	Index        int
	Program      *Program
	Instructions []Instruction

	// state for bytecode generation
	Registers  map[ast.Decl]Register
	PrevResult Register // last result
	NextFree   Location // next assignable location

	Arguments []Register
}

func (p *Procedure) Extend(node ast.Node) {
	endRegister := Rg(0, None)
	switch n := node.(type) {
	case *ast.TopScope:
		for _, decl := range n.Decls {
			p.Extend(decl)
		}

	case *ast.Block:
		for _, subnode := range n.Nodes {
			p.Extend(subnode)
		}

	case *ast.AsmBlock:
		asm := new(AssemblyArgs)
		asm.Source = n.Source
		asm.InputBindings = n.Inputs
		asm.InputRegisters = make([]Register, len(n.Inputs))
		for i, binding := range n.Inputs {
			decl := binding.Name.Decl
			utils.Assert(decl != nil, "An unresolved identifier survived until bytecode generation")

			reg, exists := p.Registers[decl]
			utils.Assert(exists, "A register was not allocated for a declaration before use in inline assembly")

			asm.InputRegisters[i] = reg
		}

		if len(n.Outputs) > 0 {
			binding := n.Outputs[0] // currently limited to one output

			decl := binding.Name.Decl
			utils.Assert(decl != nil, "An unresolved identifier survived until bytecode generation")

			reg, exists := p.Registers[decl]
			utils.Assert(exists, "A register was not allocated for a declaration before use in inline assembly")

			asm.OutputRegister = reg
			asm.OutputBinding = binding
			asm.HasOutput = true
		}

		instruction := Inst(CALL_ASM, asm)
		p.Instructions = append(p.Instructions, instruction)

	case *ast.ImmutableDecl:
		if defn, ok := n.Defn.(*ast.ConstantDefn); ok {
			p.Extend(defn.Expr)
			p.Registers[n] = p.PrevResult
		}

	case *ast.MutableDecl:
		if n.Expr != nil {
			p.Extend(n.Expr)
			p.Registers[n] = p.PrevResult
		} else {
			p.Registers[n] = Rg(p.AssignLocation(), typeFromAst(n.Expr.GetType()))
			utils.NotImplemented("Bytecode generation for declaration without initialization")
		}

	case *ast.EvalStmt:
		p.Extend(n.Expr)

	case *ast.ReturnStmt:
		// FIXME: because of the current hack used for return values "return" only works if it is the last item in the function!
		p.Extend(n.Value)

	case *ast.AssignStmt:
		utils.Assert(len(n.Left) == len(n.Right), "An unbalanced assignment survived until bytecode generation")

		// simple assignment
		if len(n.Right) == 1 {
			p.Extend(n.Right[0])
			rhs := p.PrevResult

			if expr, ok := n.Left[0].(*ast.Identifier); ok {
				utils.Assert(expr.Decl != nil, "An unresolved identifier survived until bytecode generation")
				lhs, exists := p.Registers[expr.Decl]
				utils.Assert(exists, "A register was not allocated for a name before use in an expression")

				// copy from rhs to lhs (cast as needed)
				rightType := n.Right[0].GetType()
				p.insertCast(rhs, rightType, expr.Type)
				p.Instructions = append(p.Instructions, Inst(COPY, Unary(p.PrevResult, lhs)))
			} else {
				panic("TODO: Handle non-identifier expressions as the assignee in assignment")
			}

			return
		}

		// parallel assignment
		tmps := make([]Register, len(n.Right))
		for i, expr := range n.Right {
			p.Extend(expr)
			rhs := p.PrevResult

			// copy from rhs to temporary
			tmps[i] = Rg(p.AssignLocation(), rhs.Typ)
			p.Instructions = append(p.Instructions, Inst(COPY, Unary(rhs, tmps[i])))
		}
		for i, expr := range n.Left {
			if e, ok := expr.(*ast.Identifier); ok {
				utils.Assert(e.Decl != nil, "An unresolved identifier survived until bytecode generation")
				lhs, exists := p.Registers[e.Decl]
				utils.Assert(exists, "A register was not allocated for a name before use in an expression")

				// copy from temporary to lhs (cast as needed)
				p.insertCast(tmps[i], n.Right[i].GetType(), e.Type)
				p.Instructions = append(p.Instructions, Inst(COPY, Unary(p.PrevResult, lhs)))
			} else {
				panic("TODO: Handle non-identifier expressions as the assignee in assignment")
			}
		}

	case *ast.TextLiteral:
		utils.Assert(n.Value != ast.UnparsedValue, "An unparsed value survived until bytecode generation")
		register := Rg(p.AssignLocation(), Pointer)

		text, ok := n.Value.([]byte)
		utils.Assert(ok, "A text literal is not a byte slice during bytecode generation")

		name := p.Program.NextConstantName()
		instruction := Inst(LOAD, ConstPtr(name, register))

		p.Program.DefineData(name, text)
		p.Instructions = append(p.Instructions, instruction)
		endRegister = register

	case *ast.NumberLiteral:
		utils.Assert(n.Value != ast.UnparsedValue, "An unparsed value survived until bytecode generation")
		register := Rg(p.AssignLocation(), typeFromAst(n.Type))

		switch n.Value.(type) {
		case int64, uint64, float64:
		default:
			utils.AssertionFailed("A number literal is not an int64, uint64, or float64 value during bytecode generation")
		}

		name := p.Program.NextConstantName()
		instruction := Inst(LOAD, Constant(name, register))

		p.Program.DefineData(name, Pack(n.Value))
		p.Instructions = append(p.Instructions, instruction)
		endRegister = register

	case *ast.Identifier:
		utils.Assert(n.Decl != nil, "An unresolved identifier survived until bytecode generation")
		register, exists := p.Registers[n.Decl]
		utils.Assert(exists, "A register was not allocated for a declaration before use in an expression")
		endRegister = register

	case *ast.GroupExpr:
		p.Extend(n.Subexpr)
		endRegister = p.PrevResult

	case *ast.InfixExpr:
		// TODO: casts should probably be added to the AST elsewhere and only processed here

		// instructions to evaluate arguments
		p.Extend(n.Left)
		p.insertCast(p.PrevResult, n.Left.GetType(), n.Type)
		left := p.PrevResult
		p.Extend(n.Right)
		p.insertCast(p.PrevResult, n.Right.GetType(), n.Type)
		right := p.PrevResult

		var op Opcode
		switch n.Type {
		// TODO: probably shouldn't have any "inferred" types by the time we get here
		case ast.InferredNumber, ast.InferredSigned, ast.InferredUnsigned, ast.InferredFloat:
			switch n.Operator {
			case ast.BuiltinAdd:
				op = ADD
			case ast.BuiltinSubtract:
				op = SUBTRACT
			case ast.BuiltinMultiply:
				op = MULTIPLY
			case ast.BuiltinDivide:
				op = DIVIDE
			}
		default:
			panic("TODO: Unhandle expression type in bytecode generator")
		}

		out := Rg(p.AssignLocation(), typeFromAst(n.Type))
		infix := Inst(op, Binary(left, right, out))

		// bytecode to evaluate operator
		p.Instructions = append(p.Instructions, infix)
		endRegister = out

	case *ast.ProcedureExpr:
		proc := p.Program.NewProcedure()
		if defn, ok := n.Parent.(*ast.ConstantDefn); ok {
			if decl, ok := defn.Parent.(*ast.ImmutableDecl); ok {
				p.Program.Text[decl.Name.Literal] = proc.Index
			}
		}

		proc.Arguments = make([]Register, len(n.Params))
		for i, param := range n.Params {
			if param.Expr != nil {
				utils.NotImplemented("Bytecode generation for parameter default values")
			}
			register := Rg(proc.AssignLocation(), typeFromAst(param.Type))
			proc.Arguments[i] = register
			proc.Registers[param] = register
		}
		proc.Extend(n.Block)

	case *ast.CallExpr:
		name, ok := n.Procedure.(*ast.Identifier)
		if !ok {
			utils.NotImplemented("Bytecode generation for procedure calls to procedure pointers; procedure must be known at compile time")
		}

		child := p.Program.Procedures[p.Program.Text[name.Literal]]
		utils.Assert(len(n.Arguments) == len(child.Arguments), "A procedure call with an incorrect number of arguments survived until bytecode generation")

		ins := make([]Register, len(n.Arguments))
		for i, arg := range n.Arguments {
			p.Extend(arg)
			ins[i] = p.PrevResult
		}

		// TODO: figure out whether we actually have a return value or not
		out := Rg(-1, None)
		// if HAS_RETURN {
		out = Rg(p.AssignLocation(), None) // TODO: get the return type of the procedure
		// }
		p.Instructions = append(p.Instructions, Inst(CALL, Proc(child, out, ins)))
		endRegister = out

	default:
		panic("TODO: Unhandled node type in bytecode.Generate")
	}

	// set the result of this expression
	p.PrevResult = endRegister
}

func (p *Procedure) AssignLocation() Location {
	utils.Assert(p.NextFree < OutOfRegisters, "Ran out of assignable registers.")
	location := p.NextFree
	p.NextFree += 1
	return location
}

func typeFromAst(t ast.Type) Type {
	switch t {
	case ast.InferredFloat:
		return Float64
	case ast.InferredUnsigned:
		return Uint64
	case ast.InferredSigned, ast.InferredNumber, ast.BuiltinInt, ast.BuiltinInt64:
		return Int64
	case ast.BuiltinText:
		return Pointer // TODO: eventually text should be a pointer/length struct
	default:
		utils.NotImplemented("bytecode generation for for non-numeric/non-builtin types")
		return None
	}
}

// TODO: Handle nonnumeric types
func (p *Procedure) insertCast(in Register, from ast.Type, to ast.Type) {
	p.PrevResult = in
	if from == to {
		return
	}

	// TODO: maybe insert overflow check for integer conversions?
	// FIXME: right now "inferred numbers" are accept upto uint64 max,
	//        but here I want to (and do) treat them as signed
	switch to {
	case ast.InferredNumber, ast.InferredSigned:
		out := Rg(p.AssignLocation(), Int64)
		cast := Inst(CAST_I64, Unary(in, out))
		p.Instructions = append(p.Instructions, cast)
		p.PrevResult = out
	case ast.InferredUnsigned:
		out := Rg(p.AssignLocation(), Uint64)
		cast := Inst(CAST_U64, Unary(in, out))
		p.Instructions = append(p.Instructions, cast)
		p.PrevResult = out
	case ast.InferredFloat:
		out := Rg(p.AssignLocation(), Float64)
		cast := Inst(CAST_F64, Unary(in, out))
		p.Instructions = append(p.Instructions, cast)
		p.PrevResult = out
	default:
		panic("TODO: unhandled type in insertCast")
	}
}
