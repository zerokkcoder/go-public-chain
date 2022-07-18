package blc

// 创建创世区块
func (cli *CLI) createGenesisBlockChain(address string) {
	blockChain := CreateBlockChainWithGenesisBlock(address)
	defer blockChain.DB.Close()

	utxoSet := &UTXOSet{blockChain}
	utxoSet.ResetUTXOSet()
}
