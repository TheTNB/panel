package request

type DatabaseUserCreate struct {
	ServerID   uint     `form:"server_id" json:"server_id" validate:"required|exists:database_servers,id"`
	Username   string   `form:"username" json:"username" validate:"required"`
	Password   string   `form:"password" json:"password" validate:"required"`
	Host       string   `form:"host" json:"host"`
	Privileges []string `form:"privileges" json:"privileges"`
	Remark     string   `form:"remark" json:"remark"`
}

type DatabaseUserUpdate struct {
	ID         uint     `form:"id" json:"id" validate:"required|exists:database_users,id"`
	Password   string   `form:"password" json:"password"`
	Privileges []string `form:"privileges" json:"privileges"`
	Remark     string   `form:"remark" json:"remark"`
}

type DatabaseUserUpdateRemark struct {
	ID     uint   `form:"id" json:"id" validate:"required|exists:database_users,id"`
	Remark string `form:"remark" json:"remark"`
}
