package images

import (
	"embed"
)

const (
	Alpine = "alpine.tar"
)

//go:embed alpine.tar
var Images embed.FS
