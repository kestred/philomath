package parser

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/scanner"
	"github.com/kestred/philomath/code/token"
)

// FIXME: The grammar feels like a hack because of the way functions are parsed
//
// Either allow semicolons to be skipped most of the time or also require
// semicolons at the end of struct/etc definitions.
//
// I've failed semicolon handling if I don't avoid the following Go weirdness:
//
//   (invalid)
//  foo := 3
//       + 4;
//
//   (valid)
//  foo := 3 +
//         4;
//
// It might also be ok to consider both of the above invalid and elide
// semicolons, but still allow a statement that looks like this:
//
//  foo := (13
//        + 4
//        - 7)
//

const MaxErrors = 8 // try to choose an actionable number of errors

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

func Make(filename string, trace bool, src []byte) Parser {
	p := Parser{}
	p.Init(filename, trace, src)
	return p
}

// Init prepares the parser p to convert a text src into an ast by starting
// a scanner, and scanning the the first token from the source.
func (p *Parser) Init(filename string, trace bool, src []byte) {
	scanError := func(pos token.Position, msg string) {
		p.error(pos, msg)
	}

	p.filename = filename
	p.scanner.Init(filename, src, scanError)
	p.operators.InitBuiltin()
	p.next()

	// don't trace first token
	p.trace = trace
}

func (p *Parser) ParseTop() *ast.TopScope {
	defer p.recoverStopped()
	var decls []ast.Decl
	for p.tok != token.END {
		decls = append(decls, p.parseDeclaration())
		for p.tok == token.SEMICOLON {
			p.next() // eat extra semicolons
		}
	}
	return ast.Top(decls)
}

func (p *Parser) ParseEvaluable() ast.Evaluable {
	defer p.recoverStopped()
	return p.parseEvaluable()
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
	if n > MaxErrors {
		p.stopParsing()
	}

	p.Errors = append(p.Errors, &ParseError{pos, msg})
}

func (p *Parser) expect(tok token.Token) bool {
	if p.tok != tok {
		// TODO: While this is easy to program, it makes for absolutely terrible
		// error messages in every single case that can be used.
		// Eventually, anywhere this is used should be replaced with a thought out message.
		p.error(p.scanner.Pos(), fmt.Sprintf(`Expected '%v' but received '%v'.`, tok, p.tok))

		/* TODO: Improvement for error messages...

		   It would probably produce better errors if:
		     When expecting a delimiter, don't eat a non-delimiter token
		     When expecting a non-delimiter, don't eat a delimiter token

		   The theory here is that most syntax errors occur at the boundary between
		   the spoken language and punctuation.  In particular, we tend to write the
		   the spoken part of the code correctly, but then have missing or extra punctuation.
		   By avoiding eating non-delimiters when expecting a delimiter, we prevent
		   ourselves from eating the first part of the next valid expression when
		   the delimiter is missing.  On the other hand, if we are expecting
		   a non-delimiter token and find a delimiter, then it is likely that we forgot
		   to fill information in when writing an expression, like the type of an argument.

		   In some cases that might result in an infinite loop since errors from
		   the same source line are discared.  One way of avoiding the possibility
		   of an infinite loop would be by quitting if the token position hasn't
		   changed since at all since the last error (or after a few errors).

		   In the aborted infinite loop case, it probably isn't worth parsing the
		   rest of this file because it is clearly irregular.  If we haven't reported
		   that many errors, it could still be worth parsing the other files so that
		   we can report an actionable number of errors.

		   TODO: Advanced mode error recovery...

		   Recovery (eg. error quality) might also be better if a parsing
		   function could register a closing delimiter and then parsing unwinds
		   immediately to that function if the delimiter is found.
		*/
		p.next()
		return false
	} else {
		p.next()
		return true
	}
}

// TODO: This makes for absolutely terrible error messages.
// Eventually, anywhere this is used should be replaced with a thought out message.
func (p *Parser) expected(what string) {
	p.error(p.scanner.Pos(), fmt.Sprintf(`Expected '%v' but received '%v'.`, what, p.tok))
	p.next() // eat something to make sure we don't infinite loop
}

func (p *Parser) next() {
	if p.trace {
		pc, _, line, _ := runtime.Caller(1)
		path := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		name := path[len(path)-1]
		// ignore expect and expected
		if len(name) >= 6 && name[0:6] == "expect" {
			pc, _, line, _ = runtime.Caller(2)
			path = strings.Split(runtime.FuncForPC(pc).Name(), ".")
			name = path[len(path)-1]
		}
		caller := "Parser." + name
		lit := p.lit
		if len(lit) > 7 {
			lit = lit[0:6] + "~"
		}
		// NOTE: For some irritating reason, running `go test` will always
		//       hide stderr so make sure we use stdout
		fmt.Printf(" %7.7s : %-14s @ %v:%v\n", lit, p.tok, caller, line)
	}

	p.pos, p.tok, p.lit = p.scanner.Scan()
}

