package blc

type BlockChain struct {
	Blocks []*Block // 存储有序的区块
}

// 1. 创建带有创世区块的区块链
func CreateBlockChainWithGenesisBlock() *BlockChain {
	// 创建创世区块
	genesisBlock := CreateGenesisBlock("Genesis data......")
	// 返回区块链对象
	return &BlockChain{[]*Block{genesisBlock}}
}
