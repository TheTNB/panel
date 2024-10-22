package embed

import "embed"

//go:embed all:frontend/*
var PublicFS embed.FS

//go:embed all:website/*
var WebsiteFS embed.FS
