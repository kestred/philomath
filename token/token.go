package token

import (
	"fmt"
	"strconv"
)

// Position describes an arbitrary source position including the name, line,
// and column location. A Position is valid if the line number is > 0.
type Position struct {
	Name   string // source name, if any
	Offset int    // offset, starting at 0
	Line   int    // line number, starting at 1
	Column int    // column number, starting at 1
}

// IsValid reports whether the position is valid.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

// String returns a string in one of several forms:
//
//	name:line:column    valid position with name
//	line:column         valid position without name
//	name                invalid position with name
//	-                   invalid position without name
//
func (pos Position) String() string {
	s := pos.Name
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	}
	if s == "" {
		s = "-"
	}
	return s
}

// Token is the set of lexical tokens in Philomath
type Token int

const (
	// Special tokens
	INVALID Token = iota
	END

	// Identifier
	IDENT

	// Literals
	NUMBER
	TEXT

	// Operators
	COLON    // :
	CONS     // ::
	EQUALS   // =
	ASTERISK // *
	SLASH    // /
	PLUS     // +
	HYPHEN   // -
	ARROW    // ->
	COMMA    // ,
	PERIOD   // .

	// Delimiters
	LEFT_PAREN    // (
	LEFT_BRACKET  // [
	LEFT_BRACE    // {
	RIGHT_PAREN   // )
	RIGHT_BRACKET // ]
	RIGHT_BRACE   // }

	// Keywords
	keywords_begin

	IF     // if
	FOR    // for
	IN     // in
	BREAK  // break
	RETURN // return

	STRUCT // struct
	MODULE // module

	keywords_end
)

var tokens = [...]string{
	INVALID: "Invalid token",
	END:     "End of source",

	IDENT: "Identifier",

	NUMBER: "Number",
	TEXT:   "Text",

	COLON:    ":",
	CONS:     "::",
	EQUALS:   "=",
	ASTERISK: "*",
	SLASH:    "/",
	PLUS:     "+",
	HYPHEN:   "-",
	ARROW:    "->",
	COMMA:    ",",
	PERIOD:   ".",

	LEFT_PAREN:    "(",
	LEFT_BRACKET:  "[",
	LEFT_BRACE:    "{",
	RIGHT_PAREN:   ")",
	RIGHT_BRACKET: "]",
	RIGHT_BRACE:   "}",

	IF:     "if",
	FOR:    "for",
	IN:     "in",
	BREAK:  "break",
	RETURN: "return",

	STRUCT: "struct",
	MODULE: "module",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "Token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keywords_begin + 1; i < keywords_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword)
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

func (tok Token) HasLiteral() bool {
	return (IDENT <= tok && tok <= TEXT) || tok.IsKeyword()
}

func (tok Token) IsKeyword() bool {
	return keywords_begin < tok && tok < keywords_end
}
