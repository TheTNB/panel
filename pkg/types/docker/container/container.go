package container

import (
	"github.com/tnb-labs/panel/pkg/types/docker"
)

// Container contains response of Engine API:
// GET "/containers/json"
type Container struct {
	ID         string `json:"Id"`
	Names      []string
	Image      string
	ImageID    string
	Command    string
	Created    int64
	Ports      []docker.Port
	SizeRw     int64 `json:",omitempty"`
	SizeRootFs int64 `json:",omitempty"`
	Labels     map[string]string
	State      string
	Status     string
	HostConfig struct {
		NetworkMode string            `json:",omitempty"`
		Annotations map[string]string `json:",omitempty"`
	}
	Mounts []MountPoint
}

// MountPoint represents a mount point configuration inside the container.
// This is used for reporting the mountpoints in use by a container.
type MountPoint struct {
	// Type is the type of mount, see `Type<foo>` definitions in
	// github.com/docker/docker/api/types/mount.Type
	Type string `json:",omitempty"`

	// Name is the name reference to the underlying data defined by `Source`
	// e.g., the volume name.
	Name string `json:",omitempty"`

	// Source is the source location of the mount.
	//
	// For volumes, this contains the storage location of the volume (within
	// `/var/lib/docker/volumes/`). For bind-mounts, and `npipe`, this contains
	// the source (host) part of the bind-mount. For `tmpfs` mount points, this
	// field is empty.
	Source string

	// Destination is the path relative to the container root (`/`) where the
	// Source is mounted inside the container.
	Destination string

	// Driver is the volume driver used to create the volume (if it is a volume).
	Driver string `json:",omitempty"`

	// Mode is a comma separated list of options supplied by the user when
	// creating the bind/volume mount.
	//
	// The default is platform-specific (`"z"` on Linux, empty on Windows).
	Mode string

	// RW indicates whether the mount is mounted writable (read-write).
	RW bool

	// Propagation describes how mounts are propagated from the host into the
	// mount point, and vice-versa. Refer to the Linux kernel documentation
	// for details:
	// https://www.kernel.org/doc/Documentation/filesystems/sharedsubtree.txt
	//
	// This field is not used on Windows.
	Propagation string
}
