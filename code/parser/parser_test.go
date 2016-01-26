package parser

import (
	"testing"

	"github.com/kestred/philomath/code/ast"
	"github.com/stretchr/testify/assert"
)

func TestParseMain(t *testing.T) {
	input := `main :: () { somevar := 1; }`
	parser := Make("example", false, []byte(input))
	top := parser.ParseTop()
	if assert.Empty(t, parser.Errors) {
		expected := ast.Top([]ast.Decl{
			ast.Immutable("main", ast.Constant(
				ast.ProcExp(nil, nil, ast.Blok([]ast.Evaluable{
					ast.Mutable("somevar", nil, ast.NumLit("1")),
				})),
			)),
		})
		assert.Equal(t, expected, top)
	}
}

func TestParseError(t *testing.T) {
	var parser Parser
	parser.Init("error.phi", false, []byte(`1 * (2 + 3} - 4`))
	parser.ParseEvaluable()
	if assert.True(t, len(parser.Errors) > 0, "Expected some errors but found none.") {
		assert.Equal(t, "error.phi:1:12: Expected ')' but received '}'.", parser.Errors[0].Error())
	}

	parser = Parser{}
	parser.Init("error.phi", false, []byte(`{ 1 - 4 }`))
	parser.ParseEvaluable()
	if assert.True(t, len(parser.Errors) > 0, "Expected some errors but found none.") {
		assert.Equal(t, "error.phi:1:10: Expected ';' but received '}'.", parser.Errors[0].Error())
	}
}

func parseAny(t *testing.T, input string) ast.Node {
	p := Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	return node
}

func TestParseDeclarations(t *testing.T) {
	expected := ast.Blok([]ast.Evaluable{
		ast.Mutable("foo", nil, ast.NumLit("3")),
		ast.Immutable("baz", ast.Constant(ast.NumLit("1"))),
		ast.Eval(ast.InExp(
			ast.InExp(ast.NumLit("2"), ast.BuiltinAdd, ast.Ident("foo")),
			ast.BuiltinAdd,
			ast.Ident("baz"),
		)),
	})

	assert.Equal(t, expected, parseAny(t, `{
		foo := 3;      # mutable declaration
		baz :: 1;      # constant definition
		2 + foo + baz; # evaluated statement
	}`))
}

func TestParseProcedures(t *testing.T) {
	expected := ast.Blok([]ast.Evaluable{
		ast.Immutable("immutable", ast.Constant(
			ast.ProcExp(nil, nil, ast.Blok([]ast.Evaluable{
				ast.Mutable("dragonfruit", nil, ast.NumLit("3.0")),
				ast.Eval(ast.InExp(ast.Ident("dragonfruit"), ast.BuiltinDivide, ast.NumLit("1.0"))),
			})),
		)),
		ast.Mutable("mutable", nil,
			ast.ProcExp(nil, nil, ast.Blok([]ast.Evaluable{
				ast.Eval(ast.InExp(
					ast.InExp(ast.NumLit("0700"), ast.BuiltinAdd, ast.NumLit("0040")),
					ast.BuiltinAdd,
					ast.NumLit("0004"),
				)),
			})),
		),
		ast.Immutable("outer", ast.Constant(
			ast.ProcExp(nil, nil, ast.Blok([]ast.Evaluable{
				ast.Immutable("inner", ast.Constant(ast.ProcExp(nil, nil, ast.Blok(nil)))),
			})),
		)),
		ast.Immutable("short", ast.Constant(
			ast.ProcExp(nil, nil, ast.Blok([]ast.Evaluable{
				ast.Eval(ast.TxtLit(`"Hello world"`)),
			})),
		)),
	})

	assert.Equal(t, expected, parseAny(t, `{
		# immutable procedure
		immutable :: () {
			dragonfruit := 3.0;
			dragonfruit / 1.0;
		}

		# mutable procedure
		mutable := () {
			0700 + 0040 + 0004;
		};

		# nested procedures
		outer :: () {
			# empty procedure
			inner :: () {
			}
		}

		# short procedure
		short :: (): "Hello world";
	}`))
}

func TestParseBlocks(t *testing.T) {
	expected := ast.Blok([]ast.Evaluable{
		ast.Blok([]ast.Evaluable{
			ast.Mutable("bar", nil, ast.Ident("foo")),
			ast.Eval(ast.InExp(ast.NumLit("0755"), ast.BuiltinSubtract, ast.Ident("baz"))),
			ast.Assign(
				[]ast.Expr{ast.Ident("foo")}, nil,
				[]ast.Expr{ast.InExp(ast.Ident("baz"), ast.BuiltinMultiply, ast.NumLit("4"))},
			),
			ast.Assign(
				[]ast.Expr{ast.Ident("bar"), ast.Ident("foo")}, nil,
				[]ast.Expr{
					ast.InExp(ast.Ident("foo"), ast.BuiltinAdd, ast.NumLit("27")),
					ast.Ident("bar"),
				},
			),
		}),
		ast.Eval(ast.InExp(ast.NumLit("8.4e-5"), ast.BuiltinDivide, ast.NumLit("0.5"))),
		ast.Blok(nil),
	})

	assert.Equal(t, expected, parseAny(t, `{
		; # ignore extra semicolons occuring before a statement

		# a nested block
		{
			bar := foo;
			0755 - baz;

			foo = baz * 4;		        # assignment
			bar, foo = foo + 27, bar; # parallel assignment
		}

		# ignore extra semicolons
		; ;
		8.4e-5 / 0.5;; ;
		;

		# empty block
		{
		}
	}`))
}

