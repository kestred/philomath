package semantics

import (
	"testing"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/code"
	"github.com/kestred/philomath/code/parser"
	"github.com/stretchr/testify/assert"
)

func inferExpression(t *testing.T, input string) ast.Expr {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	expr := p.ParseExpression()
	if len(p.Errors) > 0 {
		assert.Fail(t, "Unexpected parse error", p.Errors[0].Error())
	}

	section := code.PrepareTree(expr, nil)
	InferTypes(&section)
	return expr
}

func inferBlock(t *testing.T, input string) *ast.Block {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	block := p.ParseBlock()
	if len(p.Errors) > 0 {
		assert.Fail(t, "Unexpected parse error", p.Errors[0].Error())
	}

	section := code.PrepareTree(block, nil)
	ResolveNames(&section)
	InferTypes(&section)
	return block
}

func TestLiteralValues(t *testing.T) {
	assert.Equal(t, uint64(22), inferExpression(t, `22`).(*ast.NumberLiteral).Value)
	assert.Equal(t, uint64(0755), inferExpression(t, `0755`).(*ast.NumberLiteral).Value)
	// TODO: Implement hexidecimal scanning
	// assert.Equal(t, uint64(0xff), inferExpression(`0xff`).(*ast.NumberLiteral).Value)
	assert.Equal(t, float64(.32), inferExpression(t, `.32`).(*ast.NumberLiteral).Value)
	assert.Equal(t, float64(3.2), inferExpression(t, `3.2`).(*ast.NumberLiteral).Value)
	assert.Equal(t, float64(0.32), inferExpression(t, `0.32`).(*ast.NumberLiteral).Value)
	assert.Equal(t, float64(3e2), inferExpression(t, `3e2`).(*ast.NumberLiteral).Value)
	assert.Equal(t, float64(3e+2), inferExpression(t, `3e+2`).(*ast.NumberLiteral).Value)
	assert.Equal(t, float64(3e-2), inferExpression(t, `3e-2`).(*ast.NumberLiteral).Value)
}

func TestInferLiterals(t *testing.T) {
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `22`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `0755`).GetType())
	// TODO: Implement hexidecimal scanning
	// assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `0xff`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `.32`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `3.2`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `0.32`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `3e2`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `3e+2`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `3e-2`).GetType())
}

func TestInferArithmetic(t *testing.T) {
	// Prefix Operators
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `+7`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-07`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `+07`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `+7.0`).GetType())

	// Group Expressions
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `(7)`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `(07)`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `(-7)`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `(7.0)`).GetType())

	// Binary Operators
	//  - combinations (num x num)
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 + 7`).GetType())
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 - 7`).GetType())
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 * 7`).GetType())
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 / 7`).GetType())
	//  - combinations (num x signed)
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 + 07`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 - 07`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 * 07`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 / 07`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 + 7`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 - 7`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 * 7`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 / 7`).GetType())
	//  - combinations (num x signed)
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 + -7`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 - -7`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 * -7`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 / -7`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 + 7`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 - 7`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 * 7`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 / 7`).GetType())
	//  - combinations (unsigned x signed)
	assert.Equal(t, ast.UncastableType, inferExpression(t, `07 + -7`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `07 - -7`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `07 * -7`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `07 / -7`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-7 + 07`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-7 - 07`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-7 * 07`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-7 / 07`).GetType())
	//  - combinations (num x float)
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 + 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 - 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 * 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 / 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 + 7`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 - 7`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 * 7`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 / 7`).GetType())
	//  - combinations (unsigned x float)
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 + 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 - 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 * 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 / 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 + 07`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 - 07`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 * 07`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 / 07`).GetType())
	//  - combinations (signed x float)
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 + -7`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 - -7`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 * -7`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 / -7`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 + 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 - 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 * 7.0`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 / 7.0`).GetType())

	// Propogate unknown Type
	assert.Equal(t, ast.UncastableType, inferExpression(t, `+(-7 + 07)`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-(-7 + 07)`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `7 + (-7 + 07)`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `7 - (-7 + 07)`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `7 * (-7 + 07)`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `7 / (-7 + 07)`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `(-7 + 07) + 7`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `(-7 + 07) - 7`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `(-7 + 07) * 7`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `(-7 + 07) / 7`).GetType())
}

