package bytecode

import (
	"strconv"
	"unsafe"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/utils"
)

type Instruction struct {
	Op    Opcode
	Out   Register
	Left  Register
	Right Register
}

type Opcode int16
type Register int16

// Constant and Metadata are used to visually distinguish between register
// assignment and constant indexes, and to make it easier to search for their use.
func Constant(i int) Register { return Register(i) }
func Metadata(i int) Register { return Register(i) }

const OutOfRegisters = 32767
const (
	NOOP Opcode = iota

	/* ASIDE: MOVE vs COPY

	   I had a hard time deciding whether to call this instruction "move" or "copy".
	   In other bytecode representions or instruction sets, the operation is
	   called MOVE, but the description is then "copy A to B".

	   Semantically, a proper MOVE opcode means that value in the Left register is
	   moved to the Out register and the Left register becomes an undefined value.
	   In practice, the interperter would just leave the value in the left register
		 except during debugging, when it could be overwritten with a canary value.

	   For now, I've chosen a COPY opcode because the bytecode generator is
	   currently treating each register similar to single-static assignment form.
	   If I have a smaller register pool and start re-using registers later,
		 then I may end up switching to or adding a MOVE opcode.
	*/
	COPY_VALUE
	LOAD_CONST

	CALL
	CALL_ASM
	RETURN

	I64_ADD
	I64_SUBTRACT
	I64_MULTIPLY
	I64_DIVIDE

	U64_ADD
	U64_SUBTRACT
	U64_MULTIPLY
	U64_DIVIDE

	F64_ADD
	F64_SUBTRACT
	F64_MULTIPLY
	F64_DIVIDE

	CONVERT_F64_TO_I64
	CONVERT_F64_TO_U64
	CONVERT_I64_TO_F64
	CONVERT_U64_TO_F64
)

