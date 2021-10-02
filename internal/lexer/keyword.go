package lexer

type keyword string

const (
	kwSelect keyword = "select"
	kwFrom   keyword = "from"
	kwAs     keyword = "as"
	kwTable  keyword = "table"
	kwCreate keyword = "create"
	kwInsert keyword = "insert"
	kwInto   keyword = "into"
	kwValues keyword = "values"
	kwInt    keyword = "int"
	kwText   keyword = "text"
)
