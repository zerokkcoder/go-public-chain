package blc

import "fmt"

// 输出所有钱包地址
func (cli *CLI) addressList(nodeID string) {
	fmt.Println("输出所有钱包地址:")
	wallets, _ := NewWallets(nodeID)
	for address := range wallets.Wallets {
		fmt.Println(address)
	}
}
