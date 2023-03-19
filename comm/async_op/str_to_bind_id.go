package async_op

import (
	"hash/crc32"
)

func StrToBindId(strVal string) int {

	i := int(crc32.ChecksumIEEE([]byte(strVal)))
	if i < 0 {
		i = -1
	}
	return i

}
