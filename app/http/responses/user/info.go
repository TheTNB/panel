package responses

type Info struct {
	ID       uint     `json:"id"`
	Role     []string `json:"role"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
}
