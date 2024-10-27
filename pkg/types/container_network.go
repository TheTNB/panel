package types

import "time"

type ContainerNetwork struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Driver     string               `json:"driver"`
	IPv6       bool                 `json:"ipv6"`
	Internal   bool                 `json:"internal"`
	Attachable bool                 `json:"attachable"`
	Ingress    bool                 `json:"ingress"`
	Scope      string               `json:"scope"`
	CreatedAt  time.Time            `json:"created_at"`
	IPAM       ContainerNetworkIPAM `json:"ipam"`
	Options    []KV                 `json:"options"`
	Labels     []KV                 `json:"labels"`
}

// ContainerNetworkIPAM represents IP Address Management
type ContainerNetworkIPAM struct {
	Driver  string
	Options map[string]string // Per network IPAM driver options
	Config  []ContainerNetworkIPAMConfig
}

// ContainerNetworkIPAMConfig represents IPAM configurations
type ContainerNetworkIPAMConfig struct {
	Subnet     string            `json:"subnet"`
	IPRange    string            `json:"ip_range"`
	Gateway    string            `json:"gateway"`
	AuxAddress map[string]string `json:"AuxiliaryAddresses,omitempty"`
}

type ContainerNetworkInspect struct {
	Name       string    `json:"Name"`
	Id         string    `json:"Id"`
	Created    time.Time `json:"Created"`
	Scope      string    `json:"Scope"`
	Driver     string    `json:"Driver"`
	EnableIPv6 bool      `json:"EnableIPv6"`
	IPAM       ContainerNetworkIPAM
	Internal   bool              `json:"Internal"`
	Attachable bool              `json:"Attachable"`
	Ingress    bool              `json:"Ingress"`
	ConfigOnly bool              `json:"ConfigOnly"`
	Options    map[string]string `json:"Options"`
	Labels     map[string]string `json:"Labels"`
}
