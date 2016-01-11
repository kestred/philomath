package parser

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/kestred/philomath/ast"
	"github.com/kestred/philomath/scanner"
	"github.com/kestred/philomath/token"
)

type ParseError struct {
	Pos token.Position
	Msg string
}

func (e *ParseError) Error() string {
	return e.Pos.String() + ": " + e.Msg
}

// A parser holds the parser's internal state while processing
// a given text.  It can be allocated as part of another data
// structure but must be initialized via Init before use.
type Parser struct {
	// immutable state
	filename  string
	scanner   scanner.Scanner
	operators Operators
	trace     bool

	// parsing state
	pos int         // next token offset
	tok token.Token // next token type
	lit string      // next token literal

	errors []error
}

// Init prepares the parser p to convert a text src into an ast by starting
// a scanner, and scanning the the first token from the source.
func (p *Parser) Init(filename string, trace bool, src []byte) {
	scanError := func(pos token.Position, msg string) {
		p.errors = append(p.errors, &ParseError{pos, msg})
	}

	p.filename = filename
	p.trace = trace
	p.scanner.Init(filename, src, scanError)
	p.operators.InitBuiltin()
	p.next()
}

func (p *Parser) Parse() ast.Node {
	return &ast.Root{}
}

func (p *Parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()

	if p.trace {
		pc, _, _, _ := runtime.Caller(1)
		path := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		caller := "Parser." + path[len(path)-1]
		fmt.Errorf("%s : %-14s @ %v\n", p.lit, p.tok, caller)
	}
}

type OperatorType uint8
type Associative uint8
type Precedence int8

const (
	Nullary OperatorType = iota
	UnaryPrefix
	UnaryPostfix
	BinaryInfix
)

const (
	NonAssociative Associative = iota
	LeftAssociative
	RightAssociative
)

const (
	AssignmentPrecedence Precedence = 0
	LogicalPrecedence    Precedence = 15
	RelationPrecedence   Precedence = 31
	InfixPrecedence      Precedence = 47
	PrefixPrecedence     Precedence = 95
	PostfixPrecedence    Precedence = 111
	MaxPrecedence        Precedence = 127
)

type Operator struct {
	Name        string
	Literal     string
	Type        OperatorType
	Associative Associative
	Precedence  Precedence
}

type Operators struct {
	literals map[string][]Operator
}

func (o *Operators) InitBuiltin() {
	o.literals = make(map[string][]Operator)
	o.defineHACKY(Operator{"Logical Or", "or", BinaryInfix, LeftAssociative, LogicalPrecedence})
	o.defineHACKY(Operator{"Logical And", "and", BinaryInfix, LeftAssociative, LogicalPrecedence + 1})
	o.defineHACKY(Operator{"Add", "+", BinaryInfix, LeftAssociative, InfixPrecedence})
	o.defineHACKY(Operator{"Add", "+", BinaryInfix, LeftAssociative, InfixPrecedence})
	o.defineHACKY(Operator{"Subtract", "-", BinaryInfix, LeftAssociative, InfixPrecedence})
	o.defineHACKY(Operator{"Multiply", "*", BinaryInfix, LeftAssociative, InfixPrecedence + 1})
	o.defineHACKY(Operator{"Divide", "/", BinaryInfix, LeftAssociative, InfixPrecedence + 1})
	o.defineHACKY(Operator{"Remainder", "%", BinaryInfix, LeftAssociative, InfixPrecedence + 1})
	o.defineHACKY(Operator{"Positive", "+", UnaryPrefix, RightAssociative, PrefixPrecedence})
	o.defineHACKY(Operator{"Negative", "-", UnaryPrefix, RightAssociative, PrefixPrecedence})
	o.defineHACKY(Operator{"Address Of", "^", UnaryPrefix, RightAssociative, PrefixPrecedence})
	o.defineHACKY(Operator{"Dereference", "~", UnaryPrefix, RightAssociative, PrefixPrecedence})
	o.defineHACKY(Operator{"Descope", ".", UnaryPostfix, LeftAssociative, PostfixPrecedence})
}