func parseExpr(t *testing.T, input string) ast.Expr {
	p := Make("example", false, []byte(input))
	node := p.ParseEvaluable()
	assert.Empty(t, p.Errors, "Unexpected parser errors")
	return node.(*ast.EvalStmt).Expr
}

func TestParseArithmetic(t *testing.T) {
	var expected ast.Expr

	// add follows multiply
	expected = ast.InExp(
		ast.InExp(
			ast.NumLit("2"),
			ast.BuiltinMultiply,
			ast.NumLit("3"),
		),
		ast.BuiltinAdd,
		ast.NumLit("4"),
	)

	assert.Equal(t, expected, parseExpr(t, `2 * 3 + 4;`))

	// multiply follows add
	expected = ast.InExp(
		ast.NumLit("2"),
		ast.BuiltinAdd,
		ast.InExp(
			ast.NumLit("3"),
			ast.BuiltinMultiply,
			ast.NumLit("4"),
		),
	)

	assert.Equal(t, expected, parseExpr(t, `2 + 3 * 4;`))

	// multiply follows grouped add
	expected = ast.InExp(
		ast.GrpExp(ast.InExp(
			ast.NumLit("2"),
			ast.BuiltinAdd,
			ast.NumLit("3"),
		)),
		ast.BuiltinMultiply,
		ast.NumLit("4"),
	)

	assert.Equal(t, expected, parseExpr(t, `(2 + 3) * 4;`))

	// add and subtract associativity
	expected = ast.InExp(
		ast.InExp(
			ast.InExp(
				ast.InExp(
					ast.InExp(
						ast.InExp(
							ast.NumLit("2"),
							ast.BuiltinAdd,
							ast.NumLit("3"),
						),
						ast.BuiltinAdd,
						ast.NumLit("4"),
					),
					ast.BuiltinSubtract,
					ast.NumLit("5"),
				),
				ast.BuiltinAdd,
				ast.NumLit("6"),
			),
			ast.BuiltinSubtract,
			ast.NumLit("7"),
		),
		ast.BuiltinSubtract,
		ast.NumLit("8"),
	)

	assert.Equal(t, expected, parseExpr(t, `2 + 3 + 4 - 5 + 6 - 7 - 8;`))

	// multiply, divide, and modulus associativity
	expected = ast.InExp(
		ast.InExp(
			ast.InExp(
				ast.InExp(
					ast.InExp(
						ast.InExp(
							ast.NumLit("2"),
							ast.BuiltinDivide,
							ast.NumLit("3"),
						),
						ast.BuiltinDivide,
						ast.NumLit("4"),
					),
					ast.BuiltinMultiply,
					ast.NumLit("5"),
				),
				ast.BuiltinMultiply,
				ast.NumLit("6"),
			),
			ast.BuiltinRemainder,
			ast.NumLit("7"),
		),
		ast.BuiltinDivide,
		ast.NumLit("8"),
	)

	assert.Equal(t, expected, parseExpr(t, `2 / 3 / 4 * 5 * 6 % 7 / 8;`))

	// signed addition
	expected = ast.InExp(
		ast.PreExp(ast.BuiltinNegative, ast.NumLit("2")),
		ast.BuiltinAdd,
		ast.PreExp(ast.BuiltinPositive, ast.NumLit("4")),
	)

	assert.Equal(t, expected, parseExpr(t, `-2 + +4;`))

	// signed subtraction
	expected = ast.InExp(
		ast.PreExp(ast.BuiltinNegative, ast.NumLit("2")),
		ast.BuiltinSubtract,
		ast.PreExp(ast.BuiltinPositive, ast.NumLit("4")),
	)

	assert.Equal(t, expected, parseExpr(t, `-2 - +4;`))

	// signed multiplication
	expected = ast.InExp(
		ast.PreExp(ast.BuiltinNegative, ast.NumLit("2")),
		ast.BuiltinMultiply,
		ast.PreExp(ast.BuiltinPositive, ast.NumLit("4")),
	)

	assert.Equal(t, expected, parseExpr(t, `-2 * +4;`))

	// signed division
	expected = ast.InExp(
		ast.PreExp(ast.BuiltinNegative, ast.NumLit("2")),
		ast.BuiltinDivide,
		ast.PreExp(ast.BuiltinPositive, ast.NumLit("4")),
	)

	assert.Equal(t, expected, parseExpr(t, `-2 / +4;`))
}
