package blc

import "fmt"

// 获取余额
func (cli *CLI) getBalance(address string, nodeID string) {
	blockChain := BlockChainObject(nodeID)
	defer blockChain.DB.Close()

	utxoSet := &UTXOSet{blockChain}

	balance := utxoSet.GetBalance(address)
	fmt.Printf("%s 一共有%d个Token\n", address, balance)
}
