package srcfiles

import (
	"embed"
)

//go:embed **/*
var FS embed.FS
