package fail2ban

type Jail struct {
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
	MaxRetry int    `json:"max_retry"`
	FindTime int    `json:"find_time"`
	BanTime  int    `json:"ban_time"`
}
