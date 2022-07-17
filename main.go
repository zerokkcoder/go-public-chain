package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	// cli := blc.CLI{}
	// cli.Run()

	hasher := sha256.New()
	hasher.Write([]byte("huanghuanghuang"))
	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	fmt.Println(hashString)

}
