package blc

// 创建创世区块
func (cli *CLI) createGenesisBlockChain(address string, nodeID string) {
	blockChain := CreateBlockChainWithGenesisBlock(address, nodeID)
	defer blockChain.DB.Close()

	utxoSet := &UTXOSet{blockChain}
	utxoSet.ResetUTXOSet()
}
