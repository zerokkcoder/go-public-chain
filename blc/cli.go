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

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddresslist -- 输出所有钱包地址")
	fmt.Println("\tcreatewallet -- 创建钱包")
	fmt.Println("\tcreategenesisblock -address ADDRESS -- 创建创世区块")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细")
	fmt.Println("\tprintchain -- 输出区块信息")
	fmt.Println("\tgetbalance -address ADDRESS -- 获取账户余额")
	fmt.Println("\ttest -- 测试")
}

func (cli *CLI) Run() {
	isValidArgs()

	// 获取节点ID
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!\n")
		os.Exit(1)
	}
	fmt.Printf("NODE_ID: %s\n", nodeID)

	addressListCmd := flag.NewFlagSet("addresslist", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createGenesisBlockCmd := flag.NewFlagSet("creategenesisblock", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	testCmd := flag.NewFlagSet("test", flag.ExitOnError)

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
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addresslist":
		err := addressListCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "test":
		err := testCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
	if addressListCmd.Parsed() {
		// 输出所有钱包地址
		cli.addressList(nodeID)
	}

	if createWalletCmd.Parsed() {
		// 创建钱包
		cli.createWallet(nodeID)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}

		from := utils.JSONToArray(*flagFrom)
		to := utils.JSONToArray(*flagTo)

		for index, fromAddress := range from {
			if !IsValidForAddress([]byte(fromAddress)) || !IsValidForAddress([]byte(to[index])) {
				fmt.Println("地址无效....")
				printUsage()
				os.Exit(1)
			}
		}

		amount := utils.JSONToArray(*flagAmount)

		cli.send(from, to, amount,nodeID)
	}

	if printChainCmd.Parsed() {
		cli.printChain(nodeID)
	}

	if createGenesisBlockCmd.Parsed() {
		if !IsValidForAddress([]byte(*flagCreateGenesisBlockAddress)) {
			fmt.Println("地址无效....")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockChain(*flagCreateGenesisBlockAddress, nodeID)
	}

	if getBalanceCmd.Parsed() {
		if !IsValidForAddress([]byte(*getBalanceWithAddress)) {
			fmt.Println("地址无效....")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceWithAddress, nodeID)
	}
	if testCmd.Parsed() {
		cli.TestMethod(nodeID)
	}
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
