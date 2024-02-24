package blockchain

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/OkO2451/BlockC/transactions"
)

type CLI struct {
	Bc *bChain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")

	fmt.Println("  addblock -data BLOCK_DATA - add a block to the bChain")
	fmt.Println("  addblock example: addblock -data \"Send 1 BTC to Ivan\"")

	fmt.Println("  printchain - print all the blocks of the bChain")
	fmt.Println("  printchain example: printchain")

	fmt.Println("  getbalance -address ADDRESS - get balance for ADDRESS")
	fmt.Println("  getbalance example: getbalance -address Ivan")
}

// Run executes the command-line interface (CLI) for the blockchain application.
// It parses the command-line arguments and performs the corresponding actions based on the provided commands.
func (cli *CLI) Run() {
	// code here
}

// Run executes the command-line interface (CLI) for the blockchain application.
// It parses the command-line arguments and performs the corresponding actions based on the provided commands.
func (cli *CLI) Run() {
	// code...
}
func (cli *CLI) Run() {
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")
	getBalanceData := getBalanceCmd.String("address", "", "The address to get balance for")

	// transactions := getTransactions()
	// should be implemented usig a private key
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			cli.printUsage()
			os.Exit(1)
		}
		if *addBlockData == "" {
			fmt.Println("addblock requires a -data flag")
			os.Exit(1)
		}
		cli.addBlock(*addBlockData, nil)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			cli.printUsage()
			os.Exit(1)
		}
		cli.printChain()
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			cli.printUsage()
			os.Exit(1)
		}
		if *getBalanceData == "" {
			fmt.Println("getbalance requires an -address flag")
			os.Exit(1)
		}
		balance := cli.getBalance(*getBalanceData)
		fmt.Printf("Balance of '%s': %d\n", *getBalanceData, balance)
	default:
		cli.printUsage()
		os.Exit(1)
	}
}

func NewCLI(Bc *bChain) *CLI {
	return &CLI{Bc}
}

func (cli *CLI) addBlock(data string, tr []*transactions.Transaction) {
	fmt.Println("In AddBlock function CLI")
	cli.Bc.AddBlock(data, tr)
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

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain()
	UTXOSet := UTXOSet{bc}
	defer bc.Db.Close()

	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)
	if tx != nil {
		cbTx := NewCoinbaseTX(from, "")
		txs := []*transactions.Transaction{cbTx, tx}

		bc.MineBlock(txs)
		fmt.Println("Success!")
	} else {
		fmt.Println("Failed to send transaction")
	}
}

func (cli *CLI) getBalance(address string) int {
	bc := NewBlockchain(address)
	defer bc.Db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	return balance
}
