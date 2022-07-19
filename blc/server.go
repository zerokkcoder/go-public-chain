package blc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

//localhost:3000 主节点的地址
var knowNodes = []string{"localhost:3000"}
var nodeAddress string //全局变量，节点地址

// 启动服务器
func startServer(nodeID string, minerAdd string) {

	// 当前节点的IP地址
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)

	conn, err := net.Listen(PROTOCOL, nodeAddress)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	bc := BlockChainObject(nodeID)

	// 第一个终端：端口为3000,启动的就是主节点
	// 第二个终端：端口为3001，钱包节点
	// 第三个终端：端口号为3002，矿工节点
	if nodeAddress != knowNodes[0] {
		// 此节点是钱包节点或者矿工节点，需要向主节点发送请求同步数据
		sendVerson(knowNodes[0], bc)
	}

	for {

		// 接收客户端发送过来的数据
		conn, err := conn.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleConnection(conn, bc)
	}
}

func handleConnection(conn net.Conn, bc *BlockChain) {
	// 读取客户端发送过来的所有的数据
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Receive a Message:%s\n", request[:COMMANDLENGTH])

	command := bytesToCommand(request[:COMMANDLENGTH])

	switch command {
	case COMMAND_VERSION:
		handleVersion(request, bc)
	case COMMAND_ADDR:
		handleAddr(request, bc)
	case COMMAND_BLOCK:
		handleBlock(request, bc)
	case COMMAND_GETBLOCKS:
		handleGetblocks(request, bc)
	case COMMAND_GETDATA:
		handleGetData(request, bc)
	case COMMAND_INV:
		handleInv(request, bc)
	case COMMAND_TX:
		handleTx(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	defer conn.Close()
}

func handleVersion(request []byte, bc *BlockChain) {

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

func sendVerson(toAddress string, bc *BlockChain) {

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
