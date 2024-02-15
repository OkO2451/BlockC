package main

import (
	"github.com/OkO2451/BlockC/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.Db.Close()

	cli := blockchain.CLI{
		Bc: bc,
	}
	cli.Run()
}
