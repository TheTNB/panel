package request

type MonitorSetting struct {
	Enabled bool `json:"enabled"`
	Days    uint `json:"days"`
}

type MonitorList struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}
