package main

import "go-public-chain/blc"

func main() {
	// 创世区块
	blockChain := blc.CreateBlockChainWithGenesisBlock()
	defer blockChain.DB.Close()
}
