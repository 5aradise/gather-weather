package root

import "embed"

//go:embed sql/schema
var Migrations embed.FS

//go:embed public
var Public embed.FS
