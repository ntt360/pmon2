package crypto

import (
	"hash/crc32"
	"strconv"
)

func Crc32Hash(val string) string {
	return strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(val))), 16)
}
