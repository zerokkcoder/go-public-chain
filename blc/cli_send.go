package blc

// 转账
func (cli *CLI) send(from []string, to []string, amount []string, nodeID string) {

	blockChain := BlockChainObject(nodeID)
	defer blockChain.DB.Close()
	blockChain.MineNewBlock(from, to, amount, nodeID)

	utxoSet := &UTXOSet{blockChain}
	// 转账成功以后，更新未花费表
	utxoSet.Update()
}
