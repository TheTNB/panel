package request

type App struct {
	Slug    string `json:"slug" form:"slug"`
	Channel string `json:"channel" form:"channel"`
}

type AppSlug struct {
	Slug string `json:"slug" form:"slug"`
}

type AppUpdateShow struct {
	Slug string `json:"slug" form:"slug"`
	Show bool   `json:"show" form:"show"`
}
