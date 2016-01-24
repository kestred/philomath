package bytecode

import (
	"strconv"
	"unsafe"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/utils"
)

/* TODO: Handle shadowed variable names

   I'm putting this off right now because I think that by the time I get to
   the bytecode generator, I should not be operating on variable names at all
   and so I don't want to build a complicated solution for it.

   In the short term, this will be handled (for nested block scopes) as a check
   at the semantic level.  Note that the semantic check is not intended to avoid
   implementing shadowing, but rather because I think the bytecode generator
	 is not the right place to implement variable shadowing;
*/

type Instruction struct {
	Code Code

	Out   Register
	Left  Register
	Right Register
}

type Code uint16
type Register uint16

// Constant converts a constant table index to a register value.
//
// This function's main purpose is to visually distinguish between register
// assignment and constant indexes, and to make it easier to search for its use.
func Constant(i int) Register { return Register(i) }

const OutOfRegisters = 65535
const (
	NOOP Code = iota

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

type Data uintptr

func FromI64(v int64) Data   { return *(*Data)(unsafe.Pointer(&v)) }
func FromU64(v uint64) Data  { return *(*Data)(unsafe.Pointer(&v)) }
func FromF64(v float64) Data { return *(*Data)(unsafe.Pointer(&v)) }
func ToI64(v Data) int64     { return *(*int64)(unsafe.Pointer(&v)) }
func ToU64(v Data) uint64    { return *(*uint64)(unsafe.Pointer(&v)) }
func ToF64(v Data) float64   { return *(*float64)(unsafe.Pointer(&v)) }

type Scope struct {
	Constants    []Data
	NextRegister Register
	Registers    map[string]Register
}

func (s *Scope) Init() {
	s.Constants = []Data{0} // the 0th constant is always 0
	s.NextRegister = 1      // skip the 0th register
	s.Registers = make(map[string]Register)
}

func (s *Scope) AssignRegister() Register {
	utils.Assert(s.NextRegister < OutOfRegisters, "Ran out of assignable registers.")
	register := s.NextRegister
	s.NextRegister += 1
	return register
}

func FromBlock(block *ast.Block, scope *Scope) []Instruction {
	var insts []Instruction
	for _, node := range block.Nodes {
		switch n := node.(type) {
		case *ast.Block:
			insts = append(insts, FromBlock(n, scope)...)
		case *ast.ConstantDecl:
			// NOTE: an ConstantDefn is the only definiton used directly in bytecode generation
			if defn, ok := n.Defn.(*ast.ConstantDefn); ok {
				insts = append(insts, FromExpr(defn.Expr, scope)...)
				scope.Registers[n.Name.Literal] = insts[len(insts)-1].Out
			}
		case *ast.MutableDecl:
			if n.Expr != nil {
				insts = append(insts, FromExpr(n.Expr, scope)...)
				scope.Registers[n.Name.Literal] = insts[len(insts)-1].Out
			} else {
				//// TODO: zero initialization (but it won't matter until there are typed declarations)
				//scope.Registers[n.Name.Literal] = scope.AssignRegister()
				panic("TODO: Unhandled mutable declaration without expression")
			}
		case *ast.ExprStmt:
			insts = append(insts, FromExpr(n.Expr, scope)...)
		case *ast.AssignStmt:
			utils.Assert(len(n.Assignees) == len(n.Values), "An unbalanced assignment survived until bytecode generation")
			if len(n.Values) == 1 {
				insts = append(insts, FromExpr(n.Values[0], scope)...)
				rhs := insts[len(insts)-1].Out
				if expr, ok := n.Assignees[0].(*ast.ValueExpr); ok {
					lit, ok := expr.Literal.(*ast.Identifier)
					utils.Assert(ok, "Found a non-identifier literal as the assignee in assignment")
					lhs, exists := scope.Registers[lit.Literal]
					utils.Assert(exists, "A register was not allocated for a name before use in an expression")

					// copy from rhs to lhs (cast as needed)
					rhs = insertConversion(scope, &insts, rhs, n.Values[0].GetType(), expr.Type)
					insts = append(insts, Instruction{Code: COPY_VALUE, Out: lhs, Left: rhs})
				} else {
					panic("TODO: Handle non-identifier expressions as the assignee in assignment")
				}
			} else {
				tmps := make([]Register, len(n.Values))
				for i, expr := range n.Values {
					insts = append(insts, FromExpr(expr, scope)...)
					rhs := insts[len(insts)-1].Out

					// copy from rhs to temporary
					tmps[i] = scope.AssignRegister()
					insts = append(insts, Instruction{Code: COPY_VALUE, Out: tmps[i], Left: rhs})
				}
				for i, expr := range n.Assignees {
					if e, ok := expr.(*ast.ValueExpr); ok {
						lit, ok := e.Literal.(*ast.Identifier)
						utils.Assert(ok, "Found a non-identifier literal as the assignee in assignment")
						lhs, exists := scope.Registers[lit.Literal]
						utils.Assert(exists, "A register was not allocated for a name before use in an expression")

						// copy from temporary to lhs (cast as needed)
						rhs := insertConversion(scope, &insts, tmps[i], n.Values[i].GetType(), e.Type)
						insts = append(insts, Instruction{Code: COPY_VALUE, Out: lhs, Left: rhs})
					} else {
						panic("TODO: Handle non-identifier expressions as the assignee in assignment")
					}
				}
			}
		default:
			panic("TOOD: Unhandle node type in bytecode.FromBlock")
		}
	}

	return insts
}

func FromExpr(expr ast.Expr, scope *Scope) []Instruction {
	switch e := expr.(type) {
	case *ast.ValueExpr:
		switch lit := e.Literal.(type) {
		case *ast.NumberLiteral:
			utils.Assert(lit.Value != ast.UnparsedValue, "A value was not parsed before bytecode generation")
			register := scope.AssignRegister()

			var value Data
			switch v := lit.Value.(type) {
			case int64:
				value = FromI64(v)
			case uint64:
				value = FromU64(v)
			case float64:
				value = FromF64(v)
			default:
				panic("TODO: Unhandled value type")
			}

			nextConstant := len(scope.Constants)
			scope.Constants = append(scope.Constants, Data(value))
			return []Instruction{{Code: LOAD_CONST, Out: register, Left: Constant(nextConstant)}}
		case *ast.Identifier:
			register, exists := scope.Registers[lit.Literal]
			utils.Assert(exists, "A register was not allocated for a name before use in an expression")
			// FIXME: currently expressions must return an instruction with the Out register set,
			//        but it doesn't make a whole lot of sense to be emitting NOOPs for value loads
			return []Instruction{{Code: NOOP, Out: register}}
		default:
			panic("TODO: Unhandled value literal")
		}

	case *ast.GroupExpr:
		return FromExpr(e.Subexpr, scope)

	case *ast.InfixExpr:
		var insts []Instruction
		var infix Instruction

		// TODO: casts should probably be added to the AST elsewhere and only processed here

		// instructions to evaluate left hand side
		left := FromExpr(e.Left, scope)
		insts = append(insts, left...)
		infix.Left = insertConversion(scope, &insts, insts[len(insts)-1].Out, e.Left.GetType(), e.Type)

		// instructions to evaluate right hand side
		right := FromExpr(e.Right, scope)
		insts = append(insts, right...)
		infix.Right = insertConversion(scope, &insts, insts[len(insts)-1].Out, e.Right.GetType(), e.Type)

		// TODO: table lookup?
		// TODO: probably shouldn't have any "inferred" types by the time we get here
		// TODO: probably shouldn't be comparing against operator literals by the time we get here
		switch e.Operator.Literal {
		case "+":
			switch e.Type {
			// FIXME: right now "inferred numbers" are accept upto uint64 max,
			//        but here I want to (and do) treat them as signed
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
			switch e.Type {
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
			switch e.Type {
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
			switch e.Type {
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

// TODO: Handle non-numeric types
func insertConversion(scope *Scope, insts *[]Instruction, loc Register, from ast.Type, to ast.Type) Register {
	if from == to {
		return loc
	}

	// TODO: maybe insert overflow check for integer conversions?
	// TODO: table lookup?
	switch from {
	// FIXME: right now "inferred numbers" are accept upto uint64 max,
	//        but here I want to (and do) treat them as signed
	case ast.InferredNumber, ast.InferredSigned:
		if to == ast.InferredFloat {
			list := *insts
			convert := Instruction{Code: CONVERT_I64_TO_F64, Out: scope.AssignRegister()}
			convert.Left = loc
			*insts = append(list, convert)
			return convert.Out
		}
	case ast.InferredUnsigned:
		if to == ast.InferredFloat {
			list := *insts
			convert := Instruction{Code: CONVERT_U64_TO_F64, Out: scope.AssignRegister()}
			convert.Left = loc
			*insts = append(list, convert)
			return convert.Out
		}
	case ast.InferredFloat:
		switch to {
		case ast.InferredNumber, ast.InferredSigned:
			list := *insts
			convert := Instruction{Code: CONVERT_F64_TO_I64, Out: scope.AssignRegister()}
			convert.Left = loc
			*insts = append(list, convert)
			return convert.Out
		case ast.InferredUnsigned:
			list := *insts
			convert := Instruction{Code: CONVERT_F64_TO_U64, Out: scope.AssignRegister()}
			convert.Left = loc
			*insts = append(list, convert)
			return convert.Out
		}
	}

	return loc
}
