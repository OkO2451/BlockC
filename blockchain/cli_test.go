package blockchain

import (
	"testing"
)

func TestPrintUsage(t *testing.T) {

	bc := NewBlockchain()
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
