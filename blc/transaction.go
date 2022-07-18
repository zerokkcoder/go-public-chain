package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

// UXTO 交易模型
type Transaction struct {
	TxHash []byte      // 1. 交易 hash
	Vins   []*TXInput  // 2. 输入
	Vouts  []*TXOutput // 3. 输出
}

// 判断交易是否是创世区块交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

// 1. 创世区块创建时的 Transaction
func NewCoinbaseTransaction(address string) *Transaction {
	// 代表输入
	txInput := &TXInput{
		TxHash:    []byte{},
		Vout:      -1,
		Signature: nil,
		PublicKey: []byte{},
	}
	// 代表输出
	txOutput := NewTXOutput(10, address)
	txCoinbase := &Transaction{
		TxHash: []byte{},
		Vins:   []*TXInput{txInput},
		Vouts:  []*TXOutput{txOutput},
	}
	// 设置 TxHash 值
	txCoinbase.HashTransaction()

	return txCoinbase
}

// 2. 转账时产生的 Transaction
func NewSimpleTransaction(from string, to string, amount int, blockChain *BlockChain, txs []*Transaction) *Transaction {
	// 获取钱包
	wallets, _ := NewWallets()
	wallet := wallets.GetWallet(from)

	// 查找可用的UTXO
	money, spendableUTXODic := blockChain.FindSpendableUTXOs(from, amount, txs)
	var txInputs []*TXInput
	var txOutputs []*TXOutput
	// 消费
	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txInput := &TXInput{
				TxHash:    txHashBytes,
				Vout:      index,
				Signature: nil,
				PublicKey: wallet.PublicKey,
			}
			txInputs = append(txInputs, txInput)
		}
	}

	// 转账
	txOutput := NewTXOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)
	// 找零
	txOutput = NewTXOutput(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{
		TxHash: []byte{},
		Vins:   txInputs,
		Vouts:  txOutputs,
	}
	// 设置 TxHash 值
	tx.HashTransaction()

	return tx
}

// 将交易结构体序列化成字节数组
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}
