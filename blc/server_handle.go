package blc

import (
	"bytes"
	"encoding/gob"
	"log"
)

func handleVersion(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload Version

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	//Version
	//1. Version
	//2. BestHeight
	//3. 节点地址

	bestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if bestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	} else if bestHeight < foreignerBestHeight {
		// 去向主节点要信息
		//sendGetBlocks(payload.AddrFrom)
	}
}

func handleAddr(request []byte, bc *BlockChain) {

}
func handleBlock(request []byte, bc *BlockChain) {

}
func handleGetblocks(request []byte, bc *BlockChain) {

}
func handleGetData(request []byte, bc *BlockChain) {

}
func handleInv(request []byte, bc *BlockChain) {

}

func handleTx(request []byte, bc *BlockChain) {

}
