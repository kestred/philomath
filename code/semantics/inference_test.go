package semantics

import (
	"testing"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/code"
	"github.com/kestred/philomath/code/parser"
	"github.com/stretchr/testify/assert"
)

func inferAny(t *testing.T, input string) ast.Node {
	p := parser.Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	section := code.PrepareTree(node, nil)
	ResolveNames(&section)
	InferTypes(&section)
	return node
}

func TestInferDeclarations(t *testing.T) {
	block := inferAny(t, `{
		hoge :: -3;         # constant decl
		hoge + 2;           # one ident in expr

		piyo := 0.5 * hoge; # mutable decl
		piyo / hoge;        # two ident in expr
	}`).(*ast.Block)

	decl0 := block.Nodes[0].(*ast.ImmutableDecl)
	defn0 := decl0.Defn.(*ast.ConstantDefn)
	assert.Equal(t, ast.InferredSigned, defn0.Expr.GetType())

	stmt1 := block.Nodes[1].(*ast.EvalStmt)
	assert.Equal(t, ast.InferredSigned, stmt1.Expr.GetType())

	decl2 := block.Nodes[2].(*ast.MutableDecl)
	assert.Equal(t, ast.InferredFloat, decl2.Type)

	stmt3 := block.Nodes[3].(*ast.EvalStmt)
	assert.Equal(t, ast.InferredFloat, stmt3.Expr.GetType())
}

func TestInferAssignment(t *testing.T) {
	block := inferAny(t, `{
		plugh := -4;
		xyzzy := 012;

		plugh = 0.5 * plugh;
		xyzzy, plugh = (plugh / 5), xyzzy;
	}`).(*ast.Block)

	decl0 := block.Nodes[0].(*ast.MutableDecl)
	assert.Equal(t, ast.InferredSigned, decl0.Type)

	decl1 := block.Nodes[1].(*ast.MutableDecl)
	assert.Equal(t, ast.InferredUnsigned, decl1.Type)

	stmt2 := block.Nodes[2].(*ast.AssignStmt)
	if assert.Equal(t, 1, len(stmt2.Left)) && assert.Equal(t, 1, len(stmt2.Right)) {
		assert.Equal(t, ast.InferredSigned, stmt2.Left[0].GetType())
		assert.Equal(t, ast.InferredFloat, stmt2.Right[0].GetType())
	}

	stmt3 := block.Nodes[3].(*ast.AssignStmt)
	if assert.Equal(t, 2, len(stmt3.Left)) && assert.Equal(t, 2, len(stmt3.Right)) {
		assert.Equal(t, ast.InferredUnsigned, stmt3.Left[0].GetType())
		assert.Equal(t, ast.InferredSigned, stmt3.Left[1].GetType())
		assert.Equal(t, ast.InferredSigned, stmt3.Right[0].GetType())
		assert.Equal(t, ast.InferredUnsigned, stmt3.Right[1].GetType())
	}
}

func TestInferNestedBlock(t *testing.T) {
	block := inferAny(t, `{
		ham  := 0600;
		eggs :: -6.29;

		{
			spam := eggs / 2;
			spam - ham;
			eggs;
			ham;
		}
	}`).(*ast.Block)

	decl0 := block.Nodes[0].(*ast.MutableDecl)
	assert.Equal(t, ast.InferredUnsigned, decl0.Type)

	decl1 := block.Nodes[1].(*ast.ImmutableDecl)
	defn1 := decl1.Defn.(*ast.ConstantDefn)
	assert.Equal(t, ast.InferredFloat, defn1.Expr.GetType())

	if nest, ok := block.Nodes[2].(*ast.Block); assert.True(t, ok) {
		decl0 := nest.Nodes[0].(*ast.MutableDecl)
		assert.Equal(t, ast.InferredFloat, decl0.Type)

		stmt1 := nest.Nodes[1].(*ast.EvalStmt)
		assert.Equal(t, ast.InferredFloat, stmt1.Expr.GetType())

		stmt2 := nest.Nodes[2].(*ast.EvalStmt)
		assert.Equal(t, ast.InferredFloat, stmt2.Expr.GetType())

		stmt3 := nest.Nodes[3].(*ast.EvalStmt)
		assert.Equal(t, ast.InferredUnsigned, stmt3.Expr.GetType())
	}
}

func inferLiteral(t *testing.T, input string) ast.Literal {
	p := parser.Make("example", false, []byte(input+";"))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	section := code.PrepareTree(node, nil)
	InferTypes(&section)
	return node.(*ast.EvalStmt).Expr.(ast.Literal)
}

