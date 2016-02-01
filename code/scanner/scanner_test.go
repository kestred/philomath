package scanner

import (
	"testing"

	"github.com/kestred/philomath/code/token"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	code := `Main := () { mystery.Print(world.Greetings) }`

	failOnError := func(pos token.Position, msg string) {
		assert.Fail(t, "At Line %d, Col %d: %s", pos.Line, pos.Column, msg)
	}

	s := Scanner{}
	s.Init("main", []byte(code), failOnError)

	var tokens []token.Token
	MAX_ITER := 200 // Don't loop forever
	for i := 0; i < MAX_ITER; i++ {
		_, tok, _ := s.Scan()
		tokens = append(tokens, tok)
		if tok == token.INVALID || tok == token.END {
			break
		}
	}

	assert.Equal(t, []token.Token{
		/* MAIN := ()         */ token.IDENT, token.COLON, token.EQUALS, token.LEFT_PAREN, token.RIGHT_PAREN,
		/* { mystery.Print(   */ token.LEFT_BRACE, token.IDENT, token.PERIOD, token.IDENT, token.LEFT_PAREN,
		/* world.Greetings) } */ token.IDENT, token.PERIOD, token.IDENT, token.RIGHT_PAREN, token.RIGHT_BRACE,
		token.END,
	}, tokens)
}

type scanToken struct {
	pos int
	tok token.Token
	lit string
}

type scanError struct {
	pos  token.Position
	msg  string
	prev *scanError
}

func scanOnce(src string) (scanToken, *scanError) {
	var err *scanError
	handleError := func(pos token.Position, msg string) {
		err = &scanError{pos, msg, err}
	}

	var t scanToken
	s := Scanner{}
	s.Init("scanOnce", []byte(src), handleError)
	t.pos, t.tok, t.lit = s.Scan()

	return t, err
}

func scanAll(src string) *scanError {
	var err *scanError
	handleError := func(pos token.Position, msg string) {
		err = &scanError{pos, msg, err}
	}

	var t scanToken
	s := Scanner{}
	s.Init("scanAll", []byte(src), handleError)

	for i := 0; i < 9999; i++ {
		t.pos, t.tok, t.lit = s.Scan()
		if t.tok == token.END {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}

func TestSkipsWhitesace(t *testing.T) {
	scan, err := scanOnce("\n    for food in cornucopia\n")
	assert.Nil(t, err)
	assert.Equal(t, token.FOR, scan.tok)
	assert.Equal(t, 5, scan.pos)
	assert.Equal(t, "for", scan.lit)

	scan, err = scanOnce(`
    // line comment
/* block comment
*/
	  for food in corncupia//comment
`)
	assert.Nil(t, err)
	assert.Equal(t, token.FOR, scan.tok)
	assert.Equal(t, 44, scan.pos)
	assert.Equal(t, "for", scan.lit)

	scan, err = scanOnce("\n    \r\n   \n\r for food in cornucopia\n")
	assert.Nil(t, err)
	assert.Equal(t, token.FOR, scan.tok)
	assert.Equal(t, 13, scan.pos)
	assert.Equal(t, "for", scan.lit)
}

func TestScansComment(t *testing.T) {
	scan, err := scanOnce("// line comment\nfor food in cornucopia")
	assert.Equal(t, token.FOR, scan.tok)
	assert.Nil(t, err)

	scan, err = scanOnce("/* block comment */for food in cornucopia")
	assert.Equal(t, token.FOR, scan.tok)
	assert.Nil(t, err)

	scan, err = scanOnce("/* /* nested comment */ */ for food in cornucopia")
	assert.Equal(t, token.FOR, scan.tok)
	assert.Nil(t, err)

	scan, err = scanOnce("/*\nThis comment is unterminated\n")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `unterminated comment`, err.msg)
	}
}

