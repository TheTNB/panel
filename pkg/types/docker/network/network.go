package network // import "github.com/docker/docker/api/types/network"

import (
	"time"
)

// Network is the body of the "get network" http response message.
type Network struct {
	Name       string                      // Name is the name of the network
	ID         string                      `json:"Id"` // ID uniquely identifies a network on a single machine
	Created    time.Time                   // Created is the time the network created
	Scope      string                      // Scope describes the level at which the network exists (e.g. `swarm` for cluster-wide or `local` for machine level)
	Driver     string                      // Driver is the Driver name used to create the network (e.g. `bridge`, `overlay`)
	EnableIPv6 bool                        // EnableIPv6 represents whether to enable IPv6
	IPAM       IPAM                        // IPAM is the network's IP Address Management
	Internal   bool                        // Internal represents if the network is used internal only
	Attachable bool                        // Attachable represents if the global scope is manually attachable by regular containers from workers in swarm mode.
	Ingress    bool                        // Ingress indicates the network is providing the routing-mesh for the swarm cluster.
	ConfigFrom ConfigReference             // ConfigFrom specifies the source which will provide the configuration for this network.
	ConfigOnly bool                        // ConfigOnly networks are place-holder networks for network configurations to be used by other networks. ConfigOnly networks cannot be used directly to run containers or services.
	Containers map[string]EndpointResource // Containers contains endpoints belonging to the network
	Options    map[string]string           // Options holds the network specific options to use for when creating the network
	Labels     map[string]string           // Labels holds metadata specific to the network being created
	Peers      []PeerInfo                  `json:",omitempty"` // List of peer nodes for an overlay network
	Services   map[string]ServiceInfo      `json:",omitempty"`
}

// PeerInfo represents one peer of an overlay network
type PeerInfo struct {
	Name string
	IP   string
}

// Task carries the information about one backend task
type Task struct {
	Name       string
	EndpointID string
	EndpointIP string
	Info       map[string]string
}

// ServiceInfo represents service parameters with the list of service's tasks
type ServiceInfo struct {
	VIP          string
	Ports        []string
	LocalLBIndex int
	Tasks        []Task
}

// EndpointResource contains network resources allocated and used for a
// container in a network.
type EndpointResource struct {
	Name        string
	EndpointID  string
	MacAddress  string
	IPv4Address string
	IPv6Address string
}

// ConfigReference specifies the source which provides a network's configuration
type ConfigReference struct {
	Network string
}
