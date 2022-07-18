package blc

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutputs struct {
	UTXOs []*UTXO
}

// 将TXOutputs序列化成字节数组
func (ts *TXOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(ts)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化
func DeserializeTXOutputs(txOutputsBytes []byte) *TXOutputs {
	var txOutputs TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}
	return &txOutputs
}
