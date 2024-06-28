package archive

import (
	"path/filepath"
	"strings"
)

func BaseName(chart string) *string {
	base := filepath.Base(chart)
	ext := filepath.Ext(base)
	base = strings.TrimSuffix(base, ext)
	if strings.HasSuffix(base, ".tar") {
		base = strings.TrimSuffix(base, ".tar")
	}
	return &chart
}
