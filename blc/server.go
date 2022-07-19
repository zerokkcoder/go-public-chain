package blc

import (
	"bytes"
	"fmt"
	"go-public-chain/utils"
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

		// 读取客户端发送过来的所有的数据
		request, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}

		fmt.Printf("Receive a Message:%s\n", request)
	}
}

func sendMessage(to string, from string) {

	fmt.Println("客户端向服务器发送数据......")
	conn, err := net.Dial("tcp", to)
	if err != nil {
		panic("error")
	}
	defer conn.Close()

	// 附带要发送的数据
	_, err = io.Copy(conn, bytes.NewReader([]byte(from)))
	if err != nil {
		log.Panic(err)
	}
}

func sendVerson(toAddress string, bc *BlockChain) {
	
	bestHeight := bc.GetBestHeight()
	payload := utils.GobEncode(Version{NODE_VERSION, bestHeight, nodeAddress})

	request := append(utils.CommandToBytes(VERSION), payload...)

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
