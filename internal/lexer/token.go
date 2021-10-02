package lexer

type tokenKind uint

const (
	kindKeyword tokenKind = iota
	kindSymbol
	kindIdentifier
	kindString
	kindNumeric
)

type token struct {
	kind  tokenKind
	value string
	loc   location
}

func (t *token) equals(other *token) bool {
	return t.value == other.value && t.kind == other.kind
}
