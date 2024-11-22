package types

type MySQLUser struct {
	User   string   `json:"user"`
	Host   string   `json:"host"`
	Grants []string `json:"grants"`
}