func TestErrorsRespectWhitespace(t *testing.T) {
	scan, err := scanOnce("\n\n    ~\n")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 6, err.pos.Offset)
		assert.Equal(t, 3, err.pos.Line)
		assert.Equal(t, 5, err.pos.Column)
		assert.Equal(t, `unexpected character U+007E '~'`, err.msg)
	}
}

func TestScansIdentifier(t *testing.T) {
	scan, err := scanOnce(`justletters`)
	assert.Nil(t, err)
	assert.Equal(t, token.IDENT, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `justletters`, scan.lit)

	scan, err = scanOnce(`CamelCase`)
	assert.Nil(t, err)
	assert.Equal(t, token.IDENT, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `CamelCase`, scan.lit)

	scan, err = scanOnce(`snake_case`)
	assert.Nil(t, err)
	assert.Equal(t, token.IDENT, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `snake_case`, scan.lit)

	scan, err = scanOnce(`with123`)
	assert.Nil(t, err)
	assert.Equal(t, token.IDENT, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `with123`, scan.lit)

	scan, err = scanOnce(`1992president`)
	assert.NotNil(t, err)
	assert.Equal(t, token.INVALID, scan.tok)
	assert.Equal(t, 0, scan.pos)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `identifier cannot start with digits (eg. 0, 1, ..., 9)`, err.msg)
	}

	// guess it is supposed to be a number if it is longer than a year
	scan, err = scanOnce(`12345atStart`)
	assert.NotNil(t, err)
	assert.Equal(t, token.INVALID, scan.tok)
	assert.Equal(t, 0, scan.pos)
	if assert.NotNil(t, err) {
		assert.Equal(t, 5, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 6, err.pos.Column)
		assert.Equal(t, `missing space after number before "atStart"`, err.msg)
	}

	// guess it is supposed to be a number if it starts with a decimal point
	scan, err = scanOnce(`.25atStart`)
	assert.NotNil(t, err)
	assert.Equal(t, token.INVALID, scan.tok)
	assert.Equal(t, 0, scan.pos)
	if assert.NotNil(t, err) {
		assert.Equal(t, 3, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 4, err.pos.Column)
		assert.Equal(t, `missing space after number before "atStart"`, err.msg)
	}

	// guess it is supposed to be a number if it includes an exponent
	scan, err = scanOnce(`12e5atStart`)
	assert.NotNil(t, err)
	assert.Equal(t, token.INVALID, scan.tok)
	assert.Equal(t, 0, scan.pos)
	if assert.NotNil(t, err) {
		assert.Equal(t, 4, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 5, err.pos.Column)
		assert.Equal(t, `missing space after number before "atStart"`, err.msg)
	}

	// guess it is supposed to be a number if it includes a signed exponent
	scan, err = scanOnce(`12e-5atStart`)
	assert.NotNil(t, err)
	assert.Equal(t, token.INVALID, scan.tok)
	assert.Equal(t, 0, scan.pos)
	if assert.NotNil(t, err) {
		assert.Equal(t, 5, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 6, err.pos.Column)
		assert.Equal(t, `missing space after number before "atStart"`, err.msg)
	}
}

func TestScansStrings(t *testing.T) {
	scan, err := scanOnce(`"simple"`)
	assert.Nil(t, err)
	assert.Equal(t, token.TEXT, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `"simple"`, scan.lit)

	scan, err = scanOnce(`" white space "`)
	assert.Nil(t, err)
	assert.Equal(t, token.TEXT, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `" white space "`, scan.lit)

	scan, err = scanOnce(`"quote\""`)
	assert.Nil(t, err)
	assert.Equal(t, token.TEXT, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `"quote\""`, scan.lit)

	scan, err = scanOnce(`"escaped \n\r\b\t\f"`)
	assert.Nil(t, err)
	assert.Equal(t, token.TEXT, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `"escaped \n\r\b\t\f"`, scan.lit)

	scan, err = scanOnce(`"slashes \\ //"`)
	assert.Nil(t, err)
	assert.Equal(t, token.TEXT, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `"slashes \\ //"`, scan.lit)
}

