package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"go-public-chain/blc"
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

// 字节数组翻转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// version 转字节数组
func CommandToBytes(command string) []byte {
	var bytes [blc.COMMANDLENGTH]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

// 字节数组转version
func BytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

// 将结构体序列化成字节数组
func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
