package blc

import (
	"fmt"
	"os"
)

// 转账
func (cli *CLI) send(from []string, to []string, amount []string) {
	if !DBExists() {
		fmt.Println("数据库不存在.......")
		os.Exit(1)
	}
	blockChain := BlockChainObject()
	defer blockChain.DB.Close()
	blockChain.MineNewBlock(from, to, amount)
}