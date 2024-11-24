package types

type MySQLUser struct {
	User   string   `json:"user"`
	Host   string   `json:"host"`
	Grants []string `json:"grants"`
}

type MySQLDatabase struct {
	Name      string `json:"name"`
	CharSet   string `json:"char_set"`
	Collation string `json:"collation"`
}
