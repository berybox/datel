package static

import "embed"

// FS static files for embed import
//go:embed *
var FS embed.FS
