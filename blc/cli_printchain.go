package blc

import (
	"fmt"
	"os"
)

// 打印区块链
func (cli *CLI) printChain() {
	if !DBExists() {
		fmt.Println("数据库不存在.......")
		os.Exit(1)
	}
	blockChain := BlockChainObject()
	defer blockChain.DB.Close()
	blockChain.PrintChain()
}
