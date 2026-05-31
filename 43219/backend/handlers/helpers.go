package handlers

import (
	"strconv"
	"time"
)

func timeNow() time.Time { return time.Now() }

func uintStr(v uint) string { return strconv.FormatUint(uint64(v), 10) }
