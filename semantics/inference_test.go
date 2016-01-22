package semantics

import (
	"testing"

	"github.com/kestred/philomath/ast"
	// TODO: Maybe avoid relying on parser when more code is stable?
	"github.com/kestred/philomath/parser"
	"github.com/stretchr/testify/assert"
)

func numberValue(t *testing.T, input string) interface{} {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	expr := p.ParseExpression()
	if len(p.Errors) > 0 {
		assert.Fail(t, "Unexpected parse error", p.Errors[0].Error())
	}

	InferTypes(expr)
	valExpr := expr.(*ast.ValueExpr)
	numLit := valExpr.Literal.(*ast.NumberLiteral)
	return numLit.Value
}

func TestLiteralValues(t *testing.T) {
	assert.Equal(t, uint64(22), numberValue(t, `22`))
	assert.Equal(t, uint64(0755), numberValue(t, `0755`))
	// TODO: Implement hexidecimal scanning
	// assert.Equal(t, uint64(0xff), numberValue(`0xff`))
	assert.Equal(t, float64(.32), numberValue(t, `.32`))
	assert.Equal(t, float64(3.2), numberValue(t, `3.2`))
	assert.Equal(t, float64(0.32), numberValue(t, `0.32`))
	assert.Equal(t, float64(3e2), numberValue(t, `3e2`))
	assert.Equal(t, float64(3e+2), numberValue(t, `3e+2`))
	assert.Equal(t, float64(3e-2), numberValue(t, `3e-2`))
}

func inferExpression(t *testing.T, input string) ast.Type {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	expr := p.ParseExpression()
	if len(p.Errors) > 0 {
		assert.Fail(t, "Unexpected parse error", p.Errors[0].Error())
	}

	typ := InferTypes(expr)
	assert.Equal(t, typ, expr.GetType())
	return typ
}

func inferBlock(t *testing.T, input string) *ast.Block {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	block := p.ParseBlock()
	if len(p.Errors) > 0 {
		assert.Fail(t, "Unexpected parse error", p.Errors[0].Error())
	}

	InferTypes(block)
	return block
}

func TestInferLiterals(t *testing.T) {
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `22`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `0755`))
	// TODO: Implement hexidecimal scanning
	// assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `0xff`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `.32`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `3.2`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `0.32`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `3e2`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `3e+2`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `3e-2`))
}

func TestInferArithmetic(t *testing.T) {
	// Prefix Operators
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `+7`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-07`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `+07`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `+7.0`))

	// Group Expressions
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `(7)`))
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `(07)`))
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `(-7)`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `(7.0)`))

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
	//  - combinations (num x float)
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 + 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 - 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 * 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 / 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 + 7`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 - 7`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 * 7`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 / 7`))
	//  - combinations (unsigned x float)
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 + 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 - 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 * 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 / 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 + 07`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 - 07`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 * 07`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 / 07`))
	//  - combinations (signed x float)
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 + -7`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 - -7`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 * -7`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 / -7`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 + 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 - 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 * 7.0`))
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 / 7.0`))

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

func TestInferBlock(t *testing.T) {
	block := inferBlock(t, `{
		hoge :: -3;         # constant decl
		hoge + 2;           # one ident in expr

		piyo := 0.5 * hoge; # mutable decl
		piyo / hoge;        # two ident in expr

		fuga := hogera;     # use undefined in decl
		0755 - fuga;        # propogate undefined in expr
	}`)

	if decl, ok := block.Nodes[0].(*ast.ConstantDecl); assert.True(t, ok) {
		if defn, ok := decl.Defn.(*ast.ExprDefn); assert.True(t, ok) {
			assert.Equal(t, ast.InferredSigned, defn.Expr.GetType())
		}
	}
	if stmt, ok := block.Nodes[1].(*ast.ExprStmt); assert.True(t, ok) {
		assert.Equal(t, ast.InferredSigned, stmt.Expr.GetType())
	}

	if decl, ok := block.Nodes[2].(*ast.MutableDecl); assert.True(t, ok) {
		assert.Equal(t, ast.InferredType, decl.Type) // not overridden for now
		assert.Equal(t, ast.InferredFloat, decl.Expr.GetType())
	}
	if stmt, ok := block.Nodes[3].(*ast.ExprStmt); assert.True(t, ok) {
		assert.Equal(t, ast.InferredFloat, stmt.Expr.GetType())
	}

	if decl, ok := block.Nodes[4].(*ast.MutableDecl); assert.True(t, ok) {
		assert.Equal(t, ast.InferredType, decl.Type)
		assert.Equal(t, ast.UnknownType, decl.Expr.GetType())
	}
	if stmt, ok := block.Nodes[5].(*ast.ExprStmt); assert.True(t, ok) {
		assert.Equal(t, ast.UnknownType, stmt.Expr.GetType())
	}
}
