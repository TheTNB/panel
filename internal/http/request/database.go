package request

type DatabaseCreate struct {
	ServerID   uint   `form:"server_id" json:"server_id" validate:"required,exists=database_servers id"`
	Name       string `form:"name" json:"name" validate:"required"`
	CreateUser bool   `form:"create_user" json:"create_user"`
	Username   string `form:"username" json:"username" validate:"required_if=CreateUser true"`
	Password   string `form:"password" json:"password" validate:"required_if=CreateUser true"`
	Host       string `form:"host" json:"host"`
	Remark     string `form:"remark" json:"remark"`
}

type DatabaseDelete struct {
	ServerID uint   `form:"server_id" json:"server_id" validate:"required,exists=database_servers id"`
	Name     string `form:"name" json:"name" validate:"required"`
}