func (o *Operators) defineHACKY(op Operator) {
	// TODO: Check that the operator has valid values and isn't stepping on any toes
	o.literals[op.Literal] = append(o.literals[op.Literal], op)
}

func (o *Operators) Lookup(literal string) ([]Operator, bool) {
	operators, exists := o.literals[literal]
	return operators, exists
}

func (p *Parser) ParseExpression() ast.Expr {
	return p.parseOperators(0)
}

func (p *Parser) parseOperators(precedence Precedence) ast.Expr {
	lhs := p.parseBaseExpression()
	if p.tok == token.LEFT_BRACKET {
		panic("TODO: Hande array subscript")
		// p.parse
	} else if p.tok == token.LEFT_PAREN {
		panic("TODO: Handle function call")
	} else if !p.tok.IsOperator() {
		return lhs
	}

	op := p.binaryOperator()
	consumable := MaxPrecedence
	for (op.Type == BinaryInfix || op.Type == UnaryPostfix) &&
		(precedence <= op.Precedence && op.Precedence <= consumable) {

		p.next() // consume operator
		if op.Type == BinaryInfix {
			rhs := p.parseOperators(rightPrec(op))
			lhs = &ast.InfixExpr{lhs, ast.Operator{op.Literal}, rhs}
		} else {
			lhs = &ast.PostfixExpr{lhs, ast.Operator{op.Literal}}
		}

		if p.tok.IsOperator() {
			op = p.binaryOperator()
			consumable = nextPrec(op)
		} else {
			break
		}
	}

	return lhs
}

func (p *Parser) binaryOperator() Operator {
	options, defined := p.operators.Lookup(p.lit)
	if !defined {
		panic("TODO: Handle undefined operators")
	}

	var op Operator
	for _, opt := range options {
		if opt.Type == BinaryInfix || opt.Type == UnaryPostfix {
			op = opt
			break
		}
	}

	if op.Type == Nullary {
		panic("TODO: Handle expected an infix or postfix operator")
	}

	return op
}

func rightPrec(op Operator) Precedence {
	if op.Associative == RightAssociative {
		return op.Precedence
	} else {
		return op.Precedence + 1
	}
}

func nextPrec(op Operator) Precedence {
	if op.Associative == LeftAssociative || op.Type == UnaryPostfix {
		return op.Precedence
	} else {
		return op.Precedence - 1
	}
}

func (p *Parser) parseBaseExpression() ast.Expr {
	// handle prefix expression
	if p.tok.IsOperator() {
		options, defined := p.operators.Lookup(p.lit)
		if !defined {
			panic("TODO: Handle undefined operators")
		}

		var op Operator
		for _, opt := range options {
			if opt.Type == UnaryPrefix {
				op = opt
				break
			}
		}

		if op.Type == Nullary {
			panic("TODO: Handle expected a prefix operator")
		}

		p.next() // consume operator
		expr := p.parseOperators(PrefixPrecedence)
		return &ast.PrefixExpr{ast.Operator{op.Literal}, expr}
	}

	switch p.tok {
	// handle grouped expression
	case token.LEFT_PAREN:
		p.next() // consume left parent

		expr := p.parseOperators(0)

		// Basically: expect(token.RightParen)
		if p.tok != token.RIGHT_PAREN {
			panic("TODO: Expected right paren")
		}
		p.next()

		return &ast.GroupExpr{expr}

	// handle literals
	case token.IDENT, token.TEXT:
		panic("TODO: I'll handle these eventually")

	case token.NUMBER:
		expr := &ast.ValueExpr{&ast.NumberLiteral{p.lit}}
		p.next()
		return expr
	}

	panic("NOT implemented")
	return nil // TODO: Implement
}
