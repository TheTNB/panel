package request

type AppSlug struct {
	Slug string `json:"slug" form:"slug"`
}

type AppUpdateShow struct {
	Slug string `json:"slug" form:"slug"`
	Show bool   `json:"show" form:"show"`
}
