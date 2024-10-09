package request

type App struct {
	Slug        string `json:"slug" form:"slug"`
	VersionSlug string `json:"version_slug" form:"version_slug"`
}

type AppSlug struct {
	Slug string `json:"slug" form:"slug"`
}

type AppUpdateShow struct {
	Slug string `json:"slug" form:"slug"`
	Show bool   `json:"show" form:"show"`
}
