package crypto

import (
	"crypto/md5"
	"fmt"
)

func Md5Checksum(input string) string {
	hash := md5.Sum([]byte(input))
	return fmt.Sprintf("%x", hash[:])
}
