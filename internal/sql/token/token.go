package token

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

type Token struct {
	Kind  Kind
	Value string
	Location
}

type Location struct {
	Line   int
	Column int
}

func (t *Token) equals(tk *Token) bool {
	return t.Value == tk.Value && t.Kind == tk.Kind
}
