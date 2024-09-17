package embed

import "embed"

//go:embed frontend/*
var PublicFS embed.FS

//go:embed website/*
var WebsiteFS embed.FS
