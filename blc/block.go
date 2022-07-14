package blc

import (
	"fmt"
	"time"
)

type Block struct {
	Height        int64  // 1. 区块高度
	PrevBlockHash []byte // 2. 上一个区块 HASH
	Data          []byte // 3. 交易数据
	Timestamp     int64  // 4. 时间戳
	Hash          []byte // 5. 当前区块的 HASH
	Nonce         int64  // 6. Nonce值 用于工作量证明
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
		Nonce:         0,
	}
	// 调用工作量证明的方法并且返回有效的 Hash 和 Nonce 值
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	fmt.Println()

	return block
}

// 2. 生成创世区块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
