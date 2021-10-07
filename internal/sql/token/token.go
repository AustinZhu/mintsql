package token

import "fmt"

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

func (l *Location) String() string {
	return fmt.Sprintf("%d:%d", l.Line, l.Column)
}

type Token struct {
	Kind  Kind
	Value string
	Location
}

func (t *Token) Equals(tk *Token) bool {
	return t.Value == tk.Value && t.Kind == tk.Kind
}

func (t *Token) String() string {
	switch t.Kind {
	case KindEof:
		return "EOF"
	case KindError:
		return t.Value
	default:
		return fmt.Sprintf("'%s'", t.Value)
	}
}
