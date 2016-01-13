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

func TestParseArithmetic(t *testing.T) {
	var expected ast.Expr

	// add follows multiply
	expected = ast.NewInfixExpr(
		ast.NewInfixExpr(
			ast.NewValueExpr(ast.NewNumberLiteral("2")),
			ast.Operator{"*"},
			ast.NewValueExpr(ast.NewNumberLiteral("3")),
		),
		ast.Operator{"+"},
		ast.NewValueExpr(ast.NewNumberLiteral("4")),
	)

	assert.Equal(t, expected, parseExpression(`2 * 3 + 4`))

	// multiply follows add
	expected = ast.NewInfixExpr(
		ast.NewValueExpr(ast.NewNumberLiteral("2")),
		ast.Operator{"+"},
		ast.NewInfixExpr(
			ast.NewValueExpr(ast.NewNumberLiteral("3")),
			ast.Operator{"*"},
			ast.NewValueExpr(ast.NewNumberLiteral("4")),
		),
	)

	assert.Equal(t, expected, parseExpression(`2 + 3 * 4`))

	// multiply follows grouped add
	expected = ast.NewInfixExpr(
		ast.NewGroupExpr(ast.NewInfixExpr(
			ast.NewValueExpr(ast.NewNumberLiteral("2")),
			ast.Operator{"+"},
			ast.NewValueExpr(ast.NewNumberLiteral("3")),
		)),
		ast.Operator{"*"},
		ast.NewValueExpr(ast.NewNumberLiteral("4")),
	)

	assert.Equal(t, expected, parseExpression(`(2 + 3) * 4`))

	// add and subtract associativity
	expected = ast.NewInfixExpr(
		ast.NewInfixExpr(
			ast.NewInfixExpr(
				ast.NewInfixExpr(
					ast.NewInfixExpr(
						ast.NewInfixExpr(
							ast.NewValueExpr(ast.NewNumberLiteral("2")),
							ast.Operator{"+"},
							ast.NewValueExpr(ast.NewNumberLiteral("3")),
						),
						ast.Operator{"+"},
						ast.NewValueExpr(ast.NewNumberLiteral("4")),
					),
					ast.Operator{"-"},
					ast.NewValueExpr(ast.NewNumberLiteral("5")),
				),
				ast.Operator{"+"},
				ast.NewValueExpr(ast.NewNumberLiteral("6")),
			),
			ast.Operator{"-"},
			ast.NewValueExpr(ast.NewNumberLiteral("7")),
		),
		ast.Operator{"-"},
		ast.NewValueExpr(ast.NewNumberLiteral("8")),
	)

	assert.Equal(t, expected, parseExpression(`2 + 3 + 4 - 5 + 6 - 7 - 8`))

	// multiply, divide, and modulus associativity
	expected = ast.NewInfixExpr(
		ast.NewInfixExpr(
			ast.NewInfixExpr(
				ast.NewInfixExpr(
					ast.NewInfixExpr(
						ast.NewInfixExpr(
							ast.NewValueExpr(ast.NewNumberLiteral("2")),
							ast.Operator{"/"},
							ast.NewValueExpr(ast.NewNumberLiteral("3")),
						),
						ast.Operator{"/"},
						ast.NewValueExpr(ast.NewNumberLiteral("4")),
					),
					ast.Operator{"*"},
					ast.NewValueExpr(ast.NewNumberLiteral("5")),
				),
				ast.Operator{"*"},
				ast.NewValueExpr(ast.NewNumberLiteral("6")),
			),
			ast.Operator{"%"},
			ast.NewValueExpr(ast.NewNumberLiteral("7")),
		),
		ast.Operator{"/"},
		ast.NewValueExpr(ast.NewNumberLiteral("8")),
	)

	assert.Equal(t, expected, parseExpression(`2 / 3 / 4 * 5 * 6 % 7 / 8`))

	// signed addition
	expected = ast.NewInfixExpr(
		ast.NewPrefixExpr(
			ast.Operator{"-"},
			ast.NewValueExpr(ast.NewNumberLiteral("2")),
		),
		ast.Operator{"+"},
		ast.NewPrefixExpr(
			ast.Operator{"+"},
			ast.NewValueExpr(ast.NewNumberLiteral("4")),
		),
	)

	assert.Equal(t, expected, parseExpression(`-2 + +4`))

	// signed subtraction
	expected = ast.NewInfixExpr(
		ast.NewPrefixExpr(
			ast.Operator{"-"},
			ast.NewValueExpr(ast.NewNumberLiteral("2")),
		),
		ast.Operator{"-"},
		ast.NewPrefixExpr(
			ast.Operator{"+"},
			ast.NewValueExpr(ast.NewNumberLiteral("4")),
		),
	)

	assert.Equal(t, expected, parseExpression(`-2 - +4`))

	// signed multiplication
	expected = ast.NewInfixExpr(
		ast.NewPrefixExpr(
			ast.Operator{"-"},
			ast.NewValueExpr(ast.NewNumberLiteral("2")),
		),
		ast.Operator{"*"},
		ast.NewPrefixExpr(
			ast.Operator{"+"},
			ast.NewValueExpr(ast.NewNumberLiteral("4")),
		),
	)

	assert.Equal(t, expected, parseExpression(`-2 * +4`))

	// signed division
	expected = ast.NewInfixExpr(
		ast.NewPrefixExpr(
			ast.Operator{"-"},
			ast.NewValueExpr(ast.NewNumberLiteral("2")),
		),
		ast.Operator{"/"},
		ast.NewPrefixExpr(
			ast.Operator{"+"},
			ast.NewValueExpr(ast.NewNumberLiteral("4")),
		),
	)

	assert.Equal(t, expected, parseExpression(`-2 / +4`))
}
