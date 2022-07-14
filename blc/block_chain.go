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

// 2. 增加区块到区块链
func (bc *BlockChain) AddBlockToBlockChain(data string, height int64, prevHash []byte) {
	// 创建新区块
	newBlock := NewBlock(data, height, prevHash)
	bc.Blocks = append(bc.Blocks, newBlock)
}
