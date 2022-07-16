package blc

import "fmt"

// 获取余额
func (cli *CLI) getBalance(address string) {
	blockChain := BlockChainObject()
	defer blockChain.DB.Close()
	balance := blockChain.GetBalance(address)
	fmt.Printf("%s一共有%d个Token\n", address, balance)
}
