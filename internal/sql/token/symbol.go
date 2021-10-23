package token

type symbol string

const (
	SEMICOLON symbol = ";"
	ASTERISK  symbol = "*"
	COMMA     symbol = ","
	LPAREN    symbol = "("
	RPAREN    symbol = ")"
	EQ        symbol = "="
	LE        symbol = "<"
	GE        symbol = ">"
	ADD       symbol = "+"
	SUB       symbol = "-"
	MUL       symbol = "*"
	DIV       symbol = "/"
	Symbols   string = ";*,()"
)

func (s symbol) String() string {
	return string(s)
}
