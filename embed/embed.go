package embed

import "embed"

//go:embed frontend/*
var PublicFS embed.FS
