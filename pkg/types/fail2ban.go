package types

type Fail2banJail struct {
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
	LogPath  string `json:"log_path"`
	MaxRetry int    `json:"max_retry"`
	FindTime int    `json:"find_time"`
	BanTime  int    `json:"ban_time"`
}
