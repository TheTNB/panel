package request

type App struct {
	Slug    string `json:"slug" form:"slug" validate:"required,not_exists=apps slug"`
	Channel string `json:"channel" form:"channel" validate:"required"`
}

type AppSlug struct {
	Slug string `json:"slug" form:"slug" validate:"required"`
}

type AppUpdateShow struct {
	Slug string `json:"slug" form:"slug" validate:"required,exists=apps slug"`
	Show bool   `json:"show" form:"show"`
}
