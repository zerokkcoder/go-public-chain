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
func (w *Wallet) GetAddress() []byte {
	// 1. hash160
	// 20字节
	ripemd160Hash := w.Ripemd160Hash(w.PublicKey)
	// 21字节
	versionRipemd160Hash := append([]byte{version}, ripemd160Hash...)
	// 两次的256 hash
	checkSumBytes := CheckSum(versionRipemd160Hash)
	//25
	bytes := append(versionRipemd160Hash, checkSumBytes...)

	return Base58Encode(bytes)
}

func (w *Wallet) Ripemd160Hash(publicKey []byte) []byte {
	// 1. hash256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)
	// 2. hash160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil)
}

func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:addressChecksumLen]
}

// 判断一个钱包地址是否有效
func IsValidForAddress(address []byte) bool {
	versionPublicChecksumBytes := Base58Decode(address)

	checkSumBytes := versionPublicChecksumBytes[len(versionPublicChecksumBytes)-addressChecksumLen:]
	versionRipemd160 := versionPublicChecksumBytes[:len(versionPublicChecksumBytes)-addressChecksumLen]

	checkBytes := CheckSum(versionRipemd160)

	return bytes.Compare(checkSumBytes, checkBytes) == 0
}
