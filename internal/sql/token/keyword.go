package token

type keyword string

const (
	SELECT keyword = "select"
	FROM   keyword = "from"
	WHERE  keyword = "where"
	AS     keyword = "as"
	TABLE  keyword = "table"
	CREATE keyword = "create"
	INSERT keyword = "insert"
	INTO   keyword = "into"
	VALUES keyword = "values"
	INT    keyword = "int"
	TEXT   keyword = "text"
	OR     keyword = "or"
	AND    keyword = "and"
	TRUE   keyword = "true"
	FALSE  keyword = "false"
)

func (s keyword) String() string {
	return string(s)
}

var Keywords = []string{
	string(AS),
	string(FROM),
	string(CREATE),
	string(INSERT),
	string(INT),
	string(SELECT),
	string(INTO),
	string(TABLE),
	string(TEXT),
	string(VALUES),
}
