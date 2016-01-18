package parser

import (
	"testing"

	"github.com/kestred/philomath/ast"
	"github.com/stretchr/testify/assert"
)

func parseExpression(input string) ast.Expr {
	var parser Parser
	parser.Init("test", false, []byte(input))
	return parser.ParseExpression()
}

func TestParseError(t *testing.T) {
	var parser Parser
	parser.Init("error.phi", false, []byte(`1 * (2 + 3} - 4`))
	parser.ParseExpression()
	if assert.True(t, len(parser.Errors) > 0, "Expected some errors but found none.") {
		assert.Equal(t, "error.phi:1:12: Expected '}' but recieved ')'.", parser.Errors[0].Error())
	}
}

func TestParseArithmetic(t *testing.T) {
	var expected ast.Expr

	// add follows multiply
	expected = ast.InExp(
		ast.InExp(
			ast.ValExp(ast.NumLit("2")),
			ast.Operator{"*"},
			ast.ValExp(ast.NumLit("3")),
		),
		ast.Operator{"+"},
		ast.ValExp(ast.NumLit("4")),
	)

	assert.Equal(t, expected, parseExpression(`2 * 3 + 4`))

	// multiply follows add
	expected = ast.InExp(
		ast.ValExp(ast.NumLit("2")),
		ast.Operator{"+"},
		ast.InExp(
			ast.ValExp(ast.NumLit("3")),
			ast.Operator{"*"},
			ast.ValExp(ast.NumLit("4")),
		),
	)

	assert.Equal(t, expected, parseExpression(`2 + 3 * 4`))

	// multiply follows grouped add
	expected = ast.InExp(
		ast.GrpExp(ast.InExp(
			ast.ValExp(ast.NumLit("2")),
			ast.Operator{"+"},
			ast.ValExp(ast.NumLit("3")),
		)),
		ast.Operator{"*"},
		ast.ValExp(ast.NumLit("4")),
	)

	assert.Equal(t, expected, parseExpression(`(2 + 3) * 4`))

	// add and subtract associativity
	expected = ast.InExp(
		ast.InExp(
			ast.InExp(
				ast.InExp(
					ast.InExp(
						ast.InExp(
							ast.ValExp(ast.NumLit("2")),
							ast.Operator{"+"},
							ast.ValExp(ast.NumLit("3")),
						),
						ast.Operator{"+"},
						ast.ValExp(ast.NumLit("4")),
					),
					ast.Operator{"-"},
					ast.ValExp(ast.NumLit("5")),
				),
				ast.Operator{"+"},
				ast.ValExp(ast.NumLit("6")),
			),
			ast.Operator{"-"},
			ast.ValExp(ast.NumLit("7")),
		),
		ast.Operator{"-"},
		ast.ValExp(ast.NumLit("8")),
	)

	assert.Equal(t, expected, parseExpression(`2 + 3 + 4 - 5 + 6 - 7 - 8`))

	// multiply, divide, and modulus associativity
	expected = ast.InExp(
		ast.InExp(
			ast.InExp(
				ast.InExp(
					ast.InExp(
						ast.InExp(
							ast.ValExp(ast.NumLit("2")),
							ast.Operator{"/"},
							ast.ValExp(ast.NumLit("3")),
						),
						ast.Operator{"/"},
						ast.ValExp(ast.NumLit("4")),
					),
					ast.Operator{"*"},
					ast.ValExp(ast.NumLit("5")),
				),
				ast.Operator{"*"},
				ast.ValExp(ast.NumLit("6")),
			),
			ast.Operator{"%"},
			ast.ValExp(ast.NumLit("7")),
		),
		ast.Operator{"/"},
		ast.ValExp(ast.NumLit("8")),
	)

	assert.Equal(t, expected, parseExpression(`2 / 3 / 4 * 5 * 6 % 7 / 8`))

	// signed addition
	expected = ast.InExp(
		ast.PreExp(ast.Operator{"-"}, ast.ValExp(ast.NumLit("2"))),
		ast.Operator{"+"},
		ast.PreExp(ast.Operator{"+"}, ast.ValExp(ast.NumLit("4"))),
	)

	assert.Equal(t, expected, parseExpression(`-2 + +4`))

	// signed subtraction
	expected = ast.InExp(
		ast.PreExp(ast.Operator{"-"}, ast.ValExp(ast.NumLit("2"))),
		ast.Operator{"-"},
		ast.PreExp(ast.Operator{"+"}, ast.ValExp(ast.NumLit("4"))),
	)

	assert.Equal(t, expected, parseExpression(`-2 - +4`))

	// signed multiplication
	expected = ast.InExp(
		ast.PreExp(ast.Operator{"-"}, ast.ValExp(ast.NumLit("2"))),
		ast.Operator{"*"},
		ast.PreExp(ast.Operator{"+"}, ast.ValExp(ast.NumLit("4"))),
	)

	assert.Equal(t, expected, parseExpression(`-2 * +4`))

	// signed division
	expected = ast.InExp(
		ast.PreExp(ast.Operator{"-"}, ast.ValExp(ast.NumLit("2"))),
		ast.Operator{"/"},
		ast.PreExp(ast.Operator{"+"}, ast.ValExp(ast.NumLit("4"))),
	)

	assert.Equal(t, expected, parseExpression(`-2 / +4`))
}
