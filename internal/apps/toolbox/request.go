package toolbox

type DNS struct {
	DNS1 string `form:"dns1" json:"dns1"`
	DNS2 string `form:"dns2" json:"dns2"`
}

type SWAP struct {
	Size int64 `form:"size" json:"size"`
}

type Timezone struct {
	Timezone string `form:"timezone" json:"timezone"`
}

type Hosts struct {
	Hosts string `form:"hosts" json:"hosts"`
}

type Password struct {
	Password string `form:"password" json:"password"`
}
