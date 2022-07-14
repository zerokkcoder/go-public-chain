package main

import (
	"fmt"
	"go-public-chain/blc"
)

func main() {
	genesisBlockChain := blc.CreateBlockChainWithGenesisBlock()
	fmt.Println(genesisBlockChain)
	fmt.Println(genesisBlockChain.Blocks)
	fmt.Println(genesisBlockChain.Blocks[0])
}
