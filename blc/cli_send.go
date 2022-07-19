package blc

import "fmt"

// 转账
func (cli *CLI) send(from []string, to []string, amount []string, nodeID string, mineNow bool) {

	blockChain := BlockChainObject(nodeID)
	defer blockChain.DB.Close()

	if mineNow {
		blockChain.MineNewBlock(from, to, amount, nodeID)
		utxoSet := &UTXOSet{blockChain}
		// 转账成功以后，更新未花费表
		utxoSet.Update()
	} else {
		// 把交易发送到矿工节点去进行验证
		fmt.Println("由矿工节点处理......")
	}

}
