package types

type EnvironmentPHPModule struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Installed   bool   `json:"installed"`
}

// EnvironmentPHPProcess PHP-FPM 工作进程信息
type EnvironmentPHPProcess struct {
	PID             int64   `json:"pid"`
	State           string  `json:"state"`
	StartSince      int64   `json:"start_since"`
	Requests        int64   `json:"requests"`
	RequestDuration int64   `json:"request_duration"` // 微秒
	Method          string  `json:"method"`
	URI             string  `json:"uri"`
	Script          string  `json:"script"`
	LastRequestCPU  float64 `json:"last_request_cpu"`
	LastRequestMem  int64   `json:"last_request_memory"`
}

// EnvironmentPHPOpcache OPcache 状态
type EnvironmentPHPOpcache struct {
	Enabled       bool    `json:"enabled"`
	Message       string  `json:"message"` // 未启用时的原因
	MemoryUsed    string  `json:"memory_used"`
	MemoryFree    string  `json:"memory_free"`
	MemoryWasted  string  `json:"memory_wasted"`
	WastedPercent float64 `json:"wasted_percent"`
	HitRate       float64 `json:"hit_rate"`
	Hits          int64   `json:"hits"`
	Misses        int64   `json:"misses"`
	CachedScripts int64   `json:"cached_scripts"`
	CachedKeys    int64   `json:"cached_keys"`
	MaxCachedKeys int64   `json:"max_cached_keys"`
	OomRestarts   int64   `json:"oom_restarts"`
	JitEnabled    bool    `json:"jit_enabled"`
	JitBufferSize string  `json:"jit_buffer_size"`
	JitBufferFree string  `json:"jit_buffer_free"`
}

// EnvironmentPHPComposer Composer 状态
type EnvironmentPHPComposer struct {
	Installed bool   `json:"installed"`
	Version   string `json:"version"`
	Mirror    string `json:"mirror"` // 空为官方源
}
