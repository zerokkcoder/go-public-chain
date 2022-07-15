package blc

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	BlockChain *BlockChain
}

func (cli *CLI) Run() {
	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createGenesisBlockCmd := flag.NewFlagSet("creategenesisblock", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "http://xxx.com", "交易数据")
	flagCreateGenesisBlockData := createGenesisBlockCmd.String("data", "Genesis data ...", "创世区块数据")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
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
	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createGenesisBlockCmd.Parsed() {
		if *flagCreateGenesisBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockChain(*flagCreateGenesisBlockData)
	}
}

func (cli *CLI) addBlock(data string) {
	cli.BlockChain.AddBlockToBlockChain(data)
}

func (cli *CLI) printChain() {
	cli.BlockChain.PrintChain()
}

func (cli *CLI) createGenesisBlockChain(data string) {
	CreateBlockChainWithGenesisBlock(data)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreategenesisblock -data DATA -- 创建创世区块")
	fmt.Println("\taddblock -data DATA -- 交易数据")
	fmt.Println("\tprintchain -- 输出区块信息")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
