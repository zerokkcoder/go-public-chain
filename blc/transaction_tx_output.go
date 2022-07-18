package blc

import "bytes"

type TXOutput struct {
	Value         int64
	Ripemd160Hash []byte // 公钥
}

// 判断是否时所属地址的 TXOutput
func (t *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	return bytes.Equal(t.Ripemd160Hash, pubKeyHash)
}

func NewTXOutput(value int64, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	// 设置 Ripemd160Hash
	txo.Lock(address)
	return txo
}

func (t *TXOutput) Lock(address string) {
	pubKeyHash := Base58Decode([]byte(address))
	t.Ripemd160Hash = pubKeyHash[1 : len(pubKeyHash)-4]
}
