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
	expected = &ast.InfixExpr{
		&ast.InfixExpr{
			&ast.ValueExpr{&ast.NumberLiteral{"2"}},
			ast.Operator{"*"},
			&ast.ValueExpr{&ast.NumberLiteral{"3"}},
		},
		ast.Operator{"+"},
		&ast.ValueExpr{&ast.NumberLiteral{"4"}},
	}

	assert.Equal(t, expected, parseExpression(`2 * 3 + 4`))

	// multiply follows add
	expected = &ast.InfixExpr{
		&ast.ValueExpr{&ast.NumberLiteral{"2"}},
		ast.Operator{"+"},
		&ast.InfixExpr{
			&ast.ValueExpr{&ast.NumberLiteral{"3"}},
			ast.Operator{"*"},
			&ast.ValueExpr{&ast.NumberLiteral{"4"}},
		},
	}

	assert.Equal(t, expected, parseExpression(`2 + 3 * 4`))

	// multiply follows grouped add
	expected = &ast.InfixExpr{
		&ast.GroupExpr{&ast.InfixExpr{
			&ast.ValueExpr{&ast.NumberLiteral{"2"}},
			ast.Operator{"+"},
			&ast.ValueExpr{&ast.NumberLiteral{"3"}},
		}},
		ast.Operator{"*"},
		&ast.ValueExpr{&ast.NumberLiteral{"4"}},
	}

	assert.Equal(t, expected, parseExpression(`(2 + 3) * 4`))

	// add and subtract associativity
	expected = &ast.InfixExpr{
		&ast.InfixExpr{
			&ast.InfixExpr{
				&ast.InfixExpr{
					&ast.InfixExpr{
						&ast.InfixExpr{
							&ast.ValueExpr{&ast.NumberLiteral{"2"}},
							ast.Operator{"+"},
							&ast.ValueExpr{&ast.NumberLiteral{"3"}},
						},
						ast.Operator{"+"},
						&ast.ValueExpr{&ast.NumberLiteral{"4"}},
					},
					ast.Operator{"-"},
					&ast.ValueExpr{&ast.NumberLiteral{"5"}},
				},
				ast.Operator{"+"},
				&ast.ValueExpr{&ast.NumberLiteral{"6"}},
			},
			ast.Operator{"-"},
			&ast.ValueExpr{&ast.NumberLiteral{"7"}},
		},
		ast.Operator{"-"},
		&ast.ValueExpr{&ast.NumberLiteral{"8"}},
	}

	assert.Equal(t, expected, parseExpression(`2 + 3 + 4 - 5 + 6 - 7 - 8`))

	// multiply, divide, and modulus associativity
	expected = &ast.InfixExpr{
		&ast.InfixExpr{
			&ast.InfixExpr{
				&ast.InfixExpr{
					&ast.InfixExpr{
						&ast.InfixExpr{
							&ast.ValueExpr{&ast.NumberLiteral{"2"}},
							ast.Operator{"/"},
							&ast.ValueExpr{&ast.NumberLiteral{"3"}},
						},
						ast.Operator{"/"},
						&ast.ValueExpr{&ast.NumberLiteral{"4"}},
					},
					ast.Operator{"*"},
					&ast.ValueExpr{&ast.NumberLiteral{"5"}},
				},
				ast.Operator{"*"},
				&ast.ValueExpr{&ast.NumberLiteral{"6"}},
			},
			ast.Operator{"%"},
			&ast.ValueExpr{&ast.NumberLiteral{"7"}},
		},
		ast.Operator{"/"},
		&ast.ValueExpr{&ast.NumberLiteral{"8"}},
	}

	assert.Equal(t, expected, parseExpression(`2 / 3 / 4 * 5 * 6 % 7 / 8`))

	// signed addition
	expected = &ast.InfixExpr{
		&ast.PrefixExpr{ast.Operator{"-"}, &ast.ValueExpr{&ast.NumberLiteral{"2"}}},
		ast.Operator{"+"},
		&ast.PrefixExpr{ast.Operator{"+"}, &ast.ValueExpr{&ast.NumberLiteral{"4"}}},
	}

	assert.Equal(t, expected, parseExpression(`-2 + +4`))

	// signed subtraction
	expected = &ast.InfixExpr{
		&ast.PrefixExpr{ast.Operator{"-"}, &ast.ValueExpr{&ast.NumberLiteral{"2"}}},
		ast.Operator{"-"},
		&ast.PrefixExpr{ast.Operator{"+"}, &ast.ValueExpr{&ast.NumberLiteral{"4"}}},
	}

	assert.Equal(t, expected, parseExpression(`-2 - +4`))

	// signed multiplication
	expected = &ast.InfixExpr{
		&ast.PrefixExpr{ast.Operator{"-"}, &ast.ValueExpr{&ast.NumberLiteral{"2"}}},
		ast.Operator{"*"},
		&ast.PrefixExpr{ast.Operator{"+"}, &ast.ValueExpr{&ast.NumberLiteral{"4"}}},
	}

	assert.Equal(t, expected, parseExpression(`-2 * +4`))

	// signed division
	expected = &ast.InfixExpr{
		&ast.PrefixExpr{ast.Operator{"-"}, &ast.ValueExpr{&ast.NumberLiteral{"2"}}},
		ast.Operator{"/"},
		&ast.PrefixExpr{ast.Operator{"+"}, &ast.ValueExpr{&ast.NumberLiteral{"4"}}},
	}

	assert.Equal(t, expected, parseExpression(`-2 / +4`))
}
