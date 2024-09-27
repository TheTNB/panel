package php

type UpdateConfig struct {
	Config string `form:"config" json:"config"`
}

type ExtensionSlug struct {
	Slug string `form:"slug" json:"slug"`
}
