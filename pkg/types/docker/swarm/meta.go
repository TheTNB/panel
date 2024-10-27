package swarm

import (
	"strconv"
	"time"
)

// Version represents the internal object version.
type Version struct {
	Index uint64 `json:",omitempty"`
}

// String implements fmt.Stringer interface.
func (v Version) String() string {
	return strconv.FormatUint(v.Index, 10)
}

// Meta is a base object inherited by most of the other once.
type Meta struct {
	Version   Version   `json:",omitempty"`
	CreatedAt time.Time `json:",omitempty"`
	UpdatedAt time.Time `json:",omitempty"`
}
