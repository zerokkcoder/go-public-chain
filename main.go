package main

import (
	"fmt"
	"go-public-chain/blc"
)

func main() {
	// block := blc.NewBlock("Genenis Block", 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	genesisBlock := blc.CreateGenesisBlock("Genenis Block")
	fmt.Println(genesisBlock)
}
