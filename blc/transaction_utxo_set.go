package blc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"

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
			// fmt.Printf("key=%s, value=%v\n", k, v)
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

// 返回要凑多少钱，对应的TXOutput的TX的Hash和index
func (us *UTXOSet) FindUnPackageSpendableUTXOs(from string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs {
		if !tx.IsCoinbaseTransaction() {
			for _, in := range tx.Vins {
				//是否能够解锁
				publicKeyHash := Base58Decode([]byte(from))
				ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]
				if in.UnLockRipemd160Hash(ripemd160Hash) {
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}
	}

	for _, tx := range txs {

		// Vouts
	work1:
		for index, out := range tx.Vouts {
			if out.UnLockScriptPubKeyWithAddress(from) {
				if len(spentTXOutputs) != 0 {
					for hash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if hash == txHashStr {
							var isSpendUTXO bool
							for _, outIndex := range indexArray {
								if index == outIndex {
									isSpendUTXO = true
									continue work1
								}
							}
							if !isSpendUTXO {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						}
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				}
			}
		}
	}

	return unUTXOs
}

func (us *UTXOSet) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {

	unPackageUTXOs := us.FindUnPackageSpendableUTXOs(from, txs)
	spentableUTXO := make(map[string][]int)

	var money int64 = 0

	for _, UTXO := range unPackageUTXOs {
		money += UTXO.Output.Value
		txHash := hex.EncodeToString(UTXO.TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash], UTXO.Index)
		if money >= amount {
			return money, spentableUTXO
		}
	}

	// 钱还不够
	us.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			c := b.Cursor()
		UTXOLOOP:
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTXOutputs(v)
				for _, utxo := range txOutputs.UTXOs {
					if utxo.Output.UnLockScriptPubKeyWithAddress(from) {
						money += utxo.Output.Value
						txHash := hex.EncodeToString(utxo.TxHash)
						spentableUTXO[txHash] = append(spentableUTXO[txHash], utxo.Index)
						if money >= amount {
							break UTXOLOOP
						}
					}
				}
			}
		}

		return nil
	})

	if money < int64(amount) {
		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}

	return money, spentableUTXO
}

// 更新
func (us *UTXOSet) Update() {
	// 最新的Block
	block := us.BlockChain.Iterator().Next()
	// utxos表，更新
	ins := []*TXInput{}
	outsMap := make(map[string]*TXOutputs)
	// 找到需要删除的数据
	for _, tx := range block.Txs {
		for _, in := range tx.Vins {
			ins = append(ins, in)
		}
	}
	for _, tx := range block.Txs {
		utxos := []*UTXO{}

		for index, out := range tx.Vouts {
			isSpent := false
			for _, in := range ins {
				if in.Vout == index && bytes.Equal(tx.TxHash, in.TxHash) && bytes.Equal(out.Ripemd160Hash, Ripemd160Hash(in.PublicKey)) {
					isSpent = true
					continue
				}

			}
			if !isSpent {
				utxo := &UTXO{tx.TxHash, index, out}
				utxos = append(utxos, utxo)
			}
		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.TxHash)
			outsMap[txHash] = &TXOutputs{utxos}
		}
	}
	err := us.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {
			// 删除
			for _, in := range ins {
				txOutputsBytes := b.Get(in.TxHash)
				if len(txOutputsBytes) == 0 {
					continue
				}
				txOutputs := DeserializeTXOutputs(txOutputsBytes)
				UTXOs := []*UTXO{}
				// 判断是否需要删除
				isNeedDelete := false
				for _, utxo := range txOutputs.UTXOs {
					if in.Vout == utxo.Index && bytes.Equal(utxo.Output.Ripemd160Hash, Ripemd160Hash(in.PublicKey)) {
						isNeedDelete = true
					} else {
						UTXOs = append(UTXOs, utxo)
					}
				}
				if isNeedDelete {
					b.Delete(in.TxHash)
					if len(UTXOs) > 0 {
						preTXOutputs := outsMap[hex.EncodeToString(in.TxHash)]
						preTXOutputs.UTXOs = append(preTXOutputs.UTXOs, UTXOs...)
						outsMap[hex.EncodeToString(in.TxHash)] = preTXOutputs
					}
				}
			}

			// 新增
			for keyHash, outPuts := range outsMap {
				keyHashBytes, _ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes, outPuts.Serialize())
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
