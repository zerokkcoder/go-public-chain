package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

// 标准的JSON字符串转数组
func JSONToArray(jsonString string) []string {
	// json 转 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}