var opcodes = [...]string{
	NOOP: "No operation",

	COPY_VALUE: "Copy value",
	LOAD_CONST: "Load constant",

	CALL:     "Call procedure",
	CALL_ASM: "Call assembly",
	RETURN:   "Return",

	I64_ADD:      "Signed Addition",
	I64_SUBTRACT: "Signed Subtraction",
	I64_MULTIPLY: "Signed Multiplication",
	I64_DIVIDE:   "Signed Division",

	U64_ADD:      "Unsigned Addition",
	U64_SUBTRACT: "Unsigned Subtraction",
	U64_MULTIPLY: "Unsigned Multiplication",
	U64_DIVIDE:   "Unsigned Division",

	F64_ADD:      "Float Addition",
	F64_SUBTRACT: "Float Subtraction",
	F64_MULTIPLY: "Float Multiplication",
	F64_DIVIDE:   "Float Division",

	CONVERT_F64_TO_I64: "Truncate float to signed",
	CONVERT_F64_TO_U64: "Truncate float to unsigned",
	CONVERT_I64_TO_F64: "Convert signed to float",
	CONVERT_U64_TO_F64: "Convert unsigned to float",
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

type Data uintptr

func FromI64(v int64) Data   { return *(*Data)(unsafe.Pointer(&v)) }
func FromU64(v uint64) Data  { return *(*Data)(unsafe.Pointer(&v)) }
func FromF64(v float64) Data { return *(*Data)(unsafe.Pointer(&v)) }
func ToI64(v Data) int64     { return *(*int64)(unsafe.Pointer(&v)) }
func ToU64(v Data) uint64    { return *(*uint64)(unsafe.Pointer(&v)) }
func ToF64(v Data) float64   { return *(*float64)(unsafe.Pointer(&v)) }

type Program struct {
	Constants  []Data
	Procedures []Procedure
	Metadata   []interface{}
	Data       map[string]int // offset into Constants
	Text       map[string]int // offset into Procedures
}

func NewProgram() *Program {
	prog := &Program{}
	prog.Data = map[string]int{}
	prog.Text = map[string]int{"start_": 0}
	prog.AddConstant(0)
	prog.AddMetadata(nil)
	prog.NewProcedure() // start_
	return prog
}

func (p *Program) Extend(node ast.Node) {
	p.Procedures[0].Extend(node)
}

func (p *Program) AddConstant(data Data) int {
	next := len(p.Constants)
	p.Constants = append(p.Constants, data)
	return next
}

func (p *Program) AddMetadata(data interface{}) int {
	next := len(p.Metadata)
	p.Metadata = append(p.Metadata, data)
	return next
}

func (p *Program) NewProcedure() *Procedure {
	next := len(p.Procedures)
	p.Procedures = append(p.Procedures, Procedure{
		Index:        next,
		Program:      p,
		Instructions: []Instruction{},
		Registers:    map[ast.Decl]Register{},
		NextRegister: +0,
		ExprRegister: -1,
	})
	return &p.Procedures[next]
}

type Assembly struct {
	Source  string
	Wrapper unsafe.Pointer // for interpreter
	Tempdir string

	HasOutput      bool
	OutputBinding  ast.AsmBinding
	OutputRegister Register
	InputBindings  []ast.AsmBinding
	InputRegisters []Register
}

type Procedure struct {
	Index        int
	Program      *Program
	Instructions []Instruction

	// state for bytecode generation
	Registers    map[ast.Decl]Register
	NextRegister Register // next assignable
	ExprRegister Register // last result
}

func (p *Procedure) Extend(node ast.Node) {
	endRegister := Register(-1)
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
		asm := new(Assembly)
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

		metadata := p.Program.AddMetadata(asm)
		instruction := Instruction{Op: CALL_ASM, Left: Metadata(metadata)}
		p.Instructions = append(p.Instructions, instruction)

	case *ast.ImmutableDecl:
		if defn, ok := n.Defn.(*ast.ConstantDefn); ok {
			p.Extend(defn.Expr)
			p.Registers[n] = p.ExprRegister
		}

	case *ast.MutableDecl:
		if n.Expr != nil {
			p.Extend(n.Expr)
			p.Registers[n] = p.ExprRegister
		} else {
			//// TODO: zero initialization (but it won't matter until there are type declarations)
			//p.Registers[n] = p.AssignRegister()
			panic("TODO: Unhandled mutable declaration without an expression")
		}

	case *ast.EvalStmt:
		p.Extend(n.Expr)

	case *ast.AssignStmt:
		utils.Assert(len(n.Left) == len(n.Right), "An unbalanced assignment survived until bytecode generation")

		// simple assignment
		if len(n.Right) == 1 {
			p.Extend(n.Right[0])
			rhs := p.ExprRegister

			if expr, ok := n.Left[0].(*ast.Identifier); ok {
				utils.Assert(expr.Decl != nil, "An unresolved identifier survived until bytecode generation")
				lhs, exists := p.Registers[expr.Decl]
				utils.Assert(exists, "A register was not allocated for a name before use in an expression")

				// copy from rhs to lhs (cast as needed)
				p.insertConvert(rhs, n.Right[0].GetType(), expr.Type)
				p.Instructions = append(p.Instructions, Instruction{Op: COPY_VALUE, Out: lhs, Left: p.ExprRegister})
			} else {
				panic("TODO: Handle non-identifier expressions as the assignee in assignment")
			}

			return
		}

		// parallel assignment
		tmps := make([]Register, len(n.Right))
		for i, expr := range n.Right {
			p.Extend(expr)
			rhs := p.ExprRegister

			// copy from rhs to temporary
			tmps[i] = p.AssignRegister()
			p.Instructions = append(p.Instructions, Instruction{Op: COPY_VALUE, Out: tmps[i], Left: rhs})
		}
		for i, expr := range n.Left {
			if e, ok := expr.(*ast.Identifier); ok {
				utils.Assert(e.Decl != nil, "An unresolved identifier survived until bytecode generation")
				lhs, exists := p.Registers[e.Decl]
				utils.Assert(exists, "A register was not allocated for a name before use in an expression")

				// copy from temporary to lhs (cast as needed)
				p.insertConvert(tmps[i], n.Right[i].GetType(), e.Type)
				p.Instructions = append(p.Instructions, Instruction{Op: COPY_VALUE, Out: lhs, Left: p.ExprRegister})
			} else {
				panic("TODO: Handle non-identifier expressions as the assignee in assignment")
			}
		}

	case *ast.TextLiteral:
		utils.Assert(n.Value != ast.UnparsedValue, "An unparsed value survived until bytecode generation")
		register := p.AssignRegister()

		constant := p.Program.AddConstant(Data(unsafe.Pointer(&n.Value.([]byte)[0])))
		p.Program.AddMetadata(n.Value)

		instruction := Instruction{Op: LOAD_CONST, Out: register, Left: Constant(constant)}
		p.Instructions = append(p.Instructions, instruction)
		endRegister = register

	case *ast.NumberLiteral:
		utils.Assert(n.Value != ast.UnparsedValue, "An unparsed value survived until bytecode generation")
		register := p.AssignRegister()

		var value Data
		switch v := n.Value.(type) {
		case int64:
			value = FromI64(v)
		case uint64:
			value = FromU64(v)
		case float64:
			value = FromF64(v)
		default:
			panic("TODO: Unhandled value type")
		}

		constant := p.Program.AddConstant(Data(value))
		instruction := Instruction{Op: LOAD_CONST, Out: register, Left: Constant(constant)}
		p.Instructions = append(p.Instructions, instruction)
		endRegister = register

	case *ast.Identifier:
		utils.Assert(n.Decl != nil, "An unresolved identifier survived until bytecode generation")
		register, exists := p.Registers[n.Decl]
		utils.Assert(exists, "A register was not allocated for a declaration before use in an expression")
		endRegister = register

	case *ast.GroupExpr:
		p.Extend(n.Subexpr)
		endRegister = p.ExprRegister

	case *ast.InfixExpr:
		// TODO: casts should probably be added to the AST elsewhere and only processed here
		var infix Instruction

		// instructions to evaluate arguments
		p.Extend(n.Left)
		p.insertConvert(p.ExprRegister, n.Left.GetType(), n.Type)
		infix.Left = p.ExprRegister
		p.Extend(n.Right)
		p.insertConvert(p.ExprRegister, n.Right.GetType(), n.Type)
		infix.Right = p.ExprRegister

		// TODO: table lookup?
		// TODO: probably shouldn't have any "inferred" types by the time we get here
		switch n.Operator {
		case ast.BuiltinAdd:
			switch n.Type {
			// FIXME: right now "inferred numbers" are accept upto uint64 max,
			//        but here I want to (and do) treat them as signed
			case ast.InferredNumber, ast.InferredSigned:
				infix.Op = I64_ADD
			case ast.InferredUnsigned:
				infix.Op = U64_ADD
			case ast.InferredFloat:
				infix.Op = F64_ADD
			default:
				panic("TODO: Unhandle expression type in bytecode generator")
			}
		case ast.BuiltinSubtract:
			switch n.Type {
			case ast.InferredNumber, ast.InferredSigned:
				infix.Op = I64_SUBTRACT
			case ast.InferredUnsigned:
				infix.Op = U64_SUBTRACT
			case ast.InferredFloat:
				infix.Op = F64_SUBTRACT
			default:
				panic("TODO: Unhandle expression type in bytecode generator")
			}
		case ast.BuiltinMultiply:
			switch n.Type {
			case ast.InferredNumber, ast.InferredSigned:
				infix.Op = I64_MULTIPLY
			case ast.InferredUnsigned:
				infix.Op = U64_MULTIPLY
			case ast.InferredFloat:
				infix.Op = F64_MULTIPLY
			default:
				panic("TODO: Unhandle expression type in bytecode generator")
			}
		case ast.BuiltinDivide:
			switch n.Type {
			case ast.InferredNumber, ast.InferredSigned:
				infix.Op = I64_DIVIDE
			case ast.InferredUnsigned:
				infix.Op = U64_DIVIDE
			case ast.InferredFloat:
				infix.Op = F64_DIVIDE
			default:
				panic("TODO: Unhandle expression type in bytecode generator")
			}
		}

		// bytecode to evaluate operator
		infix.Out = p.AssignRegister()
		p.Instructions = append(p.Instructions, infix)
		endRegister = infix.Out

	case *ast.ProcedureExpr:
		proc := p.Program.NewProcedure()
		if defn, ok := n.Parent.(*ast.ConstantDefn); ok {
			if decl, ok := defn.Parent.(*ast.ImmutableDecl); ok {
				p.Program.Text[decl.Name.Literal] = proc.Index
			}
		}
		proc.Extend(n.Block)

	default:
		panic("TODO: Unhandled node type in bytecode.Generate")
	}

	// set the result of this expression
	p.ExprRegister = endRegister
}

