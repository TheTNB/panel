package docker

// Port An open port on a container
type Port struct {

	// Host IP address that the container's port is mapped to
	IP string `json:"IP,omitempty"`

	// Port on the container
	// Required: true
	PrivatePort uint16 `json:"PrivatePort"`

	// Port exposed on the host
	PublicPort uint16 `json:"PublicPort,omitempty"`

	// type
	// Required: true
	Type string `json:"Type"`
}
