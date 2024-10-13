package php

type UpdateConfig struct {
	Config string `form:"config" json:"config" validate:"required"`
}

type ExtensionSlug struct {
	Slug string `form:"slug" json:"slug" validate:"required"`
}
