package request

type DashboardCurrent struct {
	Nets  []string `json:"nets" form:"nets"`
	Disks []string `json:"disks" form:"disks"`
}
