package plugins

type PHPExtension struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Installed   bool   `json:"installed"`
}

type LoadInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Fail2banJail struct {
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
	LogPath  string `json:"log_path"`
	MaxRetry int    `json:"max_retry"`
	FindTime int    `json:"find_time"`
	BanTime  int    `json:"ban_time"`
}

type S3fsMount struct {
	ID     int64  `json:"id"`
	Path   string `json:"path"`
	Bucket string `json:"bucket"`
	Url    string `json:"url"`
}
