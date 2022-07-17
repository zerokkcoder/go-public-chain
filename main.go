package main

import (
	"fmt"
	"go-public-chain/blc"
)

func main() {
	// cli := blc.CLI{}
	// cli.Run()

	wallet := blc.NewWallet()
	// address := wallet.GetAddress()
	isValid := wallet.IsValidForAddress([]byte("1MUhdPZJHDLvxRuXiX3cNBLWLW7jfa2dp"))
	fmt.Printf("1MUhdPZJHDLvxRuXiX3cNBLWLW7jfa2dp 这个地址为 %v\n", isValid)
}
