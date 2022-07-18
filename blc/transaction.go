package blc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
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

	// 进行数字签名
	blockChain.SignTransaction(tx, wallet.PrivateKey)

	return tx
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()

}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.TxHash = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// 签名
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.IsCoinbaseTransaction() {
		return
	}

	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR:Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vins {
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[inID].Signature = signature
	}
}

// 拷贝一份新的 Transaction 用于签名
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins {
		inputs = append(inputs, &TXInput{
			TxHash:    vin.TxHash,
			Vout:      vin.Vout,
			Signature: vin.Signature,
			PublicKey: vin.PublicKey,
		})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TXOutput{
			Value:         vout.Value,
			Ripemd160Hash: vout.Ripemd160Hash,
		})
	}

	txCopy := Transaction{
		TxHash: tx.TxHash,
		Vins:   inputs,
		Vouts:  outputs,
	}
	return txCopy
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

// 数字签名验证
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 私钥 TxHash
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])

		rawPublicKey := ecdsa.PublicKey{curve, &x, &y}
		if !ecdsa.Verify(&rawPublicKey, txCopy.TxHash, &r, &s) {
			return false
		}
	}
	return true
}
