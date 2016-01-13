package bytecode

import (
	"strconv"

	"github.com/kestred/philomath/ast"
)

type Code uint16
type Register uint16
type Data uintptr

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

	INT64_ADD
	INT64_SUBTRACT
	INT64_MULTIPLY
	INT64_DIVIDE

	/*
		Uint64Add
		Uint64Subtract
		Uint64Multiply
		Uint64Divide

		Float64Add
		Float64Subtract
		Float64Multiply
		Float64Divide
	*/
)

var opcodes = [...]string{
	NOOP:       "No operation",
	LOAD_CONST: "Load constant",

	INT64_ADD:      "Addition",
	INT64_SUBTRACT: "Subtraction",
	INT64_MULTIPLY: "Multiplication",
	INT64_DIVIDE:   "Division",
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
	if s.NextRegister == OutOfRegisters {
		panic("Ran out of registers.  TODO: Register re-use")
	}

	register := s.NextRegister
	s.NextRegister += 1
	return register
}

func FromExpr(expr ast.Expr, scope *Scope) []Instruction {
	switch node := expr.(type) {

	case *ast.ValueExpr:
		switch literal := node.Literal.(type) {
		case *ast.NumberLiteral:
			register := scope.AssignRegister()
			value, err := strconv.Atoi(literal.Literal)
			if err != nil {
				panic("TODO: Actually perform type checking, etc")
			}

			constIndex := Register(len(scope.Constants))
			scope.Constants = append(scope.Constants, Data(value))
			return []Instruction{{Code: LOAD_CONST, Out: register, Left: constIndex}}
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
		var infix Instruction
		switch node.Operator.Literal {
		case "+":
			infix.Code = INT64_ADD
		case "-":
			infix.Code = INT64_SUBTRACT
		case "*":
			infix.Code = INT64_MULTIPLY
		case "/":
			infix.Code = INT64_DIVIDE
		}

		var insts []Instruction

		// bytecode to evaluate left
		left := FromExpr(node.Left, scope)
		infix.Left = left[len(left)-1].Out
		insts = append(insts, left...)

		// bytecode to evaluate right
		right := FromExpr(node.Right, scope)
		infix.Right = right[len(right)-1].Out
		insts = append(insts, right...)

		// bytecode to evaluate operator
		infix.Out = scope.AssignRegister()
		insts = append(insts, infix)
		return insts

	default:
		panic("TODO: Unhandled expression type")
	}
}
