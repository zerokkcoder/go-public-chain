package blc

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"go-public-chain/utils"
	"math/big"
)

// 256位 Hash里面前面至少要有16个零
const targetBit = 16

type ProofOfWork struct {
	Block  *Block   // 验证的区块
	target *big.Int // 大数据存储，代表挖矿难度
}

// 创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	// 1. 创建一个初始值为1的target
	target := big.NewInt(1)
	// 2. 右移 256-targetBit
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{block, target}
}

// 判断 Hash 是否有效
func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int
	hashInt.SetBytes(pow.Block.Hash)
	return pow.target.Cmp(&hashInt) == 1
}

func (pow ProofOfWork) Run() ([]byte, int64) {
	// 1. 将 Block的属性拼接成字节数组
	// 2. 生成 Hash
	// 3. 判断 hash 有效性， 如果满足条件，跳出循环
	nonce := 0
	var hashInt big.Int // 存储新生成的 hash
	var hash [32]byte
	for {
		// 准备数据
		dataBytes := pow.prepareData(int64(nonce))
		// 生成 Hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)
		// 将 hash存储到hashInt
		hashInt.SetBytes(hash[:])
		// 判断hashInt是否小于Block里面的target
		if pow.target.Cmp(&hashInt) == 1 {
			break
		}

		nonce++
	}

	return hash[:], int64(nonce)
}

// 数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.Data,
		utils.IntToHex(pow.Block.Timestamp),
		utils.IntToHex(int64(targetBit)),
		utils.IntToHex(int64(nonce)),
		utils.IntToHex(int64(pow.Block.Height)),
	}, []byte{})

	return data
}
