package token

import (
	"fmt"
	"strings"
)

type Kind uint

const (
	KindKeyword Kind = iota
	KindSymbol
	KindIdentifier
	KindString
	KindNumeric
	KindError
	KindEof
)

type Location struct {
	Line   int
	Column int
}

func (l Location) String() string {
	return fmt.Sprintf("%d:%d", l.Line, l.Column)
}

type Token struct {
	Value string
	Kind  Kind
	Location
}

func (t *Token) Equals(tk *Token) bool {
	if tk == nil || t.Kind != tk.Kind {
		return false
	}
	if t.Kind == KindKeyword && strings.ToLower(t.Value) == strings.ToLower(tk.Value) {
		return true
	}
	return t.Value == tk.Value
}

func (t *Token) NotEquals(tk *Token) bool {
	if tk == nil || t.Kind != tk.Kind {
		return true
	}
	if t.Kind == KindKeyword && strings.ToLower(t.Value) == strings.ToLower(tk.Value) {
		return false
	}
	return t.Value != tk.Value
}

func IsKind(tk *Token, kinds ...Kind) (ok bool) {
	if tk == nil {
		return false
	}
	for _, k := range kinds {
		ok = ok || k == tk.Kind
	}
	return
}

func NotKind(tk *Token, kinds ...Kind) (ok bool) {
	if tk == nil {
		return true
	}
	ok = true
	for _, k := range kinds {
		ok = ok && k != tk.Kind
	}
	return
}

func IsEnd(tk *Token) bool {
	return tk == nil || tk.Kind == KindEof
}

func NotEnd(tk *Token) bool {
	return tk != nil && tk.Kind != KindEof
}

func (t Token) String() string {
	switch t.Kind {
	case KindEof:
		return "EOF"
	case KindError:
		return t.Value
	default:
		return fmt.Sprintf("'%s'", t.Value)
	}
}

func Error(t *Token, msg string, expect ...interface{}) error {
	if t == nil {
		return fmt.Errorf("error: missing token, expected %s", expect)
	}
	return fmt.Errorf("%s: error: %s, expected %s, got '%s'", t.Location, msg, expect, t.Value)
}

func NewSymbol(val symbol) *Token {
	return &Token{
		Kind:  KindSymbol,
		Value: string(val),
	}
}

func NewKeyword(val Keyword) *Token {
	return &Token{
		Kind:  KindKeyword,
		Value: string(val),
	}
}

type Stream struct {
	stream []*Token
	cur    int
}

func (ts *Stream) Add(token *Token) {
	ts.stream = append(ts.stream, token)
}

func (ts *Stream) Next() (t *Token) {
	if ts.cur+1 >= len(ts.stream) {
		return nil
	}
	ts.cur++
	t = ts.stream[ts.cur]
	return t
}

func (ts *Stream) Peek() (t *Token) {
	if ts.cur >= len(ts.stream) {
		return nil
	}
	ts.cur++
	t = ts.stream[ts.cur]
	ts.cur--
	return t
}

func NewStream() *Stream {
	return &Stream{
		stream: make([]*Token, 0),
		cur:    -1,
	}
}
