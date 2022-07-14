package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

// 将 int64 转换为 字节数组
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
