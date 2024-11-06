package request

type DatabaseCreate struct {
	ServerID string `form:"server_id" json:"server_id" validate:"required,exists=database_servers id"`
	Name     string `form:"name" json:"name" validate:"required"`
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
	Remark   string `form:"remark" json:"remark"`
}

type DatabaseUpdate struct {
	ID       string `form:"id" json:"id" validate:"required,exists=databases id"`
	Name     string `form:"name" json:"name" validate:"required"`
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
	Remark   string `form:"remark" json:"remark"`
}
