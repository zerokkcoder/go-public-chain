package blc

import "fmt"

func (cli *CLI) createWallet() {
	wallets, _ := NewWallets()
	address := wallets.CreateNewWallet()
	fmt.Printf("新钱包地址: %s,钱包个数: %d\n", address, len(wallets.Wallets))
}
