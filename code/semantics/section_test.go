package semantics

/*
import (
	"reflect"
	"testing"

	"github.com/kestred/philomath/code/parser"
	"github.com/stretchr/testify/assert"
)

func parseExample(t *testing.T, input string) Node {
	p := parser.Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	return node
}

func TestFlattenBlock(t *testing.T) {
	expected := []Node{
		&Block{},
		// foo := -3;
		&MutableDecl{},
		&Identifier{},
		&BaseType{},
		&PrefixExpr{},
		&OperatorDefn{},
		&NumberLiteral{},
		// baz :: 1;
		&ImmutableDecl{},
		&Identifier{},
		&ConstantDefn{},
		&NumberLiteral{},
		// (2 + foo) + baz;
		&EvalStmt{},
		&InfixExpr{},
		&GroupExpr{},
		&InfixExpr{},
		&NumberLiteral{},
		&OperatorDefn{},
		&Identifier{},
		&OperatorDefn{},
		&Identifier{},
		// {
		&Block{},
		// bar := foo;
		&MutableDecl{},
		&Identifier{},
		&BaseType{},
		&Identifier{},
		// 0755 - baz;
		&EvalStmt{},
		&InfixExpr{},
		&NumberLiteral{},
		&OperatorDefn{},
		&Identifier{},
		// foo = baz * 4;
		&AssignStmt{},
		&Identifier{},
		&InfixExpr{},
		&Identifier{},
		&OperatorDefn{},
		&NumberLiteral{},
		// bar, foo = foo + 27, bar;
		&AssignStmt{},
		&Identifier{},
		&Identifier{},
		&InfixExpr{},
		&Identifier{},
		&OperatorDefn{},
		&NumberLiteral{},
		&Identifier{},
		// }
	}

	block := parseExample(t, `{
		foo := -3;     // mutable declaration
		baz :: 1;      // constant definition
		(2 + foo) + baz; // evaluated statement

		// a nested block
		{
			bar := foo;
			0755 - baz;

			foo = baz * 4;		        // assignment
			bar, foo = foo + 27, bar; // parallel assignment
		}
	}`)

	section := FlattenTree(block, nil)
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
*/
