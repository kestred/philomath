package bytecode

import (
	"strconv"
	"unsafe"

	"github.com/kestred/philomath/ast"
	"github.com/kestred/philomath/utils"
)

type Code uint16
type Register uint16
type Data uintptr

// Constant converts a constant table index to a register value.
//
// This function's main purpose is to visually distinguish between register
// assignment and constant indexes, and to make it easier to search for its use.
func Constant(i int) Register { return Register(i) }

type Instruction struct {
	Code Code

	Out   Register
	Left  Register
	Right Register
}

const OutOfRegisters = 65535

const (
	NOOP Code = iota

	LOAD_CONST

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
	NOOP:       "No operation",
	LOAD_CONST: "Load constant",

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

func (code Code) String() string {
	s := ""
	if 0 <= code && code < Code(len(opcodes)) {
		s = opcodes[code]
	}
	if s == "" {
		s = "Code(" + strconv.Itoa(int(code)) + ")"
	}
	return s
}

type Scope struct {
	Constants    []Data
	Registers    map[string]Register
	NextRegister Register
}

func (s *Scope) Init() {
	s.Constants = []Data{0}
	s.Registers = make(map[string]Register)
	s.NextRegister = 1 // skip the 0th register
}

func (s *Scope) AssignRegister() Register {
	utils.Assert(s.NextRegister < OutOfRegisters, "Ran out of assignable registers.")
	register := s.NextRegister
	s.NextRegister += 1
	return register
}

func FromExpr(expr ast.Expr, scope *Scope) []Instruction {
	switch node := expr.(type) {

	case *ast.ValueExpr:
		switch literal := node.Literal.(type) {
		case *ast.NumberLiteral:
			utils.Assert(literal.Value != ast.UnparsedValue, "A value was not parsed before bytecode generation")
			register := scope.AssignRegister()

			var value Data
			switch v := literal.Value.(type) {
			// NOTE: These can't be combined for some noxious reason
			case int64:
				value = *(*Data)(unsafe.Pointer(&v))
			case uint64:
				value = *(*Data)(unsafe.Pointer(&v))
			case float64:
				value = *(*Data)(unsafe.Pointer(&v))
			default:
				panic("TODO: Unhandled value type")
			}

			nextConstant := len(scope.Constants)
			scope.Constants = append(scope.Constants, Data(value))
			return []Instruction{{Code: LOAD_CONST, Out: register, Left: Constant(nextConstant)}}
		case *ast.Ident:
			panic("TODO: I haven't done declarations... so this identifier isn't that useful")
			// NOTE: Find or assign register
			/*
				var exists bool
				register, exists = scope.Registers[literal]
				if !exists {
					register = scope.AssignRegister()
					scope.Registers[literal] = register
				}
			*/
		default:
			panic("TODO: Unhandled value literal")
		}

	case *ast.GroupExpr:
		return FromExpr(node.Subexpr, scope)

	case *ast.InfixExpr:
		var insts []Instruction
		var infix Instruction

		// TODO: casts should probably be added to the AST elsewhere and only processed here

		// instructions to evaluate left hand side
		left := FromExpr(node.Left, scope)
		insts = append(insts, left...)
		insertConversion(scope, &insts, node.Left.GetType(), node.Type)
		infix.Left = insts[len(insts)-1].Out

		// instructions to evaluate right hand side
		right := FromExpr(node.Right, scope)
		insts = append(insts, right...)
		insertConversion(scope, &insts, node.Right.GetType(), node.Type)
		infix.Right = insts[len(insts)-1].Out

		// TODO: table lookup?
		// TODO: probably shouldn't have any "inferred" types by the time we get here
		// TODO: probably shouldn't be comparing against operator literals by the time we get here
		switch node.Operator.Literal {
		case "+":
			switch node.Type {
			// FIXME: right now "inferred numbers" are parsed into uint64, but here I treat them as signed
			case ast.InferredNumber, ast.InferredSigned:
				infix.Code = I64_ADD
			case ast.InferredUnsigned:
				infix.Code = U64_ADD
			case ast.InferredFloat:
				infix.Code = F64_ADD
			default:
				panic("TODO: Unhandle expression type in bytecode generator")
			}
		case "-":
			switch node.Type {
			case ast.InferredNumber, ast.InferredSigned:
				infix.Code = I64_SUBTRACT
			case ast.InferredUnsigned:
				infix.Code = U64_SUBTRACT
			case ast.InferredFloat:
				infix.Code = F64_SUBTRACT
			default:
				panic("TODO: Unhandle expression type in bytecode generator")
			}
		case "*":
			switch node.Type {
			case ast.InferredNumber, ast.InferredSigned:
				infix.Code = I64_MULTIPLY
			case ast.InferredUnsigned:
				infix.Code = U64_MULTIPLY
			case ast.InferredFloat:
				infix.Code = F64_MULTIPLY
			default:
				panic("TODO: Unhandle expression type in bytecode generator")
			}
		case "/":
			switch node.Type {
			case ast.InferredNumber, ast.InferredSigned:
				infix.Code = I64_DIVIDE
			case ast.InferredUnsigned:
				infix.Code = U64_DIVIDE
			case ast.InferredFloat:
				infix.Code = F64_DIVIDE
			default:
				panic("TODO: Unhandle expression type in bytecode generator")
			}
		}

		// bytecode to evaluate operator
		infix.Out = scope.AssignRegister()
		insts = append(insts, infix)
		return insts

	default:
		panic("TODO: Unhandled expression type")
	}
}

func insertConversion(scope *Scope, insts *[]Instruction, from ast.Type, to ast.Type) {
	if from == to {
		return
	}

	// TODO: maybe insert overflow check for integer conversions?
	//   I don't think lossy conversions will be allowed by the type checker so for now, do nothing

	// TODO: table lookup?
	switch from {
	// FIXME: right now "inferred numbers" are parsed into uint64, but here I treat them as signed
	case ast.InferredNumber, ast.InferredSigned:
		if to == ast.InferredFloat {
			list := *insts
			convert := Instruction{Code: CONVERT_I64_TO_F64, Out: scope.AssignRegister()}
			convert.Left = list[len(list)-1].Out
			*insts = append(list, convert)
		}
	case ast.InferredUnsigned:
		if to == ast.InferredFloat {
			list := *insts
			convert := Instruction{Code: CONVERT_U64_TO_F64, Out: scope.AssignRegister()}
			convert.Left = list[len(list)-1].Out
			*insts = append(list, convert)
		}
	case ast.InferredFloat:
		switch to {
		case ast.InferredNumber, ast.InferredSigned:
			list := *insts
			convert := Instruction{Code: CONVERT_F64_TO_I64, Out: scope.AssignRegister()}
			convert.Left = list[len(list)-1].Out
			*insts = append(list, convert)
		case ast.InferredUnsigned:
			list := *insts
			convert := Instruction{Code: CONVERT_F64_TO_U64, Out: scope.AssignRegister()}
			convert.Left = list[len(list)-1].Out
			*insts = append(list, convert)
		}
	}
}
