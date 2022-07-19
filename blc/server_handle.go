package blc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
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

func sendVersion(toAddress string, bc *BlockChain) {

	bestHeight := bc.GetBestHeight()
	payload := gobEncode(Version{NODE_VERSION, bestHeight, nodeAddress})

	request := append(commandToBytes(COMMAND_VERSION), payload...)

	sendData(toAddress, request)
}

func sendData(to string, data []byte) {

	fmt.Println("客户端向服务器发送数据......")
	conn, err := net.Dial("tcp", to)
	if err != nil {
		panic("error")
	}
	defer conn.Close()

	// 附带要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}
