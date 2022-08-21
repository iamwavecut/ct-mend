// Package resources to embed
package resources

import (
	"embed"
)

//go:embed all:*
var FS embed.FS
