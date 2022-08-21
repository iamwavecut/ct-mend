// Package resources to embed
package resources

import (
	"embed"
)

//go:embed *
var FS embed.FS
