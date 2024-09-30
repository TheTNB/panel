package request

type PanelSetting struct {
	Name        string `json:"name"`
	Locale      string `json:"locale"`
	Entrance    string `json:"entrance"`
	WebsitePath string `json:"website_path"`
	BackupPath  string `json:"backup_path"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Port        int    `json:"port"`
	HTTPS       bool   `json:"https"`
	Cert        string `json:"cert"`
	Key         string `json:"key"`
}
