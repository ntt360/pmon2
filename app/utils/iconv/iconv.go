package iconv

import "strconv"

func MustInt(val string) int {
	valInt, _ := strconv.Atoi(val)
	return valInt
}
