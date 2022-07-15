package main

import (
	"go-public-chain/blc"
)

func main() {

	// 创世区块
	blockchain := blc.CreateBlockChainWithGenesisBlock()
	defer blockchain.DB.Close()

	cli := blc.CLI{BlockChain: blockchain}
	cli.Run()

}
