package blc

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const dbName = "blockchain.db"  // 数据库名字
const blockTableName = "blocks" // 表名

type BlockChain struct {
	Tip []byte   // 最新区块的hash
	DB  *bolt.DB // 数据库
}

// 1. 创建带有创世区块的区块链
func CreateBlockChainWithGenesisBlock(address string) *BlockChain {
	// 判断数据库是否存在
	if DBExists() {
		fmt.Println("创世区块已经存在...")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块.......")

	// 创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte

	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			// 创建一个 coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(address)
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
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

			genesisHash = genesisBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &BlockChain{genesisHash, db}
}

// 查找一个地址对应的TxOutput未花费的所有 TXOutput
func (bc *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UXTO {

	var unUTXOs []*UXTO

	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs {

		// Vouts
	work1:
		for index, out := range tx.Vouts {
			if out.UnLockScriptPubKeyWithAddress(address) {
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
								utxo := &UXTO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						}
					}
				} else {
					utxo := &UXTO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				}
			}
		}
	}

	blockIterator := bc.Iterator()

	for {
		block := blockIterator.Next()

		// txHash
		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]
			// Vins
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					// 判断是否解锁
					pubKeyHash := Base58Decode([]byte(address))
					ripemd160Hash := pubKeyHash[1 : len(pubKeyHash)-4]
					if in.UnLockRipemd160Hash(ripemd160Hash) {
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}
			// Vouts
		work2:
			for index, out := range tx.Vouts {
				if out.UnLockScriptPubKeyWithAddress(address) {
					if len(spentTXOutputs) != 0 {
						var isSpendUTXO bool
						for txHash, indexArray := range spentTXOutputs {
							for _, i := range indexArray {
								if index == i && txHash == hex.EncodeToString(tx.TxHash) {
									isSpendUTXO = true
									continue work2
								}
							}
						}
						if !isSpendUTXO {
							utxo := &UXTO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					} else {
						utxo := &UXTO{tx.TxHash, index, out}
						unUTXOs = append(unUTXOs, utxo)
					}
				}
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return unUTXOs
}

// 查询余额
func (bc *BlockChain) GetBalance(address string) int64 {
	utxos := bc.UnUTXOs(address, []*Transaction{})

	var amount int64

	for _, utxo := range utxos {
		amount += utxo.Output.Value
	}

	return amount
}

// 转账时查找可用的UTXO
func (bc *BlockChain) FindSpendableUTXOs(from string, amount int, txs []*Transaction) (int64, map[string][]int) {
	// 1. 先获取所有的 UTXO
	utxos := bc.UnUTXOs(from, txs)
	var spendableUTXODic = make(map[string][]int)
	// 2. 遍历 utxos
	var value int64
	for _, utxo := range utxos {
		value += utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXODic[hash] = append(spendableUTXODic[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {
		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}
	return value, spendableUTXODic
}

// 挖掘新的区块
func (bc *BlockChain) MineNewBlock(from []string, to []string, amount []string) {
	// $ go run .\main.go send -from '[\"huanggz\"]' -to '[\"lisi\"]' -amount '[\"6\"]'
	// [huanggz]
	// [lisi]
	// [6]

	// 建立交易
	var txs []*Transaction
	for index, address := range from {
		amountInt, _ := strconv.Atoi(amount[index])
		// 1. 通过相关算法建立 Transaction 数组
		tx := NewSimpleTransaction(address, to[index], amountInt, bc, txs)
		txs = append(txs, tx)
	}

	// 2. 建立新的区块
	// 获取最新区块 height 和 Hash
	var block *Block
	err := bc.DB.View(func(tx *bolt.Tx) error {
		// 获取表
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("tip"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	block = NewBlock(txs, block.Height+1, block.Hash)

	// 3. 将新区块存储到数据库
	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			b.Put(block.Hash, block.Serialize())
			b.Put([]byte("tip"), block.Hash)
			bc.Tip = block.Hash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// 2. 增加区块到区块链
func (bc *BlockChain) AddBlockToBlockChain(txs []*Transaction) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 1. 获取表
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// 先获取最新区块
			blockBytes := b.Get(bc.Tip)
			lastBlock := DeserializeBlock(blockBytes)
			// 2. 创建新区块
			newBlock := NewBlock(txs, lastBlock.Height+1, lastBlock.Hash)
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

// 迭代器
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.Tip, bc.DB}
}

// 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 遍历区块链
func (bc *BlockChain) PrintChain() {
	// 获取迭代器
	blockChainIterator := bc.Iterator()
	for {
		block := blockChainIterator.Next()

		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 PM"))
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.Txs {
			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%v\n", in.PublicKey)
			}
			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Printf("%d\n", out.Value)
				fmt.Printf("%v\n", out.Ripemd160Hash)
			}
		}

		fmt.Println("-----------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

// 返回 BlockChain 对象
func BlockChainObject() *BlockChain {
	// 创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip = b.Get([]byte("tip"))
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return &BlockChain{tip, db}
}
