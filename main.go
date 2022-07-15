package main

import "go-public-chain/blc"

func main() {
	// 创世区块
	blockChain := blc.CreateBlockChainWithGenesisBlock()
	defer blockChain.DB.Close()

	// 新区块
	blockChain.AddBlockToBlockChain("send 100rmb to zhansan")
	blockChain.AddBlockToBlockChain("send 400rmb to lisi")
	blockChain.AddBlockToBlockChain("send 300rmb to wangwu")
}
