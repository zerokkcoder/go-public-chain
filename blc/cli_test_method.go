package blc

func (cli *CLI) TestMethod() {

	blockChain := BlockChainObject()
	defer blockChain.DB.Close()

	utxoSet := &UTXOSet{blockChain}
	utxoSet.ResetUTXOSet()
}
