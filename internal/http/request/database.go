package request

type DatabaseCreate struct {
	ServerID   uint   `form:"server_id" json:"server_id" validate:"required|exists:database_servers,id"`
	Name       string `form:"name" json:"name" validate:"required"`
	CreateUser bool   `form:"create_user" json:"create_user"`
	Username   string `form:"username" json:"username" validate:"requiredIf:CreateUser,true"`
	Password   string `form:"password" json:"password" validate:"requiredIf:CreateUser,true"`
	Host       string `form:"host" json:"host"`
	Comment    string `form:"comment" json:"comment"`
}

type DatabaseDelete struct {
	ServerID uint   `form:"server_id" json:"server_id" validate:"required|exists:database_servers,id"`
	Name     string `form:"name" json:"name" validate:"required"`
}

type DatabaseComment struct {
	ServerID uint   `form:"server_id" json:"server_id" validate:"required|exists:database_servers,id"`
	Name     string `form:"name" json:"name" validate:"required"`
	Comment  string `form:"comment" json:"comment"`
}
