package supervisor

type Process struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Pid    string `json:"pid"`
	Uptime string `json:"uptime"`
}
