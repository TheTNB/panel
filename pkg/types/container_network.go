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
	Driver  string                       `json:"driver"`
	Options []KV                         `json:"options"`
	Config  []ContainerNetworkIPAMConfig `json:"config"`
}

// ContainerNetworkIPAMConfig represents IPAM configurations
type ContainerNetworkIPAMConfig struct {
	Subnet     string            `json:"subnet"`
	IPRange    string            `json:"ip_range"`
	Gateway    string            `json:"gateway"`
	AuxAddress map[string]string `json:"aux_address"`
}
