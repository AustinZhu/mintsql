package token

type keyword string

const (
	Select keyword = "select"
	From   keyword = "from"
	As     keyword = "as"
	Table  keyword = "table"
	Create keyword = "create"
	Insert keyword = "insert"
	Into   keyword = "into"
	Values keyword = "values"
	Int    keyword = "int"
	Text   keyword = "text"
)
