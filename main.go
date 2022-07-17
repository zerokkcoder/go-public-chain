package main

import (
	"fmt"

	"golang.org/x/crypto/ripemd160"
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
	hasher := ripemd160.New()
	hasher.Write([]byte("huanghuanghuang"))
	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	fmt.Println(hashString)

}
