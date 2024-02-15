package main

import (
	"github.com/OkO2451/BlockC/blockchain"
	
)

func main() {
	bc := NewBlockchain()
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
