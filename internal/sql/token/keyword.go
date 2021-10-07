package token

type Keyword string

const (
	Select Keyword = "select"
	From   Keyword = "from"
	As     Keyword = "as"
	Table  Keyword = "table"
	Create Keyword = "create"
	Insert Keyword = "insert"
	Into   Keyword = "into"
	Values Keyword = "values"
	Int    Keyword = "int"
	Text   Keyword = "text"
)
