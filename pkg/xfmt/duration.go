package xfmt

import (
	"fmt"
	"time"
)

func Duration(d time.Duration) string {
	td := int(d.Seconds())
	return fmt.Sprintf("%02d:%02d", td/60, td%60)
}
