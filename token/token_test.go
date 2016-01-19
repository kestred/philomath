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

	assert.Equal(t, "Number", NUMBER.String())
	assert.Equal(t, "Text", TEXT.String())

	assert.Equal(t, "Operator", OPERATOR.String())
	assert.Equal(t, ".", PERIOD.String())

	assert.Equal(t, ":", COLON.String())
	assert.Equal(t, "::", CONS.String())
	assert.Equal(t, ";", SEMICOLON.String())
	assert.Equal(t, ",", COMMA.String())
	assert.Equal(t, "=", EQUALS.String())
	assert.Equal(t, "->", ARROW.String())

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

func TestIsOperator(t *testing.T) {
	assert.Equal(t, false, INVALID.IsOperator())
	assert.Equal(t, false, END.IsOperator())

	assert.Equal(t, false, IDENT.IsOperator())

	assert.Equal(t, false, NUMBER.IsOperator())
	assert.Equal(t, false, TEXT.IsOperator())

	assert.Equal(t, true, OPERATOR.IsOperator())
	assert.Equal(t, true, PERIOD.IsOperator())

	assert.Equal(t, false, COLON.IsOperator())
	assert.Equal(t, false, CONS.IsOperator())
	assert.Equal(t, false, SEMICOLON.IsOperator())
	assert.Equal(t, false, COMMA.IsOperator())
	assert.Equal(t, false, EQUALS.IsOperator())
	assert.Equal(t, false, ARROW.IsOperator())

	assert.Equal(t, false, LEFT_PAREN.IsOperator())
	assert.Equal(t, false, LEFT_BRACKET.IsOperator())
	assert.Equal(t, false, LEFT_BRACE.IsOperator())
	assert.Equal(t, false, RIGHT_PAREN.IsOperator())
	assert.Equal(t, false, RIGHT_BRACKET.IsOperator())
	assert.Equal(t, false, RIGHT_BRACE.IsOperator())

	assert.Equal(t, false, IF.IsOperator())
	assert.Equal(t, false, FOR.IsOperator())
	assert.Equal(t, false, IN.IsOperator())
	assert.Equal(t, false, BREAK.IsOperator())
	assert.Equal(t, false, RETURN.IsOperator())

	assert.Equal(t, false, STRUCT.IsOperator())
	assert.Equal(t, false, MODULE.IsOperator())
}

func TestIsKeyword(t *testing.T) {
	assert.Equal(t, false, INVALID.IsKeyword())
	assert.Equal(t, false, END.IsKeyword())

	assert.Equal(t, false, IDENT.IsKeyword())

	assert.Equal(t, false, NUMBER.IsKeyword())
	assert.Equal(t, false, TEXT.IsKeyword())

	assert.Equal(t, false, OPERATOR.IsKeyword())
	assert.Equal(t, false, PERIOD.IsKeyword())

	assert.Equal(t, false, COLON.IsKeyword())
	assert.Equal(t, false, CONS.IsKeyword())
	assert.Equal(t, false, SEMICOLON.IsKeyword())
	assert.Equal(t, false, COMMA.IsKeyword())
	assert.Equal(t, false, EQUALS.IsKeyword())
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
