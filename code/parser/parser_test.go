package parser

import (
	"testing"

	"github.com/kestred/philomath/code/ast"
	"github.com/stretchr/testify/assert"
)

func parseExpression(t *testing.T, input string) ast.Expr {
	var p Parser
	p.Init("example", false, []byte(input))
	expr := p.ParseExpression()
	if len(p.Errors) > 0 {
		t.Fatalf("Unexpected parse error\n\n%v", p.Errors[0].Error())
	}
	return expr
}

func parseBlock(t *testing.T, input string) *ast.Block {
	var p Parser
	p.Init("example", false, []byte(input))
	block := p.ParseBlock()
	if len(p.Errors) > 0 {
		t.Fatalf("Unexpected parse error\n\n%v", p.Errors[0].Error())
	}
	return block
}

func TestParseError(t *testing.T) {
	var parser Parser
	parser.Init("error.phi", false, []byte(`1 * (2 + 3} - 4`))
	parser.ParseExpression()
	if assert.True(t, len(parser.Errors) > 0, "Expected some errors but found none.") {
		assert.Equal(t, "error.phi:1:12: Expected ')' but recieved '}'.", parser.Errors[0].Error())
	}

	parser = Parser{}
	parser.Init("error.phi", false, []byte(`{ 1 - 4 }`))
	parser.ParseBlock()
	if assert.True(t, len(parser.Errors) > 0, "Expected some errors but found none.") {
		assert.Equal(t, "error.phi:1:10: Expected ';' but recieved '}'.", parser.Errors[0].Error())
	}
}

func TestParseArithmetic(t *testing.T) {
	var expected ast.Expr

	// add follows multiply
	expected = ast.InExp(
		ast.InExp(
			ast.ValExp(ast.NumLit("2")),
			ast.BuiltinMultiply,
			ast.ValExp(ast.NumLit("3")),
		),
		ast.BuiltinAdd,
		ast.ValExp(ast.NumLit("4")),
	)

	assert.Equal(t, expected, parseExpression(t, `2 * 3 + 4`))

	// multiply follows add
	expected = ast.InExp(
		ast.ValExp(ast.NumLit("2")),
		ast.BuiltinAdd,
		ast.InExp(
			ast.ValExp(ast.NumLit("3")),
			ast.BuiltinMultiply,
			ast.ValExp(ast.NumLit("4")),
		),
	)

	assert.Equal(t, expected, parseExpression(t, `2 + 3 * 4`))

	// multiply follows grouped add
	expected = ast.InExp(
		ast.GrpExp(ast.InExp(
			ast.ValExp(ast.NumLit("2")),
			ast.BuiltinAdd,
			ast.ValExp(ast.NumLit("3")),
		)),
		ast.BuiltinMultiply,
		ast.ValExp(ast.NumLit("4")),
	)

	assert.Equal(t, expected, parseExpression(t, `(2 + 3) * 4`))

	// add and subtract associativity
	expected = ast.InExp(
		ast.InExp(
			ast.InExp(
				ast.InExp(
					ast.InExp(
						ast.InExp(
							ast.ValExp(ast.NumLit("2")),
							ast.BuiltinAdd,
							ast.ValExp(ast.NumLit("3")),
						),
						ast.BuiltinAdd,
						ast.ValExp(ast.NumLit("4")),
					),
					ast.BuiltinSubtract,
					ast.ValExp(ast.NumLit("5")),
				),
				ast.BuiltinAdd,
				ast.ValExp(ast.NumLit("6")),
			),
			ast.BuiltinSubtract,
			ast.ValExp(ast.NumLit("7")),
		),
		ast.BuiltinSubtract,
		ast.ValExp(ast.NumLit("8")),
	)

	assert.Equal(t, expected, parseExpression(t, `2 + 3 + 4 - 5 + 6 - 7 - 8`))

	// multiply, divide, and modulus associativity
	expected = ast.InExp(
		ast.InExp(
			ast.InExp(
				ast.InExp(
					ast.InExp(
						ast.InExp(
							ast.ValExp(ast.NumLit("2")),
							ast.BuiltinDivide,
							ast.ValExp(ast.NumLit("3")),
						),
						ast.BuiltinDivide,
						ast.ValExp(ast.NumLit("4")),
					),
					ast.BuiltinMultiply,
					ast.ValExp(ast.NumLit("5")),
				),
				ast.BuiltinMultiply,
				ast.ValExp(ast.NumLit("6")),
			),
			ast.BuiltinRemainder,
			ast.ValExp(ast.NumLit("7")),
		),
		ast.BuiltinDivide,
		ast.ValExp(ast.NumLit("8")),
	)

	assert.Equal(t, expected, parseExpression(t, `2 / 3 / 4 * 5 * 6 % 7 / 8`))

	// signed addition
	expected = ast.InExp(
		ast.PreExp(ast.BuiltinNegative, ast.ValExp(ast.NumLit("2"))),
		ast.BuiltinAdd,
		ast.PreExp(ast.BuiltinPositive, ast.ValExp(ast.NumLit("4"))),
	)

	assert.Equal(t, expected, parseExpression(t, `-2 + +4`))

	// signed subtraction
	expected = ast.InExp(
		ast.PreExp(ast.BuiltinNegative, ast.ValExp(ast.NumLit("2"))),
		ast.BuiltinSubtract,
		ast.PreExp(ast.BuiltinPositive, ast.ValExp(ast.NumLit("4"))),
	)

	assert.Equal(t, expected, parseExpression(t, `-2 - +4`))

	// signed multiplication
	expected = ast.InExp(
		ast.PreExp(ast.BuiltinNegative, ast.ValExp(ast.NumLit("2"))),
		ast.BuiltinMultiply,
		ast.PreExp(ast.BuiltinPositive, ast.ValExp(ast.NumLit("4"))),
	)

	assert.Equal(t, expected, parseExpression(t, `-2 * +4`))

	// signed division
	expected = ast.InExp(
		ast.PreExp(ast.BuiltinNegative, ast.ValExp(ast.NumLit("2"))),
		ast.BuiltinDivide,
		ast.PreExp(ast.BuiltinPositive, ast.ValExp(ast.NumLit("4"))),
	)

	assert.Equal(t, expected, parseExpression(t, `-2 / +4`))
}

