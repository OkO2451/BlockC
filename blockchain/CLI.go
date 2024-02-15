package blockchain

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	Bc *blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) Run() {

	

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")


	fmt.Printf("Your command is: %v with flag: %v\n", os.Args[1], os.Args[2])

	fmt.Printf("We are processing the command line\n")
	if len(os.Args) < 2 {
		cli.printUsage()
		fmt.Printf("You have not entered a command\n")
		os.Exit(1)
	} else {
		fmt.Printf("You have entered a command\n")
	}

	fmt.Printf("Your command is: %v\n", os.Args[1])

	switch os.Args[1] {
	case "addblock":
		fmt.Printf("You have entered addblock\n")
		addBlockCmd.Parse(os.Args[2:])
		if addBlockCmd.Parsed() {
        if *addBlockData == "" {
            fmt.Println("addblock requires a -data flag")
        } else {
            fmt.Printf("Block data: %s\n", *addBlockData)
            cli.addBlock(*addBlockData)
        }
    }
	case "printchain":
		fmt.Printf("You have entered printchain\n")
		printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		
	}

	

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func NewCLI(Bc *blockchain) *CLI {
	return &CLI{Bc}
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	fmt.Println("In AddBlock function CLI")
	cli.Bc.AddBlock(data)
	fmt.Println("Block added")
}

func (cli *CLI) printChain() {
	bci := cli.Bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
