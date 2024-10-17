package firewall

type Operation string

var (
	OperationAdd    Operation = "add"    // 添加
	OperationRemove Operation = "remove" // 移除
)

type Protocol string

var (
	ProtocolTCP    Protocol = "tcp"     // tcp
	ProtocolUDP    Protocol = "udp"     // udp
	ProtocolTCPUDP Protocol = "tcp/udp" // tcp/udp
)

type Strategy string

var (
	StrategyAccept Strategy = "accept" // 接受
	StrategyDrop   Strategy = "drop"   // 丢弃
	StrategyReject Strategy = "reject" // 拒绝
)

type Direction string

var (
	DirectionIn  Direction = "in"  // 传入
	DirectionOut Direction = "out" // 传出
)

type FireInfo struct {
	Family    string    `json:"family"`     // ipv4 ipv6
	Address   string    `json:"address"`    // 源地址或目标地址
	PortStart uint      `json:"port_start"` // 1-65535
	PortEnd   uint      `json:"port_end"`   // 1-65535
	Protocol  Protocol  `json:"protocol"`   // tcp udp tcp/udp
	Strategy  Strategy  `json:"strategy"`   // accept drop reject
	Direction Direction `json:"direction"`  // in out 入站或出站
}

type FireForwardInfo struct {
	Address    string   `json:"address"`    // 源地址
	Port       uint     `json:"port"`       // 1-65535
	Protocol   Protocol `json:"protocol"`   // tcp udp tcp/udp
	TargetIP   string   `json:"targetIP"`   // 目标地址
	TargetPort string   `json:"targetPort"` // 1-65535
}

type Forward struct {
	Protocol   Protocol `json:"protocol"`   // tcp udp tcp/udp
	Port       uint     `json:"port"`       // 1-65535
	TargetIP   string   `json:"targetIP"`   // 目标地址
	TargetPort uint     `json:"targetPort"` // 1-65535
}
