package scanner

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/kestred/philomath/code/token"
	"github.com/kestred/philomath/code/utils"
)

const bom = 0xFEFF // byte order mark, only permitted as very first character

func isValid(ch rune) bool {
	return (ch >= '\u0020' ||
		ch == '\u0009' ||
		ch == '\u000A' ||
		ch == '\u000D')
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// An ErrorHandler may be provided to Scanner.Init. If a syntax error is
// encountered and a handler was installed, the handler is called with a
// position and an error message. The position points to the beginning of
// the offending token.
type ErrorHandler func(pos token.Position, msg string)

// A Scanner holds the scanner's internal state while processing
// a given text.  It can be allocated as part of another data
// structure but must be initialized via Init before use.
type Scanner struct {
	// immutable state
	src      []byte
	err      ErrorHandler
	filename string

	// scanning state
	char       rune         // current character
	offset     int          // byte offset to current char
	readOffset int          // reading offset (position after current character)
	lineOffset int          // current line offset
	line       int          // current line number
	lines      []token.Line // previous line numbers and offsets
}

// Init prepares the scanner s to tokenize the text src by setting the
// scanner at the beginning of src.
//
// Calls to Scan will invoke the error handler err if they encounter a
// syntax error and err is not nil.
//
// Init may call err if there is an error in the first character of the file.
func (s *Scanner) Init(filename string, src []byte, err ErrorHandler) {
	s.src = src
	s.err = err
	s.filename = filename

	s.char = ' '
	s.offset = 0
	s.readOffset = 0
	s.lineOffset = 0
	s.line = 1

	s.next()
	if s.char == bom {
		s.next()
	}
}

// Peek pretends to scan the next token.  Calling Peek then Scan is a little
// over twice as much effort as just Scan, so avoid calling Peek when possible.
// The error handler will not be called while peeking.
//
// Added to support LL(2) parsing of declarations vs statements in block.
// All other rules are LL(1) parsable, so this should be called infrequently.
func (s *Scanner) Peek() (tok token.Token) {
	// save scanner state
	char := s.char
	offset := s.offset
	readOffset := s.readOffset
	lineOffset := s.lineOffset
	line := s.line
	lines := s.lines
	err := s.err
	s.err = nil // ignore errors

	_, tok, _ = s.Scan()

	// restore scanner state
	s.char = char
	s.offset = offset
	s.readOffset = readOffset
	s.lineOffset = lineOffset
	s.line = line
	s.lines = lines
	s.err = err

	return
}

// Scan scans the next token and returns the token position, the token, and its
// literal string. The source end is indicated by the END token.
func (s *Scanner) Scan() (pos int, tok token.Token, lit string) {
scanAgain:
	s.skipWhitespace()

	pos = s.offset
	ch := s.char
	switch {
	case ch == '_':
		lit = s.scanOperator()
		tok = token.OPERATOR
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = token.IDENT
		if len(lit) > 1 {
			tok = token.Lookup(lit)
		}
	case isDigit(ch):
		tok, lit = s.scanNumber(false)
	default:
		s.next() // always make progress
		lit = string(ch)

		switch ch {
		case -1:
			tok = token.END
		case '#':
			// NOTE: Ignore leading "#!/usr/bin/philomath run" et. al.
			if s.char == '!' {
				// TODO: Error unless the shebang includes "run"
				// TODO: Error unless the shebang is on the first line
				s.next()
				if s.scanComment(false) {
					goto scanAgain
				} else {
					tok = token.INVALID
				}
			} else if isLetter(s.char) {
				lit = s.scanIdentifier()
				tok = token.DIRECTIVE
			} else {
				s.error(pos, fmt.Sprintf("unexpected character %#U", ch))
				tok = token.INVALID
			}
			// TODO: Scan directive
		case '"':
			tok, lit = s.scanText()
		case '.':
			if isDigit(s.char) {
				tok, lit = s.scanNumber(true)
			} else {
				tok = token.PERIOD
			}
		case '+', '*', '%':
			tok = token.OPERATOR
		case '/':
			if s.char == '/' {
				s.next()
				if s.scanComment(false) {
					goto scanAgain
				} else {
					tok = token.INVALID
				}
			} else if s.char == '*' {
				s.next()
				if s.scanComment(true) {
					goto scanAgain
				} else {
					tok = token.INVALID
				}
			} else {
				tok = token.OPERATOR
			}
		case '-':
			if s.char == '>' {
				s.next()
				tok = token.ARROW
				lit = "->"
			} else {
				tok = token.OPERATOR
			}
		case ':':
			numColons := 1
			for s.char == ':' {
				numColons += 1
				s.next()
			}

			switch numColons {
			case 1:
				tok = token.COLON
			case 2:
				tok = token.CONS
				lit = "::"
			default:
				s.error(pos, "too many colons for '::'")
				tok = token.INVALID
				lit = string(s.src[pos:s.offset])
			}
		case ';':
			tok = token.SEMICOLON
		case ',':
			tok = token.COMMA
		case '=':
			tok = token.EQUALS
		case '(':
			tok = token.LEFT_PAREN
		case '[':
			tok = token.LEFT_BRACKET
		case '{':
			tok = token.LEFT_BRACE
		case ')':
			tok = token.RIGHT_PAREN
		case ']':
			tok = token.RIGHT_BRACKET
		case '}':
			tok = token.RIGHT_BRACE
		default:
			if isValid(ch) {
				s.error(pos, fmt.Sprintf("unexpected character %#U", ch))
			} else if ch != bom { // next reports unexpected BOMs - don't repeat
				s.error(pos, fmt.Sprintf("invalid character %#U", ch))
			}
			tok = token.INVALID
		}
	}

	return
}

func (s *Scanner) Pos() token.Position {
	// Get length of current line in UTF-8 characters
	column := 1 + utf8.RuneCount(s.src[s.lineOffset:s.offset])
	return token.Position{
		Name:   s.filename,
		Offset: s.offset,
		Line:   s.line,
		Column: column,
	}
}

func (s *Scanner) LineAt(offset int) token.Line {
	if offset >= s.lineOffset {
		return token.Line{s.line, s.lineOffset, string(s.src[s.lineOffset:offset])}
	}

	var result token.Line
	for i := len(s.lines) - 1; i >= 0; i-- {
		if s.lines[i].Offset <= offset {
			result = s.lines[i]
			break
		}
	}

	utils.Assert(result.IsValid(), fmt.Sprintf("Failed to find line info at file[%s]:%d", s.filename, offset))
	return result
}

func (s *Scanner) SourceAt(start, end int) string {
	return string(s.src[start : end+1])
}

func (s *Scanner) error(offset int, msg string) {
	if s.err != nil {
		line := s.LineAt(offset)
		column := 1 + utf8.RuneCount(s.src[line.Offset:offset])
		pos := token.Position{
			Name:   s.filename,
			Offset: offset,
			Line:   line.Line,
			Column: column,
		}

		s.err(pos, msg)
	}
}

func (s *Scanner) next() {
	if s.readOffset < len(s.src) {
		s.offset = s.readOffset

		wasCarriageReturn := false
		if s.char == '\n' {
			line := token.Line{s.line, s.lineOffset, string(s.src[s.lineOffset:s.offset])}
			s.lines = append(s.lines, line)
			s.line += 1
			s.lineOffset = s.offset
		} else if s.char == '\r' {
			wasCarriageReturn = true
			line := token.Line{s.line, s.lineOffset, string(s.src[s.lineOffset:s.offset])}
			s.lines = append(s.lines, line)
			s.line += 1
			s.lineOffset = s.offset
		}

		r, width := rune(s.src[s.readOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "unexpected character: U+0000")
		case r >= 0x80:
			// not ASCII
			r, width = utf8.DecodeRune(s.src[s.readOffset:])
			if r == utf8.RuneError && width == 1 {
				s.error(s.offset, "invalid UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "invalid byte order mark")
			}
		}
		s.readOffset += width
		s.char = r

		if s.char == '\n' && wasCarriageReturn {
			s.line -= 1
		}
	} else {
		s.offset = len(s.src)
		if s.char == '\n' || s.char == '\r' {
			s.lineOffset = s.offset
		}
		s.char = -1 // eof
	}
}

func (s *Scanner) skipWhitespace() {
	for s.char == ' ' || s.char == '\t' || s.char == '\n' || s.char == '\r' {
		s.next()
	}
}

func (s *Scanner) scanOperator() string {
	offset := s.offset
	for isLetter(s.char) || s.char == '_' {
		s.next()
	}

	return string(s.src[offset:s.offset])
}

func (s *Scanner) scanIdentifier() string {
	offset := s.offset
	for isLetter(s.char) || isDigit(s.char) || s.char == '_' {
		s.next()
	}

	return string(s.src[offset:s.offset])
}

// scanComment eats a (line or block) comment.
// In case of an error, it returns false. Otherwise it returns true.
func (s *Scanner) scanComment(block bool) bool {
	// initial "//" or "/*" already consumed
	offset := s.offset - 2
	if block {
		/* block comment */
		s.next()

		nesting := 0
		for s.char >= 0 {
			ch := s.char
			s.next()

			if ch == '/' {
				if s.char == '*' {
					s.next()
					nesting += 1
				}
			} else if ch == '*' {
				if s.char == '/' {
					s.next()
					if nesting > 0 {
						nesting -= 1
					} else {
						break
					}
				}
			}
		}
	} else {
		// // line comment
		for (s.char != '\n' && s.char != '\r') && s.char >= 0 {
			s.next()
		}
	}

	if s.char < 0 {
		s.error(offset, "unterminated comment")
		return false
	}

	return true
}

func (s *Scanner) scanMantissa() {
	for isDigit(s.char) {
		s.next()
	}
}

func (s *Scanner) scanNumber(afterDecimal bool) (token.Token, string) {
	tok := token.NUMBER
	offset := s.offset
	likeNumber := false
	if afterDecimal {
		offset -= 1
		likeNumber = true
	}

	s.scanMantissa()
	if s.char == '.' && !afterDecimal { // TODO: maybe an error?
		likeNumber = true
		s.next()

		decOffset := s.offset
		s.scanMantissa()
		if s.offset == decOffset {
			s.error(offset, "missing digits after decimal point in number")
			tok = token.INVALID
		}
	}
	if s.char == 'e' {
		likeNumber = true
		s.next()

		if s.char == '+' || s.char == '-' {
			s.next()
		}
		expOffset := s.offset
		s.scanMantissa()
		if s.offset == expOffset {
			s.error(offset, "missing digits after exponent in number")
			tok = token.INVALID
		}
	}

	if isLetter(s.char) {
		charOffset := s.offset
		for isLetter(s.char) || s.char == '_' {
			s.next()
		}

		// check if it seems like a number (exponent, decimal point, etc)
		// also, if it is a long number (say, longer than a year), guess it is a number
		if likeNumber || (charOffset-offset > 4) {
			name := string(s.src[charOffset:s.offset])
			msg := fmt.Sprintf(`missing space after number before "%s"`, name)
			s.error(charOffset, msg)
		} else {
			s.error(offset, "identifier cannot start with digits (eg. 0, 1, ..., 9)")
		}
		tok = token.INVALID
	}

	return tok, string(s.src[offset:s.offset])
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}

// scanEscape parses an escape sequence where rune is the accepted
// escaped quote. In case of a syntax error, it stops at the offending
// character (without consuming it) and returns false. Otherwise
// it returns true.
func (s *Scanner) scanEscape(quote rune) bool {
	offset := s.offset

	var n int
	var base, max uint32
	switch s.char {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.next()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.next()
		n, base, max = 2, 16, 255
	case 'u':
		s.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if s.char < 0 {
			msg = "unterminated escape sequence"
		}
		s.error(offset, msg)
		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(s.char))
		if d >= base {
			msg := fmt.Sprintf("unexpected character in escape sequence: %#U", s.char)
			if s.char < 0 {
				msg = "unterminated escape sequence"
			}
			s.error(s.offset, msg)
			return false
		}
		x = x*base + d
		s.next()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		s.error(offset, "escape sequence is invalid Unicode code point")
		return false
	}

	return true
}

func (s *Scanner) scanText() (token.Token, string) {
	// opening quote already consumed
	offset := s.offset - 1
	tok := token.TEXT

	for {
		ch := s.char
		if ch == '\n' || ch == '\r' || ch < 0 {
			tok = token.INVALID
			s.error(offset, "unterminated text literal")
			break
		} else if !isValid(ch) {
			msg := fmt.Sprintf("invalid character in text literal: %#U", ch)
			s.error(s.offset, msg)
		}
		s.next()

		if ch == '"' {
			break
		}
		if ch == '\\' {
			s.scanEscape('"')
		}
	}

	return tok, string(s.src[offset:s.offset])
}
