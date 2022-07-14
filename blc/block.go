package blc

import (
	"bytes"
	"crypto/sha256"
	"go-public-chain/utils"
	"strconv"
	"time"
)

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

	// 设置 Hash
	block.SetHash()

	return block
}

// SetHash 设置 Hash
func (b *Block) SetHash() {
	// 1. Height 转 []byte
	heightBytes := utils.IntToHex(b.Height)
	// 2. 时间戳 转 []byte
	timeString := strconv.FormatInt(b.Timestamp, 2) // 转换成二进制数据
	timeBytes := []byte(timeString)                 // 转换成 []byte
	// 3. 拼接所有属性
	blockBytes := bytes.Join([][]byte{heightBytes, b.PrevBlockHash, b.Data, timeBytes, b.Hash}, []byte{})
	// 4. 生产 Hash
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}
