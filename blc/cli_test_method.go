package blc

import "fmt"

func (cli *CLI) TestMethod() {

	blockChain := BlockChainObject()
	defer blockChain.DB.Close()

	utxoMap := blockChain.FindUTXOMap()
	fmt.Println(utxoMap)
}
