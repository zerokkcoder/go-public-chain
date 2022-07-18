package blc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

// Wallet 存储 private 和 public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey // 1. 私钥
	PublicKey  []byte           // 2. 公钥
}

// 判断一个钱包地址是否有效
func IsValidForAddress(address []byte) bool {
	version_public_checksumBytes := Base58Decode(address)

	//25
	//4
	//21
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes)-addressChecksumLen:]

	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes)-addressChecksumLen]

	checkBytes := CheckSum(version_ripemd160)

	return bytes.Equal(checkSumBytes, checkBytes)
}

// 返回钱包地址
func (w *Wallet) GetAddress() []byte {
	// 1. hash160
	// 20字节
	ripemd160Hash := w.Ripemd160Hash(w.PublicKey)
	// 21字节
	version_ripemd160Hash := append([]byte{version}, ripemd160Hash...)
	// 两次的256 hash
	checkSumBytes := CheckSum(version_ripemd160Hash)
	//25
	bytes := append(version_ripemd160Hash, checkSumBytes...)

	return Base58Encode(bytes)
}

func CheckSum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen]
}

func (w *Wallet) Ripemd160Hash(publicKey []byte) []byte {
	//1. 256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	//2. 160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil)
}

// 创建一个钱包
func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{privateKey, publicKey}
}

// 通过私钥产生一个公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
