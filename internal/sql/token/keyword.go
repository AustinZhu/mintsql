package token

type Keyword string

const (
	SELECT Keyword = "select"
	FROM   Keyword = "from"
	AS     Keyword = "as"
	TABLE  Keyword = "table"
	CREATE Keyword = "create"
	INSERT Keyword = "insert"
	INTO   Keyword = "into"
	VALUES Keyword = "values"
	INT    Keyword = "int"
	TEXT   Keyword = "text"
)

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
