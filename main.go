package main

import (
	"fmt"
	"public-blockchain/blc"
)

func main() {
	block := blc.NewBlock("Genenis Block", 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	fmt.Println(block)
}
