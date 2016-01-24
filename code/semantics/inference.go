package semantics

// NOTE: Only performing inference; type checking is a separate step
// TODO: Break-out operator overload resolution
// TODO: Should literals continue to be parsed here, or elsewhere?
// TODO: Stop assuming declarations will be in order.
//       They must be in blocks, but do not need to be at file/module scope

import (
	"strconv"
	"strings"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/utils"
)

func InferTypes(node ast.Node) ast.Type {
	context := MakeTypeMap()
	switch n := node.(type) {
	case *ast.Block:
		return inferTypesInBlock(n, context)
	case ast.Expr:
		return inferTypesInExpr(n, context)
	default:
		panic("TODO: Only pass me a Block or an Expr right now!")
	}
}

func inferTypesInBlock(block *ast.Block, context TypeMap) ast.Type {
	// NOTE: If we have first-class union types, its reasonable to say that the
	//       type of a block is the union of all return statement expression types.
	//       Otherwise, still useful to collect for error messages.
	// var returnTypes []ast.Type
	for _, node := range block.Nodes {
		switch n := node.(type) {
		case *ast.Block:
			inferTypesInBlock(n, context.Reference())
		case *ast.ConstantDecl:
			// NOTE: An ExprDefn is the only definiton with a type
			if defn, ok := n.Defn.(*ast.ExprDefn); ok {
				context.Set(n.Name.Literal, inferTypesInExpr(defn.Expr, context))
			}
		case *ast.MutableDecl:
			typ := inferTypesInExpr(n.Expr, context)
			// TODO: Not sure that I should replace the type of the decl if it is inferred...
			if n.Type == ast.InferredType {
				context.Set(n.Name.Literal, typ)
			} else {
				context.Set(n.Name.Literal, n.Type)
			}
		case *ast.ExprStmt:
			inferTypesInExpr(n.Expr, context)
		case *ast.AssignStmt:
			if len(n.Assignees) != len(n.Values) {
				panic("TODO: Handle unbalanced assignment")
			}

			for i := range n.Assignees {
				inferTypesInExpr(n.Assignees[i], context)
				inferTypesInExpr(n.Values[i], context)
			}
		case *ast.ReturnStmt:
			// NOTE: Once we actually have a return statement, the type of its
			//       expression would be this block's type; otherwise we return none.
		default:
			panic("TOOD: Unhandled stmt/decl in semantics.inferTypesInBlock")
		}
	}

	return ast.BuiltinEmpty
}

// TODO: Probably need to be careful to avoid smashing non-inferred types
//       I'll just wait until it becomes a bug and deal with it then
func inferTypesInExpr(expr ast.Expr, context TypeMap) ast.Type {
	switch e := expr.(type) {
	case *ast.PostfixExpr:
		subtype := inferTypesInExpr(e.Subexpr, context)
		e.Type = inferPostfixType(e.Operator, subtype)
		return e.Type
	case *ast.InfixExpr:
		left := inferTypesInExpr(e.Left, context)
		right := inferTypesInExpr(e.Right, context)
		e.Type = inferInfixType(e.Operator, left, right)
		return e.Type
	case *ast.PrefixExpr:
		subtype := inferTypesInExpr(e.Subexpr, context)
		e.Type = inferPrefixType(e.Operator, subtype)
		return e.Type
	case *ast.GroupExpr:
		e.Type = inferTypesInExpr(e.Subexpr, context)
		return e.Type
	case *ast.ValueExpr:
		switch literal := e.Literal.(type) {
		case *ast.Identifier:
			// NOTE: Assuming declarations will be in order (will stop being true eventually)
			typ := context.Get(literal.Literal)
			if typ != nil {
				e.Type = typ
			} else {
				e.Type = ast.UnknownType
			}
			return e.Type
		case *ast.NumberLiteral:
			var err error
			lit := literal.Literal
			if len(lit) > 2 && lit[0:2] == "0x" {
				e.Type = ast.InferredUnsigned
				literal.Value, err = strconv.ParseUint(lit, 16, 0)
				if err == strconv.ErrRange {
					panic("TODO: Handle hexadecimal literal can't be represented by uint64")
				}
				utils.AssertNil(err, "Failed parsing hexadecimal literal")
			} else if strings.Contains(lit, ".") || strings.Contains(lit, "e") {
				e.Type = ast.InferredFloat
				literal.Value, err = strconv.ParseFloat(lit, 0)
				if err == strconv.ErrRange {
					panic("TODO: Handle floating point literal can't be represented by float64")
				}
				utils.AssertNil(err, "Failed parsing float literal")
			} else if lit[0] == '0' {
				e.Type = ast.InferredUnsigned
				literal.Value, err = strconv.ParseUint(lit, 8, 0)
				if err == strconv.ErrRange {
					panic("TODO: Handle octal literal can't be represented by uint64")
				}
				utils.AssertNil(err, "Failed parsing octal literal")
			} else {
				e.Type = ast.InferredNumber
				literal.Value, err = strconv.ParseUint(lit, 10, 0)
				if err == strconv.ErrRange {
					panic("TODO: Handle decimal literal can't be represented by uint64")
				}
				utils.AssertNil(err, "Failed parsing decimal literal")
			}
			return e.Type
		default:
			panic("TODO: Unhandled value literal")
		}
	default:
		panic("TODO: Handle type inferences for this expr")
	}
}

