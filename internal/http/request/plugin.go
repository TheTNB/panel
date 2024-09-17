package request

type PluginSlug struct {
	Slug string `json:"slug" form:"slug"`
}

type PluginUpdateShow struct {
	Slug string `json:"slug" form:"slug"`
	Show bool   `json:"show" form:"show"`
}
