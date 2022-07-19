package blc

import "fmt"

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := NewWallets(nodeID)
	address := wallets.CreateNewWallet(nodeID)
	fmt.Printf("新钱包地址: %s,钱包个数: %d\n", address, len(wallets.Wallets))
}
