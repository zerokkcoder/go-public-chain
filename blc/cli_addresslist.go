package blc

import "fmt"

// 输出所有钱包地址
func (cli *CLI) addressList() {
	fmt.Println("输出所有钱包地址:")
	wallets, _ := NewWallets()
	for address := range wallets.Wallets {
		fmt.Println(address)
	}
}
