package fail2ban

type Add struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	MaxRetry    int    `json:"maxretry"`
	FindTime    int    `json:"findtime"`
	BanTime     int    `json:"bantime"`
	WebsiteName string `json:"website_name"`
	WebsiteMode string `json:"website_mode"`
	WebsitePath string `json:"website_path"`
}

type Delete struct {
	Name string `json:"name"`
}

type BanList struct {
	Name string `json:"name"`
}

type Unban struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

type SetWhiteList struct {
	IP string `json:"ip"`
}
