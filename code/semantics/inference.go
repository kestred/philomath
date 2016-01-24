package semantics

import (
	"strconv"
	"strings"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/code"
	"github.com/kestred/philomath/code/utils"
)

// TODO: Break-out operator overload resolution

func InferTypes(cs *code.Section) {
	inferTypesRecursive(cs.Root)
}

func inferTypesRecursive(node ast.Node) ast.Type {
	switch n := node.(type) {
	case *ast.Block:
		for _, subnode := range n.Nodes {
			inferTypesRecursive(subnode)
		}
	case *ast.ConstantDecl:
		if defn, ok := n.Defn.(*ast.ConstantDefn); ok {
			inferTypesRecursive(defn.Expr)
		}
	case *ast.MutableDecl:
		typ := inferTypesRecursive(n.Expr)
		if n.Type == ast.InferredType {
			n.Type = typ
		}
	case *ast.EvalStmt:
		inferTypesRecursive(n.Expr)
	case *ast.AssignStmt:
		if len(n.Left) != len(n.Right) {
			panic("TODO: Handle unbalanced assignment")
		}
		for i := range n.Left {
			inferTypesRecursive(n.Left[i])
			inferTypesRecursive(n.Right[i])
		}
	case *ast.PostfixExpr:
		subtype := inferTypesRecursive(n.Subexpr)
		n.Type = inferPostfixType(n.Operator, subtype)
		return n.Type
	case *ast.InfixExpr:
		left := inferTypesRecursive(n.Left)
		right := inferTypesRecursive(n.Right)
		n.Type = inferInfixType(n.Operator, left, right)
		return n.Type
	case *ast.PrefixExpr:
		subtype := inferTypesRecursive(n.Subexpr)
		n.Type = inferPrefixType(n.Operator, subtype)
		return n.Type
	case *ast.GroupExpr:
		n.Type = inferTypesRecursive(n.Subexpr)
		return n.Type
	case *ast.Identifier:
		utils.Assert(n.Decl != nil, "An unresolved identifier survived until type inferrence")
		switch d := n.Decl.(type) {
		case *ast.ConstantDecl:
			n.Type = d.Defn.(*ast.ConstantDefn).Expr.GetType()
		case *ast.MutableDecl:
			n.Type = d.Type
		}
		return n.Type
	case *ast.NumberLiteral:
		n.Type, n.Value = parseNumber(n.Literal)
		return n.Type
	default:
		utils.InvalidCodePath()
	}
	return ast.BuiltinEmpty
}

func inferPrefixType(op *ast.OperatorDefn, typ ast.Type) ast.Type {
	if isError(typ) {
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
	if isError(typ) {
		return typ
	}

	switch op {
	default:
		panic("TODO: Unhandled postfix operator in type inference")
	}
}

func inferInfixType(op *ast.OperatorDefn, left ast.Type, right ast.Type) ast.Type {
	if isError(left) {
		return left
	} else if isError(right) {
		return right
	}

	switch op {
	case ast.BuiltinAdd, ast.BuiltinSubtract, ast.BuiltinMultiply, ast.BuiltinDivide:
		// TODO: Implement operator overload resolution
		return castNumbers(left, right)
	default:
		panic("TODO: Unhandled infix operator in type inference")
	}
}

// TODO: Should literals continue to be parsed here, or elsewhere?
func parseNumber(num string) (ast.Type, interface{}) {
	if len(num) > 2 && num[0:2] == "0x" {
		val, err := strconv.ParseUint(num, 16, 0)
		if err == strconv.ErrRange {
			panic("TODO: Handle hexadecimal literal can't be represented by uint64")
		}
		utils.AssertNil(err, "Failed parsing hexadecimal literal")
		return ast.InferredUnsigned, val
	} else if strings.Contains(num, ".") || strings.Contains(num, "e") {
		val, err := strconv.ParseFloat(num, 0)
		if err == strconv.ErrRange {
			panic("TODO: Handle floating point literal can't be represented by float64")
		}
		utils.AssertNil(err, "Failed parsing float literal")
		return ast.InferredFloat, val
	} else if num[0] == '0' {
		val, err := strconv.ParseUint(num, 8, 0)
		if err == strconv.ErrRange {
			panic("TODO: Handle octal literal can't be represented by uint64")
		}
		utils.AssertNil(err, "Failed parsing octal literal")
		return ast.InferredUnsigned, val
	} else {
		val, err := strconv.ParseUint(num, 10, 0)
		if err == strconv.ErrRange {
			panic("TODO: Handle decimal literal can't be represented by uint64")
		}
		utils.AssertNil(err, "Failed parsing decimal literal")
		return ast.InferredNumber, val
	}
}

func castNumbers(left ast.Type, right ast.Type) ast.Type {
	// NOTE: typechecking will come through later and assert that the implicit
	//       casts are either safe (or "safe-enough")
	if !maybeNumber(left) || !maybeNumber(right) {
		// can't cast non-number to number
		return ast.UncastableType
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
			return ast.UncastableType
		}
	} else if isSigned(right) {
		if left == ast.InferredNumber {
			// cast generic to signed
			return right
		} else {
			// can't cast unsigned to signed
			return ast.UncastableType
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

func isError(typ ast.Type) bool {
	switch typ {
	case
		ast.UninferredType,
		ast.UnresolvedType,
		ast.UncastableType:
		return true
	default:
		return false
	}
}

func maybeNumber(typ ast.Type) bool {
	return typ == ast.InferredNumber || typ == ast.InferredType ||
		isFloat(typ) || isSigned(typ) || isUnsigned(typ)
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