func (p *Parser) parseBlock() *ast.Block {
	var directives []string
	if p.tok == token.DIRECTIVE {
		directives = append(directives, p.lit)
		p.next()
	}

	if p.tok == token.COLON {
		if len(directives) > 0 {
			p.error(p.scanner.Pos(), "Use of short block syntax is not allowed after specifying directives")
		}
		p.next() // eat ":"
		stmt := p.parseStatement()
		return ast.Blok([]ast.Evaluable{stmt})
	}

	// TODO: Proper directive handling
	if len(directives) > 0 && directives[0] == "asm" {
		var start, end, depth int
		depth = 1
		start = p.pos + 1
		linepos := p.scanner.Pos()
		p.expect(token.LEFT_BRACE)
		for depth > 0 {
			switch p.tok {
			case token.LEFT_BRACE:
				depth += 1
			case token.RIGHT_BRACE:
				depth -= 1
				end = p.pos - 1
			case token.END:
				p.error(linepos, `Missing a matching '}' to close an "#asm" block`)
				p.stopParsing()
			}
			p.next()
		}

		asm := ast.Asm(p.scanner.SourceAt(start, end))
		parseAssembly(asm)
		return ast.Blok([]ast.Evaluable{asm})
	} else {
		p.expect(token.LEFT_BRACE)
		for p.tok == token.SEMICOLON {
			p.next() // eat leading semicolons
		}

		var stmts []ast.Evaluable
		for p.tok != token.RIGHT_BRACE && p.tok != token.END {
			stmts = append(stmts, p.parseEvaluable())
			for p.tok == token.SEMICOLON {
				p.next() // eat extra semicolons
			}
		}
		p.expect(token.RIGHT_BRACE)
		return ast.Blok(stmts)
	}
}

func (p *Parser) parseEvaluable() ast.Evaluable {
	if p.tok == token.LEFT_BRACE || p.tok == token.DIRECTIVE {
		return p.parseBlock()
	} else if p.tok != token.IDENT {
		return p.parseStatement()
	}

	next := p.scanner.Peek()
	if next == token.CONS || next == token.COLON {
		return p.parseDeclaration()
	} else {
		return p.parseStatement()
	}
}

func (p *Parser) parseDeclaration() ast.Decl {
	name := p.lit
	p.expect(token.IDENT)

	if p.tok == token.COLON {
		// parse mutable decl
		p.next() // eat ":"
		if p.tok != token.EQUALS {
			panic("TODO: Handle typed declarations")
		}
		p.expect(token.EQUALS)
		expr := p.parseExpression()
		p.expect(token.SEMICOLON)
		return ast.Mutable(name, nil, expr)
	}

	// parse const decl
	p.expect(token.CONS)
	switch p.tok {
	case token.STRUCT:
		panic("TODO: Handle structs")
	case token.MODULE:
		panic("TODO: Handle modules")
	default:
		expr := p.parseExpression()
		if _, isFunc := expr.(*ast.ProcedureExpr); !isFunc {
			p.expect(token.SEMICOLON)
		}
		return ast.Immutable(name, ast.Constant(expr))
	}
}

func (p *Parser) parseStatement() ast.Stmt {
	if p.tok.IsKeyword() {
		panic("TODO: Handle keyword statements")
	}

	exprs := p.parseExpressionList()
	// TODO: Handle combined assignment (eg. a += 2);
	//       Generate an error but still continue if it is (a + = 2)
	if p.tok == token.EQUALS {
		p.next() // eat '='
		values := p.parseExpressionList()
		p.expect(token.SEMICOLON)
		return ast.Assign(exprs, nil, values)
	} else if len(exprs) == 1 {
		p.expect(token.SEMICOLON)
		return ast.Eval(exprs[0])
	} else {
		panic("TODO: Produce error for expression list w/o assignment, then recover")
	}
}

func (p *Parser) parseExpressionList() []ast.Expr {
	list := []ast.Expr{p.parseExpression()}
	for p.tok == token.COMMA {
		p.next() // eat ','
		list = append(list, p.parseExpression())
	}
	return list
}

func (p *Parser) parseExpression() ast.Expr {
	return p.parseOperators(0)
}