func inferPrefixType(op *ast.OperatorDefn, typ ast.Type) ast.Type {
	if typ == ast.UnknownType {
		return typ
	}

	switch op {
	case ast.BuiltinPositive, ast.BuiltinNegative:
		switch typ {
		case ast.InferredNumber:
			return ast.InferredSigned
		case
			ast.InferredFloat,
			ast.BuiltinFloat,
			ast.BuiltinFloat32,
			ast.BuiltinFloat64,
			ast.InferredSigned,
			ast.BuiltinInt,
			ast.BuiltinInt8,
			ast.BuiltinInt16,
			ast.BuiltinInt32,
			ast.BuiltinInt64:
			return typ
		case ast.InferredUnsigned:
			return ast.InferredSigned
		case ast.BuiltinUint:
			return ast.BuiltinInt
		case ast.BuiltinUint8:
			return ast.BuiltinInt8
		case ast.BuiltinUint16:
			return ast.BuiltinInt16
		case ast.BuiltinUint32:
			return ast.BuiltinInt32
		case ast.BuiltinUint64:
			return ast.BuiltinInt64
		default:
			panic("TODO: Implement operator overload resolution for prefix -/+")
		}
	default:
		panic("TODO: Unhandled prefix operator in type inference")
	}
}

func inferPostfixType(op *ast.OperatorDefn, typ ast.Type) ast.Type {
	if typ == ast.UnknownType {
		return typ
	}

	switch op {
	default:
		panic("TODO: Unhandled postfix operator in type inference")
	}
}

func inferInfixType(op *ast.OperatorDefn, left ast.Type, right ast.Type) ast.Type {
	if left == ast.UnknownType || right == ast.UnknownType {
		return ast.UnknownType
	}

	switch op {
	case ast.BuiltinAdd, ast.BuiltinSubtract, ast.BuiltinMultiply, ast.BuiltinDivide:
		// TODO: Implement operator overload resolution
		return castNumbers(left, right)
	default:
		panic("TODO: Unhandled infix operator in type inference")
	}
}

func castNumbers(left ast.Type, right ast.Type) ast.Type {
	// NOTE: typechecking will come through later and assert that the implicit
	//       casts are either safe (or "safe-enough")
	if !isNumber(left) || !isNumber(right) {
		// can't cast non-number to number
		return ast.UnknownType
	}

	if isFloat(left) {
		if isFloat(right) {
			// cast low-bit to high-bit float
			return promoteByOrder(left, right)
		} else {
			// cast any integer to any float
			return left
		}
	} else if isFloat(right) {
		// cast any integer to any float
		return right
	}

	if isSigned(left) {
		if isSigned(right) {
			// cast low-bit integer to high-bit integer
			return promoteByOrder(left, right)
		} else if right == ast.InferredNumber {
			// cast generic to signed
			return left
		} else {
			// can't cast unsigned to signed
			return ast.UnknownType
		}
	} else if isSigned(right) {
		if left == ast.InferredNumber {
			// cast generic to signed
			return right
		} else {
			// can't cast unsigned to signed
			return ast.UnknownType
		}
	}

	// cast low-bit unsigned to high-bit unsigned
	return promoteByOrder(left, right)
}

func promoteByOrder(left ast.Type, right ast.Type) ast.Type {
	if promotionOrder(left) >= promotionOrder(right) {
		return left
	} else {
		return right
	}
}

func promotionOrder(typ ast.Type) int {
	switch typ {
	case ast.InferredNumber:
		return 0
	case ast.InferredFloat, ast.InferredSigned, ast.InferredUnsigned:
		return 1
	case ast.BuiltinInt8, ast.BuiltinUint8:
		return 2
	case ast.BuiltinInt16, ast.BuiltinUint16:
		return 3
	case ast.BuiltinFloat32, ast.BuiltinInt32, ast.BuiltinUint32:
		return 4
	case ast.BuiltinFloat64, ast.BuiltinInt64, ast.BuiltinUint64:
		return 5
	case ast.BuiltinFloat, ast.BuiltinInt, ast.BuiltinUint:
		return 6
	default:
		return -1
	}
}

func isNumber(typ ast.Type) bool {
	return (typ == ast.InferredNumber) || isFloat(typ) || isSigned(typ) || isUnsigned(typ)
}

func isFloat(typ ast.Type) bool {
	switch typ {
	case
		ast.InferredFloat,
		ast.BuiltinFloat,
		ast.BuiltinFloat32,
		ast.BuiltinFloat64:
		return true
	default:
		return false
	}
}

func isUnsigned(typ ast.Type) bool {
	switch typ {
	case
		ast.InferredUnsigned,
		ast.BuiltinUint,
		ast.BuiltinUint8,
		ast.BuiltinUint16,
		ast.BuiltinUint32,
		ast.BuiltinUint64:
		return true
	default:
		return false
	}
}

func isSigned(typ ast.Type) bool {
	switch typ {
	case
		ast.InferredSigned,
		ast.BuiltinInt,
		ast.BuiltinInt8,
		ast.BuiltinInt16,
		ast.BuiltinInt32,
		ast.BuiltinInt64:
		return true
	default:
		return false
	}
}
