package request

type FirewallStatus struct {
	Status bool `json:"status" form:"status"`
}

type FirewallRule struct {
	Port     uint   `json:"port" validate:"required"`
	Protocol string `json:"protocol" validate:"required"`
}