func TestInferDeclarations(t *testing.T) {
	block := inferBlock(t, `{
		hoge :: -3;         # constant decl
		hoge + 2;           # one ident in expr

		piyo := 0.5 * hoge; # mutable decl
		piyo / hoge;        # two ident in expr
	}`)
	/*
			fuga := hogera;     # use undefined in decl
			0755 - fuga;        # propogate undefined in expr
		}`)
	*/
	if decl, ok := block.Nodes[0].(*ast.ImmutableDecl); assert.True(t, ok) {
		if defn, ok := decl.Defn.(*ast.ConstantDefn); assert.True(t, ok) {
			assert.Equal(t, ast.InferredSigned, defn.Expr.GetType())
		}
	}
	if stmt, ok := block.Nodes[1].(*ast.EvalStmt); assert.True(t, ok) {
		assert.Equal(t, ast.InferredSigned, stmt.Expr.GetType())
	}
	if decl, ok := block.Nodes[2].(*ast.MutableDecl); assert.True(t, ok) {
		assert.Equal(t, ast.InferredFloat, decl.Type)
	}
	if stmt, ok := block.Nodes[3].(*ast.EvalStmt); assert.True(t, ok) {
		assert.Equal(t, ast.InferredFloat, stmt.Expr.GetType())
	}
	/*
		if decl, ok := block.Nodes[4].(*ast.MutableDecl); assert.True(t, ok) {
			assert.Equal(t, ast.UncastableType, decl.Typ)
		}
		if stmt, ok := block.Nodes[5].(*ast.EvalStmt); assert.True(t, ok) {
			assert.Equal(t, ast.UncastableType, stmt.Expr.GetType())
		}
	*/
}

func TestInferAssignment(t *testing.T) {
	block := inferBlock(t, `{
		plugh := -4;
		xyzzy := 012;

		plugh = 0.5 * plugh;
		xyzzy, plugh = (plugh / 5), xyzzy;
	}`)
	if decl, ok := block.Nodes[0].(*ast.MutableDecl); assert.True(t, ok) {
		assert.Equal(t, ast.InferredSigned, decl.Type)
	}
	if decl, ok := block.Nodes[1].(*ast.MutableDecl); assert.True(t, ok) {
		assert.Equal(t, ast.InferredUnsigned, decl.Type)
	}
	if stmt, ok := block.Nodes[2].(*ast.AssignStmt); assert.True(t, ok) {
		if assert.Equal(t, 1, len(stmt.Left)) && assert.Equal(t, 1, len(stmt.Right)) {
			assert.Equal(t, ast.InferredSigned, stmt.Left[0].GetType())
			assert.Equal(t, ast.InferredFloat, stmt.Right[0].GetType())
		}
	}
	if stmt, ok := block.Nodes[3].(*ast.AssignStmt); assert.True(t, ok) {
		if assert.Equal(t, 2, len(stmt.Left)) && assert.Equal(t, 2, len(stmt.Right)) {
			assert.Equal(t, ast.InferredUnsigned, stmt.Left[0].GetType())
			assert.Equal(t, ast.InferredSigned, stmt.Left[1].GetType())
			assert.Equal(t, ast.InferredSigned, stmt.Right[0].GetType())
			assert.Equal(t, ast.InferredUnsigned, stmt.Right[1].GetType())
		}
	}
}

func TestInferNestedBlock(t *testing.T) {
	block := inferBlock(t, `{
		ham  := 0600;
		eggs :: -6.29;

		{
			spam := eggs / 2;
			spam - ham;
			eggs;
			ham;
		}
	}`)
	/*
			eggs * spam;
		}`)
	*/
	if decl, ok := block.Nodes[0].(*ast.MutableDecl); assert.True(t, ok) {
		assert.Equal(t, ast.InferredUnsigned, decl.Type)
	}
	if decl, ok := block.Nodes[1].(*ast.ImmutableDecl); assert.True(t, ok) {
		if defn, ok := decl.Defn.(*ast.ConstantDefn); assert.True(t, ok) {
			assert.Equal(t, ast.InferredFloat, defn.Expr.GetType())
		}
	}
	if nest, ok := block.Nodes[2].(*ast.Block); assert.True(t, ok) {
		if decl, ok := nest.Nodes[0].(*ast.MutableDecl); assert.True(t, ok) {
			assert.Equal(t, ast.InferredFloat, decl.Type)
		}
		if stmt, ok := nest.Nodes[1].(*ast.EvalStmt); assert.True(t, ok) {
			assert.Equal(t, ast.InferredFloat, stmt.Expr.GetType())
		}
		if stmt, ok := nest.Nodes[2].(*ast.EvalStmt); assert.True(t, ok) {
			assert.Equal(t, ast.InferredFloat, stmt.Expr.GetType())
		}
		if stmt, ok := nest.Nodes[3].(*ast.EvalStmt); assert.True(t, ok) {
			assert.Equal(t, ast.InferredUnsigned, stmt.Expr.GetType())
		}
	}
	/*
		if stmt, ok := block.Nodes[3].(*ast.EvalStmt); assert.True(t, ok) {
			assert.Equal(t, ast.UncastableType, stmt.Expr.GetType())
		}
	*/
}
