package embed

import "embed"

//go:embed public/*
var PublicFS embed.FS
