package lexer

import (
	"fmt"
	token2 "mintsql/internal/sql/token"
	"strings"
	"unicode"
	"unicode/utf8"
)

type LexFn func(*Lexer) LexFn

type Cursor struct {
	Start int
	Pos   int
	Width int
}

type Lexer struct {
	Name   string
	Input  string
	Tokens chan token2.Token
	State  LexFn

	Cursor
	token2.Location
}

func (l *Lexer) Emit(kind token2.Kind) {
	l.Tokens <- token2.Token{
		Kind:     kind,
		Value:    l.Input[l.Start:l.Pos],
		Location: l.Location,
	}
	l.Start = l.Pos
}

func (l *Lexer) Inc() {
	l.Pos++
	l.Column++
	if l.Pos >= len(l.Input) {
		l.Emit(token2.KindEof)
	}
}

func (l *Lexer) Dec() {
	l.Pos--
	l.Column--
}

func (l *Lexer) Backup() {
	l.Pos -= l.Width
	l.Column -= l.Width
}

func (l *Lexer) Current() string {
	return l.Input[l.Start:l.Pos]
}

func (l *Lexer) Remainder() string {
	return l.Input[l.Pos:]
}

func (l *Lexer) Next() rune {
	if l.Pos >= len(l.Input) {
		l.Width = 0
		return token2.EOF
	}

	res, width := utf8.DecodeRuneInString(l.Input[l.Pos:])
	l.Width = width
	l.Pos += l.Width
	l.Column++
	if res == token2.Newline {
		l.Line++
		l.Column = 0
	}
	return res
}

func (l *Lexer) Peek() rune {
	c := l.Next()
	l.Backup()
	return c
}

func (l *Lexer) Ignore() {
	l.Start = l.Pos
}

func (l *Lexer) IsEOF() bool {
	return l.Pos >= len(l.Input)
}

func (l *Lexer) IsWhiteSpace() bool {
	c, _ := utf8.DecodeRuneInString(l.Input[l.Pos:])
	return unicode.IsSpace(c)
}

func (l *Lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

func (l *Lexer) AcceptMany(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
}

func (l *Lexer) IsNewLine() bool {
	c, _ := utf8.DecodeRuneInString(l.Input[l.Pos : l.Pos+1])
	return c == token2.Newline
}

func (l *Lexer) SkipNewLine() {
	for {
		c := l.Next()
		if c == token2.EOF {
			l.Emit(token2.KindEof)
			break
		}
		if c == token2.Newline {
			l.Line++
			l.Column = 0
		} else {
			l.Dec()
			break
		}
	}
}

func (l *Lexer) SkipWhiteSpace() {
	for {
		c := l.Next()
		if c == token2.EOF {
			l.Emit(token2.KindEof)
			break
		}
		if !unicode.IsSpace(c) {
			l.Dec()
			break
		}
	}
}

func (l *Lexer) Errorf(format string, args ...interface{}) LexFn {
	l.Tokens <- token2.Token{
		Kind:     token2.KindError,
		Value:    fmt.Sprintf(format, args...),
		Location: l.Location,
	}
	return nil
}

func (l *Lexer) Shutdown() {
	close(l.Tokens)
}

func (l *Lexer) Run(init LexFn) {
	defer l.Shutdown()
	for st := init; st != nil; {
		st = st(l)
	}
}
