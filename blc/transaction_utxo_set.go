package blc

import (
	"encoding/hex"
	"fmt"
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
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err != nil {
				log.Panic(err)
			}
		}

		b, _ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {
			// 返回类型: map[string]*TXOutputs
			txOutputsMap := us.BlockChain.FindUTXOMap()

			for keyHash, outs := range txOutputsMap {
				txHash, _ := hex.DecodeString(keyHash)
				b.Put(txHash, outs.Serialize())
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (us *UTXOSet) findUTXOForAddress(address string) []*UTXO {
	var utxos []*UTXO
	us.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%v\n", k, v)
			txOutputs := DeserializeTXOutputs(v)
			for _, utxo := range txOutputs.UTXOs {
				if utxo.Output.UnLockScriptPubKeyWithAddress(address) {
					utxos = append(utxos, utxo)
				}
			}
		}

		return nil
	})
	return utxos
}

// 查询余额
func (us *UTXOSet) GetBalance(address string) int64 {
	UTXOs := us.findUTXOForAddress(address)
	var amount int64
	for _, utxo := range UTXOs {
		amount += utxo.Output.Value
	}

	return amount
}
