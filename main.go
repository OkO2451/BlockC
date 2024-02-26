package main

import (
	"net/http"

	"github.com/OkO2451/BlockC/blockchain"
	"github.com/OkO2451/BlockC/transactions"
	"github.com/gin-gonic/gin"
)

type BlockchainAPI struct {
	Bc *blockchain.BChain
}

// main is the entry point of the program.
// It initializes a blockchain with a given address,
// creates an instance of the BlockchainAPI,
// and sets up the routes for adding a block, printing the chain, and getting the balance.
// Finally, it starts the server to listen and serve on 0.0.0.0:8080.
func main() {
	address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa" // the first address in the blockchain
	bc := blockchain.NewBlockchain(address)
	api := &BlockchainAPI{Bc: bc}

	r := gin.Default()

	r.POST("/addblock", api.addBlock)
	r.GET("/printchain", api.printChain)
	r.GET("/getbalance/:address", api.getBalance)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func (api *BlockchainAPI) addBlock(c *gin.Context) {
	data := c.PostForm("data")
	tr := []*transactions.Transaction{} // TODO: Get transactions from somewhere
	api.Bc.AddBlock(data, tr)
	c.JSON(http.StatusOK, gin.H{"message": "Block added"})
}

func (api *BlockchainAPI) printChain(c *gin.Context) {
	bci := api.Bc.Iterator()
	blocks := []blockchain.Block{}

	for {
		block := bci.Next()
		blocks = append(blocks, *block)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	c.JSON(http.StatusOK, blocks)
}

func (api *BlockchainAPI) getBalance(c *gin.Context) {
	address := c.Param("address")
	balance := 0
	UTXOs := api.Bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
