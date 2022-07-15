package blc

import (
	"log"

	"github.com/boltdb/bolt"
)

const dbName = "blockchain.db"  // 数据库名字
const blockTableName = "blocks" // 表名

type BlockChain struct {
	Tip []byte   // 最新区块的hash
	DB  *bolt.DB // 数据库
}

// 1. 创建带有创世区块的区块链
func CreateBlockChainWithGenesisBlock() *BlockChain { // 创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var blockHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panic(err)
			}
		}

		if b != nil {
			// 创建创世区块
			genesisBlock := CreateGenesisBlock("Genesis data......")
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 存储最新区块的hash
			err = b.Put([]byte("tip"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			blockHash = genesisBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// 返回区块链对象
	return &BlockChain{blockHash, db}
}

// 2. 增加区块到区块链
func (bc *BlockChain) AddBlockToBlockChain(data string) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 1. 获取表
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// 先获取最新区块
			blockBytes := b.Get(bc.Tip)
			lastBlock := DeserializeBlock(blockBytes)
			// 2. 创建新区块
			newBlock := NewBlock(data, lastBlock.Height+1, lastBlock.Hash)
			// 3. 将区块序列化并且存储到数据库中
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 4. 更新数据库里面“tip”的hash
			err = b.Put([]byte("tip"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			// 5. 更新 blockchain 的Tip
			bc.Tip = newBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
