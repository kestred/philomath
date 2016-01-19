package parser

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/kestred/philomath/ast"
	"github.com/kestred/philomath/scanner"
	"github.com/kestred/philomath/token"
	"github.com/kestred/philomath/utils"
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

	// public state
	Errors []error
}

// Init prepares the parser p to convert a text src into an ast by starting
// a scanner, and scanning the the first token from the source.
func (p *Parser) Init(filename string, trace bool, src []byte) {
	scanError := func(pos token.Position, msg string) {
		p.error(pos, msg)
	}

	p.filename = filename
	p.trace = trace
	p.scanner.Init(filename, src, scanError)
	p.operators.InitBuiltin()
	p.next()
}

func (p *Parser) Parse() ast.Node {
	defer p.recoverStopped()
	return &ast.Root{}
}

func (p *Parser) ParseExpression() ast.Expr {
	defer p.recoverStopped()
	return p.parseExpression()
}

// A stopParsing panic is raised to indicate early termination.
//
// In most cases I consider panics to be a code smell when they are used for
// control flow.  In this case though, it is far easier to use a panic for
// early termination than it would be to return and check for errors everywhere.
//
// One alternative is to use an error handler like in the scanner and defer the
// error handling up one level but callbacks can also lead to unhappy times.
type stopParsing struct{}

func (p *Parser) stopParsing() {
	panic(stopParsing{})
}

func (p *Parser) recoverStopped() {
	if e := recover(); e != nil {
		if _, ok := e.(stopParsing); !ok {
			panic(e)
		}
	}
}

// TODO: international translations for compiler error messages
func (p *Parser) error(pos token.Position, msg string) {
	n := len(p.Errors)
	if n > 0 && p.Errors[n-1].(*ParseError).Pos.Line == pos.Line {
		return // discard - likely a spurious error
	}
	if n > 8 {
		p.stopParsing()
	}

	p.Errors = append(p.Errors, &ParseError{pos, msg})
}

func (p *Parser) expect(tok token.Token) {
	if p.tok != tok {
		// TODO: While this is easy to program, it makes for absolutely terrible
		// error messages in every single case that can be used.
		// Eventually, anywhere this is used should be replaced with a thought out message.
		p.error(p.scanner.Pos(), fmt.Sprintf(`Expected '%v' but recieved '%v'.`, p.tok, tok))
	}
	p.next()
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

func (p *Parser) parseExpression() ast.Expr {
	return p.parseOperators(0)
}

func (p *Parser) parseOperators(precedence Precedence) ast.Expr {
	lhs := p.parseBaseExpression()
	if p.tok == token.LEFT_BRACKET {
		panic("TODO: Hande array subscript")
	} else if p.tok == token.LEFT_PAREN {
		panic("TODO: Handle function call")
	} else if !p.tok.IsOperator() {
		return lhs
	}

	op := p.parseBinaryOperator()
	consumable := MaxPrecedence
	for (op.Type == BinaryInfix || op.Type == UnaryPostfix) &&
		(precedence <= op.Precedence && op.Precedence <= consumable) {

		p.next() // consume operator
		if op.Type == BinaryInfix {
			rhs := p.parseOperators(rightPrec(op))
			lhs = ast.InExp(lhs, ast.Operator{op.Literal}, rhs)
		} else {
			lhs = ast.PostExp(lhs, ast.Operator{op.Literal})
		}

		if p.tok.IsOperator() {
			op = p.parseBinaryOperator()
			consumable = nextPrec(op)
		} else {
			break
		}
	}

	return lhs
}

func (p *Parser) parseBinaryOperator() Operator {
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

	utils.Assert(op.Type != Nullary, "Found a nullary operator")

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

		utils.Assert(op.Type != Nullary, "Found a nullary operator")

		p.next() // consume operator
		expr := p.parseOperators(PrefixPrecedence)
		return ast.PreExp(ast.Operator{op.Literal}, expr)
	}

	switch p.tok {
	// handle grouped expression
	case token.LEFT_PAREN:
		p.next() // consume left parent
		expr := p.parseOperators(0)
		p.expect(token.RIGHT_PAREN)
		return ast.GrpExp(expr)

	// handle literals
	case token.IDENT, token.TEXT:
		panic("TODO: I'll handle these eventually")

	case token.NUMBER:
		expr := ast.ValExp(ast.NumLit(p.lit))
		p.next()
		return expr
	}

	panic("TODO: not implemented")
	return nil // TODO: Implement
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
	Overload    string
	Type        OperatorType
	Associative Associative
	Precedence  Precedence
}

