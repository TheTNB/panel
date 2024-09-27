package pureftpd

type User struct {
	Username string `json:"username"`
	Path     string `json:"path"`
}
