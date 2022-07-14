package blc

import "time"

type Block struct {
	Height        int64  // 1. 区块高度
	PrevBlockHash []byte // 2. 上一个区块 HASH
	Data          []byte // 3. 交易数据
	Timestamp     int64  // 4. 时间戳
	Hash          []byte // 5. 当前区块的 HASH
}

// 1. 创建新的区块
func NewBlock(data string, height int64, prevBlcokHash []byte) *Block {
	// 创建区块
	block := &Block{
		Height:        height,
		PrevBlockHash: prevBlcokHash,
		Data:          []byte(data),
		Timestamp:     time.Now().Unix(),
		Hash:          nil,
	}
	

	return block
}
