package semantics

import (
	"testing"

	"github.com/kestred/philomath/ast"
	// TODO: Maybe avoid relying on parser when more code is stable?
	"github.com/kestred/philomath/parser"
	"github.com/stretchr/testify/assert"
)

func numberValue(input string) interface{} {
	var p parser.Parser
	p.Init("test", false, []byte(input))
	expr := p.ParseExpression()
	InferTypes(expr)

	valExpr := expr.(*ast.ValueExpr)
	numLit := valExpr.Literal.(*ast.NumberLiteral)
	return numLit.Value
}

func TestLiteralValues(t *testing.T) {
	assert.Equal(t, uint64(22), numberValue(`22`))
	assert.Equal(t, uint64(0755), numberValue(`0755`))
	// TODO: Implement hexidecimal scanning
	// assert.Equal(t, uint64(0xff), numberValue(`0xff`))
	assert.Equal(t, float64(.32), numberValue(`.32`))
	assert.Equal(t, float64(3.2), numberValue(`3.2`))
	assert.Equal(t, float64(0.32), numberValue(`0.32`))
	assert.Equal(t, float64(3e2), numberValue(`3e2`))
	assert.Equal(t, float64(3e+2), numberValue(`3e+2`))
	assert.Equal(t, float64(3e-2), numberValue(`3e-2`))
}

func inferExpression(t *testing.T, input string) ast.Type {
	var p parser.Parser
	p.Init("test", false, []byte(input))
	expr := p.ParseExpression()
	typ := InferTypes(expr)
	assert.Equal(t, typ, expr.GetType())
	return typ
}

func TestInferLiterals(t *testing.T) {
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `22`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `0755`))
	// TODO: Implement hexidecimal scanning
	// assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `0xff`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `.32`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `3.2`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `0.32`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `3e2`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `3e+2`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `3e-2`))
}

func TestInferArithmetic(t *testing.T) {
	// Prefix Operators
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `+7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-07`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `+07`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `-7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `+7.0`))

	// Group Expressions
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `(7)`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `(07)`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `(-7)`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `(7.0)`))

	// Binary Operators
	//  - combinations (num x num)
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 + 7`))
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 - 7`))
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 * 7`))
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 / 7`))
	//  - combinations (num x signed)
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 + 07`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 - 07`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 * 07`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 / 07`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 + 7`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 - 7`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 * 7`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 / 7`))
	//  - combinations (num x signed)
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 + -7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 - -7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 * -7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 / -7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 + 7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 - 7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 * 7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 / 7`))
	//  - combinations (unsigned x signed)
	assert.Equal(t, ast.UnknownType, inferExpression(t, `07 + -7`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `07 - -7`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `07 * -7`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `07 / -7`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `-7 + 07`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `-7 - 07`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `-7 * 07`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `-7 / 07`))
	//  - combinations (num x real)
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7 + 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7 - 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7 * 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7 / 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 + 7`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 - 7`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 * 7`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 / 7`))
	//  - combinations (unsigned x real)
	assert.Equal(t, ast.InferredReal, inferExpression(t, `07 + 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `07 - 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `07 * 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `07 / 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 + 07`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 - 07`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 * 07`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 / 07`))
	//  - combinations (signed x real)
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 + -7`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 - -7`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 * -7`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `7.0 / -7`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `-7 + 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `-7 - 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `-7 * 7.0`))
	assert.Equal(t, ast.InferredReal, inferExpression(t, `-7 / 7.0`))

	// Propogate unknown Type
	assert.Equal(t, ast.UnknownType, inferExpression(t, `+(-7 + 07)`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `-(-7 + 07)`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `7 + (-7 + 07)`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `7 - (-7 + 07)`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `7 * (-7 + 07)`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `7 / (-7 + 07)`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `(-7 + 07) + 7`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `(-7 + 07) - 7`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `(-7 + 07) * 7`))
	assert.Equal(t, ast.UnknownType, inferExpression(t, `(-7 + 07) / 7`))
}
