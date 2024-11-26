package types

type PostgresUser struct {
	Role       string   `json:"role"`
	Attributes []string `json:"attributes"`
}

type PostgresDatabase struct {
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Encoding string `json:"encoding"`
	Comment  string `json:"comment"`
}
