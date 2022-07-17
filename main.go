package main

import (
	"fmt"
	"go-public-chain/blc"
)

func main() {
	// cli := blc.CLI{}
	// cli.Run()

	wallets := blc.NewWallets()
	fmt.Println(wallets.Wallets)
	wallets.CreateNewWallet()
	wallets.CreateNewWallet()
	fmt.Println(wallets.Wallets)
}
