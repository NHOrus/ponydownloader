// bitfmt.go
package derpiapi

import (
	"fmt"
)

const (
	KiB float64 = 1024
	MiB         = 1024 * 1024
	GiB         = 1024 * 1024 * 1024
	TiB         = 1024 * 1024 * 1024 * 1024
	PiB         = 1024 * 1024 * 1024 * 1024 * 1024
)

func fmtbytes(b float64) string {
	if b > 10*PiB {
		return fmt.Sprintf("%.2f PiB", b/PiB)
	}
	if b > 10*TiB {
		return fmt.Sprintf("%.2f TiB", b/TiB)
	}
	if b > 10*GiB {
		return fmt.Sprintf("%.2f GiB", b/GiB)
	}
	if b > 10*MiB {
		return fmt.Sprintf("%.2f MiB", b/MiB)
	}
	if b > 10*KiB {
		return fmt.Sprintf("%.2f KiB", b/KiB)
	}
	return fmt.Sprintf("%.0f B", b)
}
