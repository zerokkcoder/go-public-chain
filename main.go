package main

import (
	"fmt"
	"os"
)

func main() {
	// // 创世区块
	// blockChain := blc.CreateBlockChainWithGenesisBlock()
	// defer blockChain.DB.Close()

	// // 新区块
	// blockChain.AddBlockToBlockChain("send 100rmb to zhansan")
	// blockChain.AddBlockToBlockChain("send 400rmb to lisi")
	// blockChain.AddBlockToBlockChain("send 300rmb to wangwu")

	// // 遍历区块
	// blockChain.PrintChain()

	args := os.Args
	fmt.Printf("%v\n", args)
}