func (p *Parser) parseOperators(precedence ast.OpPrecedence) ast.Expr {
	lhs := p.parseBaseExpression()
	if p.tok == token.LEFT_BRACKET {
		panic("TODO: Hande array subscript")
	} else if p.tok == token.LEFT_PAREN {
		panic("TODO: Handle function call")
	} else if !p.tok.IsOperator() {
		return lhs
	}

	op := p.parseBinaryOperator()
	consumable := ast.MaxPrecedence
	for (op.Type == ast.BinaryInfix || op.Type == ast.UnaryPostfix) &&
		(precedence <= op.Precedence && op.Precedence <= consumable) {

		p.next() // eat operator
		if op.Type == ast.BinaryInfix {
			rhs := p.parseOperators(rightPrec(op))
			lhs = ast.InExp(lhs, op, rhs)
		} else {
			lhs = ast.PostExp(lhs, op)
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

func (p *Parser) parseBinaryOperator() *ast.OperatorDefn {
	options, defined := p.operators.Lookup(p.lit)
	if !defined {
		panic("TODO: Handle undefined operators")
	}

	var op *ast.OperatorDefn
	for _, opt := range options {
		if opt.Type == ast.BinaryInfix || opt.Type == ast.UnaryPostfix {
			op = opt
			break
		}
	}

	if op == nil {
		panic("TODO: Handle operator is not an infix/postfix operator")
	}

	return op
}

func rightPrec(op *ast.OperatorDefn) ast.OpPrecedence {
	if op.Associative == ast.RightAssociative {
		return op.Precedence
	} else {
		return op.Precedence + 1
	}
}

func nextPrec(op *ast.OperatorDefn) ast.OpPrecedence {
	if op.Associative == ast.LeftAssociative || op.Type == ast.UnaryPostfix {
		return op.Precedence
	} else {
		return op.Precedence - 1
	}
}

func (p *Parser) parseBaseExpression() ast.Expr {
	/* handle prefix expression */
	if p.tok.IsOperator() {
		options, defined := p.operators.Lookup(p.lit)
		if !defined {
			panic("TODO: Handle undefined operators")
		}

		var op *ast.OperatorDefn
		for _, opt := range options {
			if opt.Type == ast.UnaryPrefix {
				op = opt
				break
			}
		}

		if op == nil {
			panic("TODO: Handle operator is not a prefix operator")
		}

		p.next() // eat operator
		expr := p.parseOperators(ast.PrefixPrec)
		return ast.PreExp(op, expr)
	}

	switch p.tok {
	case token.LEFT_PAREN:
		p.next() // eat left paren
		var block *ast.Block
		var returnType ast.Type
		var parameters []ast.ProcedureParam
		if p.tok == token.IDENT && p.scanner.Peek() == token.COLON {
			for p.tok == token.IDENT {
				/* handle procedure parameters */
				// name := p.lit
				p.next() // eat ident
				p.expect(token.COLON)
				panic("TODO: Implement type parsing")
			}
		}

		if len(parameters) > 0 || p.tok == token.RIGHT_PAREN {
			p.expect(token.RIGHT_PAREN) // eat right paren

			/* handle procedure expression */
			if p.tok == token.ARROW {
				p.next() // eat arrow
				panic("TODO: Implement type parsing")
				if p.tok == token.COLON { // #[Compiler Message]
					p.error(p.scanner.Pos(), "Use of short block syntax is not allowed after specifying a return type")
				}
			}

			block = p.parseBlock()
			return ast.ProcExp(parameters, returnType, block)
		} else {
			/* handle grouped expression */
			expr := p.parseOperators(0)
			p.expect(token.RIGHT_PAREN)
			return ast.GrpExp(expr)
		}

	case token.IDENT:
		name := p.lit
		p.next() // eat ident
		return ast.Ident(name)

	case token.TEXT:
		expr := ast.TxtLit(p.lit)
		p.next() // eat text
		return expr

	case token.NUMBER:
		expr := ast.NumLit(p.lit)
		p.next() // eat number
		return expr

	default:
		p.expected("a value")
		return nil // TODO: maybe return BadExpr?
	}
}

type Operators struct {
	literals map[string][]*ast.OperatorDefn
}

func (o *Operators) InitBuiltin() {
	o.literals = make(map[string][]*ast.OperatorDefn)
	// logic operators
	o.defineHACKY(ast.BuiltinLogicalOr)
	o.defineHACKY(ast.BuiltinLogicalAnd)
	o.defineHACKY(ast.BuiltinElementOf)
	o.defineHACKY(ast.BuiltinNotElementOf)
	// comparison operators
	o.defineHACKY(ast.BuiltinIdentical)
	o.defineHACKY(ast.BuiltinEqual)
	o.defineHACKY(ast.BuiltinLess)
	o.defineHACKY(ast.BuiltinLessOrEqual)
	o.defineHACKY(ast.BuiltinGreater)
	o.defineHACKY(ast.BuiltinGreaterOrEqual)
	// arithmetic operators
	o.defineHACKY(ast.BuiltinCompare)
	o.defineHACKY(ast.BuiltinAdd)
	o.defineHACKY(ast.BuiltinSubtract)
	o.defineHACKY(ast.BuiltinMultiply)
	o.defineHACKY(ast.BuiltinDivide)
	o.defineHACKY(ast.BuiltinRemainder)
	o.defineHACKY(ast.BuiltinPositive)
	o.defineHACKY(ast.BuiltinNegative)
	// pointer operators
	o.defineHACKY(ast.BuiltinReference)
	o.defineHACKY(ast.BuiltinDereference)
}

func (o *Operators) defineHACKY(op *ast.OperatorDefn) {
	// TODO: Check that the operator has valid values and isn't stepping on any toes
	o.literals[op.Literal] = append(o.literals[op.Literal], op)
}

func (o *Operators) Lookup(literal string) ([]*ast.OperatorDefn, bool) {
	operators, exists := o.literals[literal]
	return operators, exists
}
