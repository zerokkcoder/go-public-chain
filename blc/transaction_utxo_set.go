package blc

import (
	"log"

	"github.com/boltdb/bolt"
)

const utxoTableName = "utxos.db"

type UTXOSet struct {
	BlockChain *BlockChain
}

// ResetUTXOSet 重置数据库表
// 遍历整个数据库，读取所有的未花费的UTXO，
// 然后将所有的 UTXO 存储到数据库
func (us *UTXOSet) ResetUTXOSet() {
	err := us.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			tx.DeleteBucket([]byte(utxoTableName))
			b, _ := tx.CreateBucket([]byte(utxoTableName))
			if b != nil {
				txOutputsMap := us.BlockChain.FindUTXOMap()
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
