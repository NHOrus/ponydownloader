package main

import (
	"fmt"
)

//Binary byte sizes of common values
const (
	_           = iota
	KiB float64 = 1 << (10 * iota)
	MiB
	GiB
	TiB
	PiB
)

func fmtbytes(b float64) string {
	switch {
	case b < 0:
		panic("Natural number is less than zero. Stuff is wrong")
	case b > PiB:
		return fmt.Sprintf("way too many B")
	case b > TiB:
		return fmt.Sprintf("%.2f TiB", b/TiB)
	case b > GiB:
		return fmt.Sprintf("%.2f GiB", b/GiB)
	case b > MiB:
		return fmt.Sprintf("%.2f MiB", b/MiB)
	case b > KiB:
		return fmt.Sprintf("%.2f KiB", b/KiB)
	default:
		return fmt.Sprintf("%.0f B", b)
	}
}
