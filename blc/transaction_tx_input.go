package blc

import "bytes"

type TXInput struct {
	TxHash    []byte // 1. 交易的Hash
	Vout      int    // 2. 存储 TXOutput 在 Vout里面的索引
	Signature []byte // 3. 数字签名
	PublicKey []byte // 4. 公钥，钱包里面的公钥
}

// 判断是否时所属地址的 TXInput, 判断当前输入是否和某个输出吻合
func (t *TXInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	publicKeyHash := Ripemd160Hash(t.PublicKey)
	return bytes.Equal(ripemd160Hash, publicKeyHash)
}
