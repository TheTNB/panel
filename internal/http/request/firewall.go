package request

type FirewallStatus struct {
	Status bool `json:"status" form:"status"`
}

type FirewallCreateRule struct {
	Port     uint   `json:"port"`
	Protocol string `json:"protocol"`
}