func TestInferLiterals(t *testing.T) {
	assert.Equal(t, uint64(22), inferLiteral(t, `22`).GetValue())
	assert.Equal(t, ast.InferredNumber, inferLiteral(t, `22`).GetType())
	assert.Equal(t, uint64(0755), inferLiteral(t, `0755`).GetValue())
	assert.Equal(t, ast.InferredUnsigned, inferLiteral(t, `0755`).GetType())
	// TODO: Implement hexidecimal scanning
	// assert.Equal(t, uint64(0xff), inferLiteral(`0xff`).GetValue())
	// assert.Equal(t, ast.InferredUnsigned, inferLiteral(t, `0xff`).GetType())
	assert.Equal(t, float64(.32), inferLiteral(t, `.32`).GetValue())
	assert.Equal(t, ast.InferredFloat, inferLiteral(t, `.32`).GetType())
	assert.Equal(t, float64(3.2), inferLiteral(t, `3.2`).GetValue())
	assert.Equal(t, ast.InferredFloat, inferLiteral(t, `3.2`).GetType())
	assert.Equal(t, float64(0.32), inferLiteral(t, `0.32`).GetValue())
	assert.Equal(t, ast.InferredFloat, inferLiteral(t, `0.32`).GetType())
	assert.Equal(t, float64(3e2), inferLiteral(t, `3e2`).GetValue())
	assert.Equal(t, ast.InferredFloat, inferLiteral(t, `3e2`).GetType())
	assert.Equal(t, float64(3e+2), inferLiteral(t, `3e+2`).GetValue())
	assert.Equal(t, ast.InferredFloat, inferLiteral(t, `3e+2`).GetType())
	assert.Equal(t, float64(3e-2), inferLiteral(t, `3e-2`).GetValue())
	assert.Equal(t, ast.InferredFloat, inferLiteral(t, `3e-2`).GetType())
}

func inferExpression(t *testing.T, input string) ast.Expr {
	p := parser.Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	section := code.PrepareTree(node, nil)
	InferTypes(&section)
	return node.(*ast.EvalStmt).Expr
}

func TestInferArithmetic(t *testing.T) {
	// Prefix Operators
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `+7;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-07;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `+07;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `+7.0;`).GetType())

	// Group Expressions
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `(7);`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `(07);`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `(-7);`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `(7.0);`).GetType())

	// Binary Operators
	//  - combinations (num x num)
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 + 7;`).GetType())
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 - 7;`).GetType())
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 * 7;`).GetType())
	assert.Equal(t, ast.InferredNumber, inferExpression(t, `7 / 7;`).GetType())
	//  - combinations (num x signed)
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 + 07;`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 - 07;`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 * 07;`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `7 / 07;`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 + 7;`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 - 7;`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 * 7;`).GetType())
	assert.Equal(t, ast.InferredUnsigned, inferExpression(t, `07 / 7;`).GetType())
	//  - combinations (num x signed)
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 + -7;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 - -7;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 * -7;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `7 / -7;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 + 7;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 - 7;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 * 7;`).GetType())
	assert.Equal(t, ast.InferredSigned, inferExpression(t, `-7 / 7;`).GetType())
	//  - combinations (unsigned x signed)
	assert.Equal(t, ast.UncastableType, inferExpression(t, `07 + -7;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `07 - -7;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `07 * -7;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `07 / -7;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-7 + 07;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-7 - 07;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-7 * 07;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-7 / 07;`).GetType())
	//  - combinations (num x float)
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 + 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 - 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 * 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7 / 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 + 7;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 - 7;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 * 7;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 / 7;`).GetType())
	//  - combinations (unsigned x float)
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 + 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 - 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 * 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `07 / 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 + 07;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 - 07;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 * 07;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 / 07;`).GetType())
	//  - combinations (signed x float)
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 + -7;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 - -7;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 * -7;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `7.0 / -7;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 + 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 - 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 * 7.0;`).GetType())
	assert.Equal(t, ast.InferredFloat, inferExpression(t, `-7 / 7.0;`).GetType())

	// Propogate unknown Type
	assert.Equal(t, ast.UncastableType, inferExpression(t, `+(-7 + 07);`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `-(-7 + 07);`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `7 + (-7 + 07);`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `7 - (-7 + 07);`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `7 * (-7 + 07);`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `7 / (-7 + 07);`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `(-7 + 07) + 7;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `(-7 + 07) - 7;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `(-7 + 07) * 7;`).GetType())
	assert.Equal(t, ast.UncastableType, inferExpression(t, `(-7 + 07) / 7;`).GetType())
}
