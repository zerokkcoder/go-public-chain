package blc

import (
	"log"

	"github.com/boltdb/bolt"
)

// 迭代器结构体
type BlockChainIterator struct {
	CurrentHash []byte   // 当前hash
	DB          *bolt.DB // 数据库
}

func (bci *BlockChainIterator) Next() *Block {
	var block *Block
	err := bci.DB.View(func(tx *bolt.Tx) error {
		// 获取表
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			currentBlockBytes := b.Get(bci.CurrentHash)
			block = DeserializeBlock(currentBlockBytes)
			bci.CurrentHash = block.PrevBlockHash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return block
}
