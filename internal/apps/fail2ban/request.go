package fail2ban

type Add struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	MaxRetry    int    `json:"maxretry" validate:"required"`
	FindTime    int    `json:"findtime" validate:"required"`
	BanTime     int    `json:"bantime" validate:"required"`
	WebsiteName string `json:"website_name"`
	WebsiteMode string `json:"website_mode"`
	WebsitePath string `json:"website_path"`
}

type Delete struct {
	Name string `json:"name" validate:"required"`
}

type BanList struct {
	Name string `json:"name" validate:"required"`
}

type Unban struct {
	Name string `json:"name" validate:"required"`
	IP   string `json:"ip" validate:"required"`
}

type SetWhiteList struct {
	IP string `json:"ip" validate:"required"`
}
