package blc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

// Wallet 存储 private 和 public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey // 1. 私钥
	PublicKey  []byte           // 2. 公钥
}

// 创建一个钱包
func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

// 通过私钥产生一个公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// 1.
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// // 返回钱包地址
// func (w Wallet) GetAddress() []byte {
// 	pubKeyHash := HashPub
// }
