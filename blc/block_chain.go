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
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
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
func (bc *BlockChain) AddBlockToBlockChain(data string, height int64, prevHash []byte) {
	// 创建或打开数据库
	// db, err := bolt.Open(dbName, 0600, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// err = db.Update(func(tx *bolt.Tx) error {
	// 	b, err = tx.CreateBucket([]byte(blockTableName))
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	if b != nil {

	// 	}

	// 	return nil
	// })

	// if err != nil {
	// 	log.Panic(err)
	// }

	// 创建新区块
	// newBlock := NewBlock(data, height, prevHash)
	// bc.Blocks = append(bc.Blocks, newBlock)
}
