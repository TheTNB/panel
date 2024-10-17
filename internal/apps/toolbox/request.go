package toolbox

import "time"

type DNS struct {
	DNS1 string `form:"dns1" json:"dns1" validate:"required"`
	DNS2 string `form:"dns2" json:"dns2" validate:"required"`
}

type SWAP struct {
	Size int64 `form:"size" json:"size" validate:"gte=0"`
}

type Timezone struct {
	Timezone string `form:"timezone" json:"timezone" validate:"required"`
}

type Time struct {
	Time time.Time `form:"time" json:"time" validate:"required"`
}

type Hostname struct {
	Hostname string `form:"hostname" json:"hostname" validate:"required,regex=^[a-zA-Z0-9][a-zA-Z0-9-]{0，61}[a-zA-Z0-9]$"`
}

type Hosts struct {
	Hosts string `form:"hosts" json:"hosts"`
}

type Password struct {
	Password string `form:"password" json:"password" validate:"required,password"`
}
