package request

type PanelSetting struct {
	Name        string `json:"name" validate:"required"`
	Locale      string `json:"locale" validate:"required"`
	Entrance    string `json:"entrance" validate:"required"`
	WebsitePath string `json:"website_path" validate:"required"`
	BackupPath  string `json:"backup_path" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password"`
	Email       string `json:"email" validate:"required"`
	Port        int    `json:"port" validate:"required,number,gte=1,lte=65535"`
	HTTPS       bool   `json:"https"`
	Cert        string `json:"cert" validate:"required"`
	Key         string `json:"key" validate:"required"`
}
