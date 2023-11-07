package responses

type Settings struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Port        string `json:"port"`
	Entrance    string `json:"entrance"`
	WebsitePath string `json:"website_path"`
	BackupPath  string `json:"backup_path"`
}
