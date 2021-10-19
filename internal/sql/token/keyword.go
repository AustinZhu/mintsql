package token

type keyword string

const (
	SELECT keyword = "select"
	FROM   keyword = "from"
	AS     keyword = "as"
	TABLE  keyword = "table"
	CREATE keyword = "create"
	INSERT keyword = "insert"
	INTO   keyword = "into"
	VALUES keyword = "values"
	INT    keyword = "int"
	TEXT   keyword = "text"
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
