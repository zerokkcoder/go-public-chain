package main

import (
	"fmt"
	"go-public-chain/blc"
)

func main() {
	// cli := blc.CLI{}
	// cli.Run()

	wallet := blc.NewWallet()
	address := wallet.GetAddress()

	fmt.Printf("address:%s\n", address)

}
