package blc

// 打印区块链
func (cli *CLI) printChain(nodeID string) {
	blockChain := BlockChainObject(nodeID)
	defer blockChain.DB.Close()
	blockChain.PrintChain()
}
