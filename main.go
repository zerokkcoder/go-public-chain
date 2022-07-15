package main

import (
	"flag"
	"fmt"
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
	flagPrintChainCmd := flag.String("printchain", "", "输出所有区块")
	flag.Parse()
	fmt.Printf("%s\n", *flagPrintChainCmd)
}
