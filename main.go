package main

import (
	"fmt"
	"go-public-chain/blc"
)

func main() {
	// cli := blc.CLI{}
	// cli.Run()

	// 256
	// 47cf44ab1c182aeb1660642934317981e4041d5ce994d3f339115361434de879
	// hasher := sha256.New()
	// hasher.Write([]byte("huanghuanghuang"))
	// hashBytes := hasher.Sum(nil)
	// hashString := fmt.Sprintf("%x", hashBytes)
	// fmt.Println(hashString)

	// ripemd160
	// 4c7da9e96f8ec28020b68f6bfe6ffff8e63a3598
	// hasher := ripemd160.New()
	// hasher.Write([]byte("huanghuanghuang"))
	// hashBytes := hasher.Sum(nil)
	// hashString := fmt.Sprintf("%x", hashBytes)
	// fmt.Println(hashString)

	// base58
	bytes := []byte("huanghuanghuang")
	encode := blc.Base58Encode(bytes)
	fmt.Printf("%x\n", encode)
	fmt.Printf("%s\n", encode)
	decode := blc.Base58Decode(encode)
	fmt.Printf("%x\n", decode)
	fmt.Printf("%s\n", decode[1:])
}