func TestReportsUsefulStringErrors(t *testing.T) {
	scan, err := scanOnce(`"`)
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `unterminated text literal`, err.msg)
	}

	scan, err = scanOnce(`"No end quote`)
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `unterminated text literal`, err.msg)
	}

	scan, err = scanOnce("\"contains unescaped \u0007 control char\"")
	assert.Equal(t, token.TEXT, scan.tok) // error is recoverable
	if assert.NotNil(t, err) {
		assert.Equal(t, 20, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 21, err.pos.Column)
		assert.Equal(t, `invalid character in text literal: U+0007`, err.msg)
	}

	scan, err = scanOnce("\"null-byte \u0000 in string\"")
	assert.Equal(t, token.TEXT, scan.tok) // error is recoverable
	if assert.NotNil(t, err) {
		assert.Equal(t, 11, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 12, err.pos.Column)
		assert.Equal(t, `invalid character in text literal: U+0000`, err.msg)
	}

	scan, err = scanOnce(`"\u`)
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `unterminated text literal`, err.msg)
		if assert.NotNil(t, err.prev) {
			err = err.prev
			assert.Equal(t, 3, err.pos.Offset)
			assert.Equal(t, 1, err.pos.Line)
			assert.Equal(t, 4, err.pos.Column)
			assert.Equal(t, `unterminated escape sequence`, err.msg)
		}
	}

	scan, err = scanOnce(`"\`)
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `unterminated text literal`, err.msg)
		if assert.NotNil(t, err.prev) {
			err = err.prev
			assert.Equal(t, 2, err.pos.Offset)
			assert.Equal(t, 1, err.pos.Line)
			assert.Equal(t, 3, err.pos.Column)
			assert.Equal(t, `unterminated escape sequence`, err.msg)
		}
	}

	scan, err = scanOnce(`"\m"`)
	assert.Equal(t, token.TEXT, scan.tok) // error is recoverable
	if assert.NotNil(t, err) {
		assert.Equal(t, 2, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 3, err.pos.Column)
		assert.Equal(t, `unknown escape sequence`, err.msg)
	}

	scan, err = scanOnce(`"\uD800"`)
	assert.Equal(t, token.TEXT, scan.tok) // error is recoverable
	if assert.NotNil(t, err) {
		assert.Equal(t, 2, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 3, err.pos.Column)
		assert.Equal(t, `escape sequence is invalid Unicode code point`, err.msg)
	}

	scan, err = scanOnce("\"multi\nline\"")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `unterminated text literal`, err.msg)
	}

	scan, err = scanOnce("\"multi\rline\"")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `unterminated text literal`, err.msg)
	}

	scan, err = scanOnce(`"bad \u1 esc"`)
	assert.Equal(t, token.TEXT, scan.tok) // error is recoverable
	if assert.NotNil(t, err) {
		assert.Equal(t, 8, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 9, err.pos.Column)
		assert.Equal(t, `unexpected character in escape sequence: U+0020 ' '`, err.msg)
	}

	scan, err = scanOnce(`"bad \u0XX1 esc"`)
	assert.Equal(t, token.TEXT, scan.tok) // error is recoverable
	if assert.NotNil(t, err) {
		assert.Equal(t, 8, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 9, err.pos.Column)
		assert.Equal(t, `unexpected character in escape sequence: U+0058 'X'`, err.msg)
	}

	scan, err = scanOnce(`"bad \uXXXX esc"`)
	assert.Equal(t, token.TEXT, scan.tok) // error is recoverable
	if assert.NotNil(t, err) {
		assert.Equal(t, 7, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 8, err.pos.Column)
		assert.Equal(t, `unexpected character in escape sequence: U+0058 'X'`, err.msg)
	}

	scan, err = scanOnce(`"bad \uFXXX esc"`)
	assert.Equal(t, token.TEXT, scan.tok) // error is recoverable
	if assert.NotNil(t, err) {
		assert.Equal(t, 8, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 9, err.pos.Column)
		assert.Equal(t, `unexpected character in escape sequence: U+0058 'X'`, err.msg)
	}

	scan, err = scanOnce(`"bad \uXXXF esc"`)
	assert.Equal(t, token.TEXT, scan.tok) // error is recoverable
	if assert.NotNil(t, err) {
		assert.Equal(t, 7, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 8, err.pos.Column)
		assert.Equal(t, `unexpected character in escape sequence: U+0058 'X'`, err.msg)
	}
}

func TestScansNumbers(t *testing.T) {
	scan, err := scanOnce("4")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "4", scan.lit)

	scan, err = scanOnce("4.123")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "4.123", scan.lit)

	scan, err = scanOnce(".4")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ".4", scan.lit)

	scan, err = scanOnce(".123")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ".123", scan.lit)

	scan, err = scanOnce("9")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "9", scan.lit)

	scan, err = scanOnce("0")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "0", scan.lit)

	scan, err = scanOnce("0.123")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "0.123", scan.lit)

	scan, err = scanOnce("123e4")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "123e4", scan.lit)

	scan, err = scanOnce("123e-4")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "123e-4", scan.lit)

	scan, err = scanOnce("123e+4")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "123e+4", scan.lit)

	scan, err = scanOnce(".123e4")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ".123e4", scan.lit)

	scan, err = scanOnce(".123e-4")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ".123e-4", scan.lit)

	scan, err = scanOnce(".123e+4")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ".123e+4", scan.lit)

	scan, err = scanOnce(".123e4567")
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ".123e4567", scan.lit)
}

func TestReportsUsefulNumberErrors(t *testing.T) {
	scan, err := scanOnce("1.")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `missing digits after decimal point in number`, err.msg)
	}

	scan, err = scanOnce("1.A")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 2, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 3, err.pos.Column)
		assert.Equal(t, `missing space after number before "A"`, err.msg)
		if assert.NotNil(t, err.prev) {
			err = err.prev
			assert.Equal(t, 0, err.pos.Offset)
			assert.Equal(t, 1, err.pos.Line)
			assert.Equal(t, 1, err.pos.Column)
			assert.Equal(t, `missing digits after decimal point in number`, err.msg)
		}
	}

	scan, err = scanOnce("1.0e")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `missing digits after exponent in number`, err.msg)
	}

	scan, err = scanOnce("1.0eA")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 4, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 5, err.pos.Column)
		assert.Equal(t, `missing space after number before "A"`, err.msg)
		if assert.NotNil(t, err.prev) {
			err = err.prev
			assert.Equal(t, 0, err.pos.Offset)
			assert.Equal(t, 1, err.pos.Line)
			assert.Equal(t, 1, err.pos.Column)
			assert.Equal(t, `missing digits after exponent in number`, err.msg)
		}
	}
}

func TestScansOperators(t *testing.T) {
	scan, err := scanOnce(`_dot_`)
	assert.Nil(t, err)
	assert.Equal(t, token.OPERATOR, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `_dot_`, scan.lit)

	scan, err = scanOnce(`_cross_`)
	assert.Nil(t, err)
	assert.Equal(t, token.OPERATOR, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `_cross_`, scan.lit)

	scan, err = scanOnce(`_seconds`)
	assert.Nil(t, err)
	assert.Equal(t, token.OPERATOR, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, `_seconds`, scan.lit)

	scan, err = scanOnce("*")
	assert.Nil(t, err)
	assert.Equal(t, token.OPERATOR, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "*", scan.lit)

	scan, err = scanOnce("/")
	assert.Nil(t, err)
	assert.Equal(t, token.OPERATOR, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "/", scan.lit)

	scan, err = scanOnce("+")
	assert.Nil(t, err)
	assert.Equal(t, token.OPERATOR, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "+", scan.lit)

	scan, err = scanOnce("-")
	assert.Nil(t, err)
	assert.Equal(t, token.OPERATOR, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "-", scan.lit)

	scan, err = scanOnce(".")
	assert.Nil(t, err)
	assert.Equal(t, token.PERIOD, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ".", scan.lit)

	scan, err = scanOnce(":")
	assert.Nil(t, err)
	assert.Equal(t, token.COLON, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ":", scan.lit)

	scan, err = scanOnce("::")
	assert.Nil(t, err)
	assert.Equal(t, token.CONS, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "::", scan.lit)

	scan, err = scanOnce(":::")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `too many colons for '::'`, err.msg)
	}

	scan, err = scanOnce(";")
	assert.Nil(t, err)
	assert.Equal(t, token.SEMICOLON, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ";", scan.lit)

	scan, err = scanOnce(",")
	assert.Nil(t, err)
	assert.Equal(t, token.COMMA, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ",", scan.lit)

	scan, err = scanOnce("=")
	assert.Nil(t, err)
	assert.Equal(t, token.EQUALS, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "=", scan.lit)

	scan, err = scanOnce("->")
	assert.Nil(t, err)
	assert.Equal(t, token.ARROW, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "->", scan.lit)
}

func TestScansDelimiters(t *testing.T) {
	scan, err := scanOnce("(")
	assert.Nil(t, err)
	assert.Equal(t, token.LEFT_PAREN, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "(", scan.lit)

	scan, err = scanOnce(")")
	assert.Nil(t, err)
	assert.Equal(t, token.RIGHT_PAREN, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, ")", scan.lit)

	scan, err = scanOnce("[")
	assert.Nil(t, err)
	assert.Equal(t, token.LEFT_BRACKET, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "[", scan.lit)

	scan, err = scanOnce("]")
	assert.Nil(t, err)
	assert.Equal(t, token.RIGHT_BRACKET, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "]", scan.lit)

	scan, err = scanOnce("{")
	assert.Nil(t, err)
	assert.Equal(t, token.LEFT_BRACE, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "{", scan.lit)

	scan, err = scanOnce("}")
	assert.Nil(t, err)
	assert.Equal(t, token.RIGHT_BRACE, scan.tok)
	assert.Equal(t, 0, scan.pos)
	assert.Equal(t, "}", scan.lit)
}

func TestReportsUsefulUnknownCharacter(t *testing.T) {
	scan, err := scanOnce("\u203B")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, "unexpected character U+203B '\u203B'", err.msg)
	}

	scan, err = scanOnce("\u200b")
	assert.Equal(t, token.INVALID, scan.tok)
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `unexpected character U+200B`, err.msg)
	}
}

func TestScannerNextCharacter(t *testing.T) {
	var err *scanError

	err = scanAll("SELECT * FROM candies\r\n  WHERE sweetness = 11\n\r\r")
	assert.Nil(t, err)

	err = scanAll(string([]byte{0x00, 0xFF}))
	if assert.NotNil(t, err) {
		assert.Equal(t, 0, err.pos.Offset)
		assert.Equal(t, 1, err.pos.Line)
		assert.Equal(t, 1, err.pos.Column)
		assert.Equal(t, `invalid character U+0000`, err.msg)
	}
}

func TestScanPos(t *testing.T) {
	var err *scanError
	handleError := func(pos token.Position, msg string) {
		err = &scanError{pos, msg, nil}
	}

	var scan scanToken
	s := Scanner{}
	s.Init("main", []byte("Main := () {\n  exit(1)\n}"), handleError)
	assert.Equal(t, token.Position{"main", 0, 1, 1}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 4, 1, 5}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 6, 1, 7}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 7, 1, 8}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 9, 1, 10}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 10, 1, 11}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 12, 1, 13}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 19, 2, 7}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 20, 2, 8}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 21, 2, 9}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 22, 2, 10}, s.Pos())
	assert.Nil(t, err)

	scan.pos, scan.tok, scan.lit = s.Scan()
	assert.Equal(t, token.Position{"main", 24, 3, 2}, s.Pos())
	assert.Nil(t, err)
}
