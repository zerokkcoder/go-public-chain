package main

import (
	"fmt"
	"go-public-chain/blc"
)

func main() {
	// 创世区块
	blockChain := blc.CreateBlockChainWithGenesisBlock()

	// 新区块
	blockChain.AddBlockToBlockChain("Send 100Rmb To zhangshan", blockChain.Blocks[len(blockChain.Blocks)-1].Height+1, blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	blockChain.AddBlockToBlockChain("Send 200Rmb To hah", blockChain.Blocks[len(blockChain.Blocks)-1].Height+1, blockChain.Blocks[len(blockChain.Blocks)-1].Hash)

	fmt.Println(blockChain)
	fmt.Println(blockChain.Blocks)
}
