package code

import (
	"reflect"
	"testing"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/parser"
	"github.com/stretchr/testify/assert"
)

func parseExample(t *testing.T, input string) ast.Node {
	p := parser.Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	return node
}

func TestFlattenBlock(t *testing.T) {
	expected := []ast.Node{
		&ast.Block{},
		// foo := -3;
		&ast.MutableDecl{},
		&ast.Identifier{},
		&ast.BaseType{},
		&ast.PrefixExpr{},
		&ast.OperatorDefn{},
		&ast.NumberLiteral{},
		// baz :: 1;
		&ast.ImmutableDecl{},
		&ast.Identifier{},
		&ast.ConstantDefn{},
		&ast.NumberLiteral{},
		// (2 + foo) + baz;
		&ast.EvalStmt{},
		&ast.InfixExpr{},
		&ast.GroupExpr{},
		&ast.InfixExpr{},
		&ast.NumberLiteral{},
		&ast.OperatorDefn{},
		&ast.Identifier{},
		&ast.OperatorDefn{},
		&ast.Identifier{},
		// {
		&ast.Block{},
		// bar := foo;
		&ast.MutableDecl{},
		&ast.Identifier{},
		&ast.BaseType{},
		&ast.Identifier{},
		// 0755 - baz;
		&ast.EvalStmt{},
		&ast.InfixExpr{},
		&ast.NumberLiteral{},
		&ast.OperatorDefn{},
		&ast.Identifier{},
		// foo = baz * 4;
		&ast.AssignStmt{},
		&ast.Identifier{},
		&ast.InfixExpr{},
		&ast.Identifier{},
		&ast.OperatorDefn{},
		&ast.NumberLiteral{},
		// bar, foo = foo + 27, bar;
		&ast.AssignStmt{},
		&ast.Identifier{},
		&ast.Identifier{},
		&ast.InfixExpr{},
		&ast.Identifier{},
		&ast.OperatorDefn{},
		&ast.NumberLiteral{},
		&ast.Identifier{},
		// }
	}

	block := parseExample(t, `{
		foo := -3;     # mutable declaration
		baz :: 1;      # constant definition
		(2 + foo) + baz; # evaluated statement

		# a nested block
		{
			bar := foo;
			0755 - baz;

			foo = baz * 4;		        # assignment
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