func TestParseBlock(t *testing.T) {
	expected := &ast.Block{[]ast.Blockable{
		ast.Mutable("foo", nil, ast.ValExp(ast.NumLit("3"))),
		ast.Constant("baz", &ast.ExprDefn{ast.ValExp(ast.NumLit("1"))}),
		&ast.ExprStmt{ast.InExp(
			ast.InExp(
				ast.ValExp(ast.NumLit("2")),
				ast.BuiltinAdd,
				ast.ValExp(ast.Ident("foo")),
			),
			ast.BuiltinAdd,
			ast.ValExp(ast.Ident("baz")),
		)},
		&ast.Block{[]ast.Blockable{
			ast.Mutable("bar", nil, ast.ValExp(ast.Ident("foo"))),
			&ast.ExprStmt{ast.InExp(
				ast.ValExp(ast.NumLit("0755")),
				ast.BuiltinSubtract,
				ast.ValExp(ast.Ident("baz")),
			)},
			&ast.AssignStmt{
				[]ast.Expr{ast.ValExp(ast.Ident("foo"))},
				nil,
				[]ast.Expr{ast.InExp(
					ast.ValExp(ast.Ident("baz")),
					ast.BuiltinMultiply,
					ast.ValExp(ast.NumLit("4")),
				)},
			},
			&ast.AssignStmt{
				[]ast.Expr{
					ast.ValExp(ast.Ident("bar")),
					ast.ValExp(ast.Ident("foo")),
				},
				nil,
				[]ast.Expr{
					ast.InExp(
						ast.ValExp(ast.Ident("foo")),
						ast.BuiltinAdd,
						ast.ValExp(ast.NumLit("27")),
					),
					ast.ValExp(ast.Ident("bar")),
				},
			},
		}},
		&ast.ExprStmt{ast.InExp(
			ast.ValExp(ast.NumLit("8.4e-5")),
			ast.BuiltinDivide,
			ast.ValExp(ast.NumLit("0.5")),
		)},
		&ast.Block{nil},
	}}

	assert.Equal(t, expected, parseBlock(t, `{
		foo := 3;      # mutable declaration
		baz :: 1;      # constant declaration
		2 + foo + baz; # expression statement

		# a nested block
		{
			bar := foo;
			0755 - baz;

			foo = baz * 4;		        # assignment statement
			bar, foo = foo + 27, bar; # parallel assignment
		}

		# ignore extra semicolons
		; ;
		8.4e-5 / 0.5;; ;
		;

		# empty block
		{
			; # ignore extra semicolons occuring before a statement
		}
	}`))
}