func (p *Procedure) AssignRegister() Register {
	utils.Assert(p.NextRegister < OutOfRegisters, "Ran out of assignable registers.")
	register := p.NextRegister
	p.NextRegister += 1
	return register
}

// TODO: Handle non-numeric types
func (p *Procedure) insertConvert(loc Register, from ast.Type, to ast.Type) {
	p.ExprRegister = loc
	if from == to {
		return
	}

	// TODO: maybe insert overflow check for integer conversions?
	// TODO: table lookup?
	switch from {
	// FIXME: right now "inferred numbers" are accept upto uint64 max,
	//        but here I want to (and do) treat them as signed
	case ast.InferredNumber, ast.InferredSigned:
		if to == ast.InferredFloat {
			convert := Instruction{Op: CONVERT_I64_TO_F64, Out: p.AssignRegister()}
			convert.Left = loc
			p.Instructions = append(p.Instructions, convert)
			p.ExprRegister = convert.Out
		}
	case ast.InferredUnsigned:
		if to == ast.InferredFloat {
			convert := Instruction{Op: CONVERT_U64_TO_F64, Out: p.AssignRegister()}
			convert.Left = loc
			p.Instructions = append(p.Instructions, convert)
			p.ExprRegister = convert.Out
		}
	case ast.InferredFloat:
		switch to {
		case ast.InferredNumber, ast.InferredSigned:
			convert := Instruction{Op: CONVERT_F64_TO_I64, Out: p.AssignRegister()}
			convert.Left = loc
			p.Instructions = append(p.Instructions, convert)
			p.ExprRegister = convert.Out
		case ast.InferredUnsigned:
			convert := Instruction{Op: CONVERT_F64_TO_U64, Out: p.AssignRegister()}
			convert.Left = loc
			p.Instructions = append(p.Instructions, convert)
			p.ExprRegister = convert.Out
		}
	}
}
