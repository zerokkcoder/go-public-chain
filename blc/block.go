package blc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Height        int64          // 1. 区块高度
	PrevBlockHash []byte         // 2. 上一个区块 HASH
	Txs           []*Transaction // 3. 交易数据
	Timestamp     int64          // 4. 时间戳
	Hash          []byte         // 5. 当前区块的 HASH
	Nonce         int64          // 6. Nonce值 用于工作量证明
}

// 1. 创建新的区块
func NewBlock(txs []*Transaction, height int64, prevBlcokHash []byte) *Block {
	// 创建区块
	block := &Block{
		Height:        height,
		PrevBlockHash: prevBlcokHash,
		Txs:           txs,
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
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

// 将区块序列化成字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// 将 Txs 转化成 字节数组
func (b *Block) HashTransactions() []byte {
	// 原始方式:
	// var txHashes [][]byte
	// var txHash [32]byte
	// for _, tx := range b.Txs {
	// 	txHashes = append(txHashes, tx.TxHash)
	// }
	// txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	// return txHash[:]

	// 默克尔树（MerkleTree）方式:
	var transactions [][]byte
	for _, tx := range b.Txs {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}
