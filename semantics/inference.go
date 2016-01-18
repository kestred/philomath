package semantics

// NOTE: Only performing inference; type checking is a separate step
// TODO: Break-out operator overload resolution
// TODO: Should literals continue to be parsed here, or elsewhere?

import (
	"strconv"
	"strings"

	"github.com/kestred/philomath/ast"
	"github.com/kestred/philomath/utils"
)

func InferTypes(expr ast.Expr) ast.Type {
	switch node := expr.(type) {
	case *ast.PostfixExpr:
		subtype := InferTypes(node.Subexpr)
		node.Type = inferPostfixType(node.Operator, subtype)
		return node.Type
	case *ast.InfixExpr:
		left := InferTypes(node.Left)
		right := InferTypes(node.Right)
		node.Type = inferInfixType(node.Operator, left, right)
		return node.Type
	case *ast.PrefixExpr:
		subtype := InferTypes(node.Subexpr)
		node.Type = inferPrefixType(node.Operator, subtype)
		return node.Type
	case *ast.GroupExpr:
		return InferTypes(node.Subexpr)
	case *ast.ValueExpr:
		switch literal := node.Literal.(type) {
		case *ast.NumberLiteral:
			var err error
			lit := literal.Literal
			if len(lit) > 2 && lit[0:2] == "0x" {
				node.Type = ast.InferredUnsigned
				literal.Value, err = strconv.ParseUint(lit, 16, 0)
				if err == strconv.ErrRange {
					panic("TODO: Handle hexadecimal literal can't be represented by uint64")
				}
				utils.AssertNil(err, "Failed parsing hexadecimal literal")
			} else if strings.Contains(lit, ".") || strings.Contains(lit, "e") {
				node.Type = ast.InferredReal
				literal.Value, err = strconv.ParseFloat(lit, 0)
				if err == strconv.ErrRange {
					panic("TODO: Handle floating point literal can't be represented by float64")
				}
				utils.AssertNil(err, "Failed parsing float literal")
			} else if lit[0] == '0' {
				node.Type = ast.InferredUnsigned
				literal.Value, err = strconv.ParseUint(lit, 8, 0)
				if err == strconv.ErrRange {
					panic("TODO: Handle octal literal can't be represented by uint64")
				}
				utils.AssertNil(err, "Failed parsing octal literal")
			} else {
				node.Type = ast.InferredNumber
				literal.Value, err = strconv.ParseUint(lit, 10, 0)
				if err == strconv.ErrRange {
					panic("TODO: Handle decimal literal can't be represented by uint64")
				}
				utils.AssertNil(err, "Failed parsing decimal literal")
			}
			return node.Type
		default:
			panic("TODO: Unhandled value literal")
		}
	default:
		panic("TODO: Handle type inferences for this expr")
	}
}

func inferPrefixType(op ast.Operator, typ ast.Type) ast.Type {
	if typ == ast.UnknownType {
		return typ
	}

	switch op.Literal {
	case "-", "+":
		switch typ {
		case ast.InferredNumber:
			return ast.InferredSigned
		case
			ast.InferredReal,
			ast.BuiltinReal,
			ast.BuiltinReal32,
			ast.BuiltinReal64,
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

func inferPostfixType(op ast.Operator, typ ast.Type) ast.Type {
	if typ == ast.UnknownType {
		return typ
	}

	switch op.Literal {
	default:
		panic("TODO: Unhandled postfix operator in type inference")
	}
}

func inferInfixType(op ast.Operator, left ast.Type, right ast.Type) ast.Type {
	if left == ast.UnknownType || right == ast.UnknownType {
		return ast.UnknownType
	}

	switch op.Literal {
	case "-", "+", "*", "/":
		// TODO: Implement operator overload resolution for prefix -/+
		return promoteNumber(left, right)
	default:
		panic("TODO: Unhandled infix operator in type inference")
	}
}

func promoteNumber(left ast.Type, right ast.Type) ast.Type {
	if !isNumber(left) {
		if !isNumber(right) {
			return ast.InferredNumber
		} else {
			return right
		}
	} else if !isNumber(right) {
		return ast.InferredNumber
	}

	if isReal(left) {
		if isReal(right) {
			return promoteByOrder(left, right)
		} else {
			return left
		}
	} else if isReal(right) {
		return right
	}

	if isSigned(left) {
		if isSigned(right) {
			return promoteByOrder(left, right)
		} else if right == ast.InferredNumber {
			return left
		} else {
			return ast.UnknownType
		}
	} else if isSigned(right) {
		if left == ast.InferredNumber {
			return right
		} else {
			return ast.UnknownType
		}
	}

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
	case ast.InferredReal, ast.InferredSigned, ast.InferredUnsigned:
		return 1
	case ast.BuiltinInt8, ast.BuiltinUint8:
		return 2
	case ast.BuiltinInt16, ast.BuiltinUint16:
		return 3
	case ast.BuiltinReal32, ast.BuiltinInt32, ast.BuiltinUint32:
		return 4
	case ast.BuiltinReal64, ast.BuiltinInt64, ast.BuiltinUint64:
		return 5
	case ast.BuiltinReal, ast.BuiltinInt, ast.BuiltinUint:
		return 6
	default:
		return -1
	}
}

func isNumber(typ ast.Type) bool {
	return (typ == ast.InferredNumber) || isReal(typ) || isSigned(typ) || isUnsigned(typ)
}

func isReal(typ ast.Type) bool {
	switch typ {
	case
		ast.InferredReal,
		ast.BuiltinReal,
		ast.BuiltinReal32,
		ast.BuiltinReal64:
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
