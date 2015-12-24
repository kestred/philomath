package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionIsValid(t *testing.T) {
	var pos Position

	pos = Position{"", 0, 0, 0}
	assert.False(t, pos.IsValid())

	pos = Position{"", 0, 1, 0}
	assert.True(t, pos.IsValid())
}

func TestPositionString(t *testing.T) {
	var pos Position

	pos = Position{"", 0, 0, 0}
	assert.Equal(t, "-", pos.String())

	pos = Position{"Src", 0, 1, 1}
	assert.Equal(t, "Src:1:1", pos.String())

	pos = Position{"Name", 15, 7, 16}
	assert.Equal(t, "Name:7:16", pos.String())
}

func TestLookup(t *testing.T) {
	// an arbitrary string
	assert.Equal(t, IDENT, Lookup("something"))

	// a manual example
	assert.Equal(t, STRUCT, Lookup("struct"))

	// all keyword tokens and no non-keyword tokens
	for i, name := range tokens {
		tok := Token(i)
		if tok.IsKeyword() {
			assert.Equal(t, tok, Lookup(name))
		} else {
			assert.Equal(t, IDENT, Lookup(name))
		}
	}
}

func TestTokenString(t *testing.T) {
	assert.Equal(t, "Invalid token", INVALID.String())
	assert.Equal(t, "End of source", END.String())

	assert.Equal(t, "Identifier", IDENT.String())

	assert.Equal(t, "Integer", INTEGER.String())
	assert.Equal(t, "Real", REAL.String())
	assert.Equal(t, "Text", TEXT.String())

	assert.Equal(t, ":", COLON.String())
	assert.Equal(t, "::", CONS.String())
	assert.Equal(t, "=", EQUALS.String())
	assert.Equal(t, "*", ASTERISK.String())
	assert.Equal(t, "/", SLASH.String())
	assert.Equal(t, "+", PLUS.String())
	assert.Equal(t, "-", MINUS.String())
	assert.Equal(t, ",", COMMA.String())
	assert.Equal(t, ".", PERIOD.String())

	assert.Equal(t, "(", LEFT_PAREN.String())
	assert.Equal(t, "[", LEFT_BRACKET.String())
	assert.Equal(t, ")", RIGHT_PAREN.String())
	assert.Equal(t, "]", RIGHT_BRACKET.String())

	assert.Equal(t, "(", LEFT_PAREN.String())
	assert.Equal(t, "[", LEFT_BRACKET.String())
	assert.Equal(t, "{", LEFT_BRACE.String())
	assert.Equal(t, ")", RIGHT_PAREN.String())
	assert.Equal(t, "]", RIGHT_BRACKET.String())
	assert.Equal(t, "}", RIGHT_BRACE.String())

	assert.Equal(t, "Token(2000)", Token(2000).String())
}

func TestHasLiteral(t *testing.T) {
	assert.Equal(t, false, INVALID.HasLiteral())
	assert.Equal(t, false, END.HasLiteral())

	assert.Equal(t, true, IDENT.HasLiteral())

	assert.Equal(t, true, NUMBER.HasLiteral())
	assert.Equal(t, true, TEXT.HasLiteral())

	assert.Equal(t, false, COLON.HasLiteral())
	assert.Equal(t, false, CONS.HasLiteral())
	assert.Equal(t, false, EQUALS.HasLiteral())
	assert.Equal(t, false, ASTERISK.HasLiteral())
	assert.Equal(t, false, SLASH.HasLiteral())
	assert.Equal(t, false, PLUS.HasLiteral())
	assert.Equal(t, false, MINUS.HasLiteral())
	assert.Equal(t, false, COMMA.HasLiteral())
	assert.Equal(t, false, PERIOD.HasLiteral())
	assert.Equal(t, false, ARROW.HasLiteral())

	assert.Equal(t, false, LEFT_PAREN.HasLiteral())
	assert.Equal(t, false, LEFT_BRACKET.HasLiteral())
	assert.Equal(t, false, LEFT_BRACE.HasLiteral())
	assert.Equal(t, false, RIGHT_PAREN.HasLiteral())
	assert.Equal(t, false, RIGHT_BRACKET.HasLiteral())
	assert.Equal(t, false, RIGHT_BRACE.HasLiteral())

	assert.Equal(t, true, IF.HasLiteral())
	assert.Equal(t, true, FOR.HasLiteral())
	assert.Equal(t, true, IN.HasLiteral())
	assert.Equal(t, true, BREAK.HasLiteral())
	assert.Equal(t, true, RETURN.HasLiteral())

	assert.Equal(t, true, STRUCT.HasLiteral())
	assert.Equal(t, true, MODULE.HasLiteral())
}

func TestIsKeyword(t *testing.T) {
	assert.Equal(t, false, INVALID.IsKeyword())
	assert.Equal(t, false, END.IsKeyword())

	assert.Equal(t, false, IDENT.IsKeyword())

	assert.Equal(t, false, NUMBER.IsKeyword())
	assert.Equal(t, false, TEXT.IsKeyword())

	assert.Equal(t, false, COLON.IsKeyword())
	assert.Equal(t, false, CONS.IsKeyword())
	assert.Equal(t, false, EQUALS.IsKeyword())
	assert.Equal(t, false, ASTERISK.IsKeyword())
	assert.Equal(t, false, SLASH.IsKeyword())
	assert.Equal(t, false, PLUS.IsKeyword())
	assert.Equal(t, false, MINUS.IsKeyword())
	assert.Equal(t, false, COMMA.IsKeyword())
	assert.Equal(t, false, PERIOD.IsKeyword())
	assert.Equal(t, false, ARROW.IsKeyword())

	assert.Equal(t, false, LEFT_PAREN.IsKeyword())
	assert.Equal(t, false, LEFT_BRACKET.IsKeyword())
	assert.Equal(t, false, LEFT_BRACE.IsKeyword())
	assert.Equal(t, false, RIGHT_PAREN.IsKeyword())
	assert.Equal(t, false, RIGHT_BRACKET.IsKeyword())
	assert.Equal(t, false, RIGHT_BRACE.IsKeyword())

	assert.Equal(t, true, IF.IsKeyword())
	assert.Equal(t, true, FOR.IsKeyword())
	assert.Equal(t, true, IN.IsKeyword())
	assert.Equal(t, true, BREAK.IsKeyword())
	assert.Equal(t, true, RETURN.IsKeyword())

	assert.Equal(t, true, STRUCT.IsKeyword())
	assert.Equal(t, true, MODULE.IsKeyword())
}
