package blc

func (cli *CLI) TestMethod(nodeID string) {

	blockChain := BlockChainObject(nodeID)
	defer blockChain.DB.Close()

	utxoSet := &UTXOSet{blockChain}
	utxoSet.ResetUTXOSet()
}
