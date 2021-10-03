package token

type symbol string

const (
	Semicolon  symbol = ";"
	Asterisk   symbol = "*"
	Comma      symbol = ","
	LeftParen  symbol = "("
	RightParen symbol = ")"
	Newline    rune   = 10
	EOF        rune   = 0
)
