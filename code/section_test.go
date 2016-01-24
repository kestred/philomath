package code

import (
	"reflect"
	"testing"

	"github.com/kestred/philomath/ast"
	"github.com/kestred/philomath/parser"
	"github.com/stretchr/testify/assert"
)

func parseBlock(t *testing.T, input string) *ast.Block {
	var p parser.Parser
	p.Init("example", false, []byte(input))
	block := p.ParseBlock()
	if len(p.Errors) > 0 {
		t.Fatalf("Unexpected parse error\n\n%v", p.Errors[0].Error())
	}
	return block
}

func TestFlattenBlock(t *testing.T) {
	expected := []ast.Node{
		&ast.Block{},
		// foo := -3;
		1: &ast.MutableDecl{},
		&ast.Ident{},
		&ast.PrefixExpr{},
		&ast.Operator{},
		&ast.ValueExpr{},
		&ast.NumberLiteral{},
		// baz :: 1;
		7: &ast.ConstantDecl{},
		&ast.Ident{},
		&ast.ExprDefn{},
		&ast.ValueExpr{},
		&ast.NumberLiteral{},
		// (2 + foo) + baz;
		12: &ast.ExprStmt{},
		&ast.InfixExpr{},
		&ast.GroupExpr{},
		&ast.InfixExpr{},
		&ast.ValueExpr{},
		&ast.NumberLiteral{},
		&ast.Operator{},
		&ast.ValueExpr{},
		&ast.Ident{},
		&ast.Operator{},
		&ast.ValueExpr{},
		&ast.Ident{},
		// {
		&ast.Block{},
		// bar := foo;
		25: &ast.MutableDecl{},
		&ast.Ident{},
		&ast.ValueExpr{},
		&ast.Ident{},
		// 0755 - baz;
		29: &ast.ExprStmt{},
		&ast.InfixExpr{},
		&ast.ValueExpr{},
		&ast.NumberLiteral{},
		&ast.Operator{},
		&ast.ValueExpr{},
		&ast.Ident{},
		// foo = baz * 4;
		36: &ast.AssignStmt{},
		&ast.ValueExpr{},
		&ast.Ident{},
		&ast.InfixExpr{},
		&ast.ValueExpr{},
		&ast.Ident{},
		&ast.Operator{},
		&ast.ValueExpr{},
		&ast.NumberLiteral{},
		// bar, foo = foo + 27, bar;
		45: &ast.AssignStmt{},
		&ast.ValueExpr{},
		&ast.Ident{},
		&ast.ValueExpr{},
		&ast.Ident{},
		&ast.InfixExpr{},
		&ast.ValueExpr{},
		&ast.Ident{},
		&ast.Operator{},
		&ast.ValueExpr{},
		&ast.NumberLiteral{},
		&ast.ValueExpr{},
		&ast.Ident{},
		// }
	}

	block := parseBlock(t, `{
		foo := -3;     # mutable declaration
		baz :: 1;      # constant declaration
		(2 + foo) + baz; # expression statement

		# a nested block
		{
			bar := foo;
			0755 - baz;

			foo = baz * 4;		        # assignment statement
			bar, foo = foo + 27, bar; # parallel assignment
		}
	}`)

	section := PrepareTree(block, nil)
	assert.Equal(t, block, section.Root)
	assert.Equal(t, block, section.Nodes[0])
	for i, example := range expected {
		if i >= len(section.Nodes) {
			t.Fatalf("Expected %v nodes, but only %v exist", len(expected), len(section.Nodes))
		} else {
			lhs := reflect.TypeOf(example).Elem().Name()
			rhs := reflect.TypeOf(section.Nodes[i]).Elem().Name()
			if lhs != rhs {
				t.Logf("Failed at %v", i)
				assert.Equal(t, lhs, rhs)
			}
		}
	}
}
