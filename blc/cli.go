package blc

import (
	"flag"
	"fmt"
	"go-public-chain/utils"
	"log"
	"os"
)

type CLI struct {
	BlockChain *BlockChain
}

func (cli *CLI) Run() {
	isValidArgs()

	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createGenesisBlockCmd := flag.NewFlagSet("creategenesisblock", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账源地址...")
	flagTo := sendBlockCmd.String("to", "", "转账目的地址...")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额址...")

	flagCreateGenesisBlockAddress := createGenesisBlockCmd.String("address", "", "创建创世区块的地址")

	getBalanceWithAddress := getBalanceCmd.String("address", "", "查询该地址的余额")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "creategenesisblock":
		err := createGenesisBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}
		cli.send(utils.JSONToArray(*flagFrom), utils.JSONToArray(*flagTo), utils.JSONToArray(*flagAmount))
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createGenesisBlockCmd.Parsed() {
		if *flagCreateGenesisBlockAddress == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockChain(*flagCreateGenesisBlockAddress)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceWithAddress == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceWithAddress)
	}
}

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

// 创建创世区块
func (cli *CLI) createGenesisBlockChain(address string) {
	blockChain := CreateBlockChainWithGenesisBlock(address)
	defer blockChain.DB.Close()
}

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

// 获取余额
func (cli *CLI) getBalance(address string) {
	fmt.Println("地址:" + address)
	txs := UnSpentTransactionWithAddress(address)

	fmt.Println(txs)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreategenesisblock -address ADDRESS -- 创建创世区块")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细")
	fmt.Println("\tprintchain -- 输出区块信息")
	fmt.Println("\tgetbalance -address ADDRESS -- 获取账户余额")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
