package token

type symbol string

const (
	SEMICOLON symbol = ";"
	ASTERISK  symbol = "*"
	COMMA     symbol = ","
	LPAREN    symbol = "("
	RPAREN    symbol = ")"
	Symbols   string = ";*,()"
)

func (s symbol) String() string {
	return string(s)
}