type Operators struct {
	literals map[string][]Operator
}

func (o *Operators) InitBuiltin() {
	o.literals = make(map[string][]Operator)
	// logic operators
	o.defineHACKY(Operator{"Logical Or", "or", "_or_", BinaryInfix, LeftAssociative, LogicalPrecedence})
	o.defineHACKY(Operator{"Logical And", "and", "_and_", BinaryInfix, LeftAssociative, LogicalPrecedence + 1})
	o.defineHACKY(Operator{"Inclusion", "in", "_in_", BinaryInfix, LeftAssociative, LogicalPrecedence + 1})
	// relation operators
	o.defineHACKY(Operator{"Identical", "is", "_is_", BinaryInfix, NonAssociative, RelationPrecedence})
	o.defineHACKY(Operator{"Equal", "==", "_eq_", BinaryInfix, NonAssociative, RelationPrecedence})
	o.defineHACKY(Operator{"Less", "<", "_lt_", BinaryInfix, NonAssociative, RelationPrecedence})
	o.defineHACKY(Operator{"Less or Equal", "<=", "_lte_", BinaryInfix, NonAssociative, RelationPrecedence})
	o.defineHACKY(Operator{"Greater", ">", "_gt_", BinaryInfix, NonAssociative, RelationPrecedence})
	o.defineHACKY(Operator{"Greater or Equal", ">=", "_gte_", BinaryInfix, NonAssociative, RelationPrecedence})
	// arithmetic operators
	o.defineHACKY(Operator{"Compare", "<=>", "_cmp_", BinaryInfix, LeftAssociative, InfixPrecedence})
	o.defineHACKY(Operator{"Add", "+", "_add_", BinaryInfix, LeftAssociative, InfixPrecedence + 1})
	o.defineHACKY(Operator{"Subtract", "-", "_sub_", BinaryInfix, LeftAssociative, InfixPrecedence + 1})
	o.defineHACKY(Operator{"Multiply", "*", "_mul_", BinaryInfix, LeftAssociative, InfixPrecedence + 2})
	o.defineHACKY(Operator{"Divide", "/", "_div_", BinaryInfix, LeftAssociative, InfixPrecedence + 2})
	o.defineHACKY(Operator{"Remainder", "%", "_rem_", BinaryInfix, LeftAssociative, InfixPrecedence + 2})
	o.defineHACKY(Operator{"Positive", "+", "pos_", UnaryPrefix, RightAssociative, PrefixPrecedence})
	o.defineHACKY(Operator{"Negative", "-", "neg_", UnaryPrefix, RightAssociative, PrefixPrecedence})
	// pointer operators
	o.defineHACKY(Operator{"Reference", "^", "ref_", UnaryPrefix, RightAssociative, PrefixPrecedence})
	o.defineHACKY(Operator{"Dereference", "~", "deref_", UnaryPrefix, RightAssociative, PrefixPrecedence})
}

func (o *Operators) defineHACKY(op Operator) {
	// TODO: Check that the operator has valid values and isn't stepping on any toes
	o.literals[op.Literal] = append(o.literals[op.Literal], op)
}

func (o *Operators) Lookup(literal string) ([]Operator, bool) {
	operators, exists := o.literals[literal]
	return operators, exists
}
