package lexer

import (
	"mintsql/internal/sql/token"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	Newline rune = 10
	EOF     rune = 0
)

type LexFn func(*Lexer) LexFn

type Cursor struct {
	Start          int
	Pos            int
	Width          int
	LastLineLength int
}

type Lexer struct {
	Input  string
	Tokens chan *token.Token
	State  LexFn

	Cursor
	token.Location
}

func (l *Lexer) Emit(kind token.Kind) {
	newToken := &token.Token{
		Kind:     kind,
		Value:    l.Input[l.Start:l.Pos],
		Location: l.Location,
	}
	if kind == token.KindEof {
		newToken.Value = "EOF"
	}
	l.Tokens <- newToken
	l.Start = l.Pos
}

func (l *Lexer) Current() string {
	return l.Input[l.Start:l.Pos]
}

func (l *Lexer) Remainder() string {
	return l.Input[l.Pos:]
}

//func (l *Lexer) Revert() {
//	l.Pos = l.Start
//	l.Width = 0
//	l.Column
//}

func (l *Lexer) Backup() {
	if l.Cursor.Pos <= l.Cursor.Start {
		return
	}
	l.Cursor.Pos -= l.Cursor.Width
	l.Location.Column--
	if res, _ := utf8.DecodeRuneInString(l.Input[l.Pos : l.Pos+l.Width]); res == Newline {
		l.Location.Line--
		l.Location.Column = l.Cursor.LastLineLength
	}
}

func (l *Lexer) Next() (res rune) {
	if l.Cursor.Pos >= len(l.Input) {
		l.Cursor.Width = 0
		return EOF
	}
	res, l.Cursor.Width = utf8.DecodeRuneInString(l.Remainder())
	l.Cursor.Pos += l.Cursor.Width
	l.Location.Column++
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

func (l *Lexer) AcceptOne(valid string) bool {
	if strings.ContainsRune(valid, l.Next()) {
		return true
	}
	l.Backup()
	return false
}

func (l *Lexer) AcceptMany(valid string) (len int) {
	for strings.ContainsRune(valid, l.Next()) {
		len++
	}
	l.Backup()
	return len
}

func (l *Lexer) Error() LexFn {
	l.Tokens <- &token.Token{
		Kind:     token.KindError,
		Value:    l.Current(),
		Location: l.Location,
	}
	return nil
}

func (l *Lexer) NextToken() *token.Token {
	t, ok := <-l.Tokens
	if ok {
		return t
	}
	return nil
}

func (l *Lexer) Run() {
	defer func() {
		l.Shutdown()
	}()
	for ; l.State != nil; l.State = l.State(l) {
	}
}

func (l *Lexer) Shutdown() {
	close(l.Tokens)
}

func New(src string, init LexFn) *Lexer {
	return &Lexer{
		Input:    src,
		Tokens:   make(chan *token.Token, len(src)/2),
		State:    init,
		Cursor:   Cursor{},
		Location: token.Location{Line: 1},
	}
}

func NewFromFile(path string, init LexFn) *Lexer {
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &Lexer{
		Input:    string(src),
		Tokens:   make(chan *token.Token, len(src)/2),
		State:    init,
		Cursor:   Cursor{},
		Location: token.Location{Line: 1},
	}
}
