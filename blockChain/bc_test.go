package blockchain

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain()
	assert.Equal(t, len(bc.blocks), 1)
	assert.Equal(t, string(bc.blocks[0].Data), "Genesis Block")
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain()
	data := "Test Block"
	bc.AddBlock(data)

	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		t.Errorf("View error: %v", err)
	}

	var lastBlock *Block

	err = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(lastHash)
		lastBlock = DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		t.Errorf("View error: %v", err)
	}
	fmt.Printf("data is : %s\n", data)
	fmt.Printf("lastBlock.Data is : %s\n", lastBlock.Data)
	fmt.Printf("data to hex: %x\n", data)

	assert.Equal(t, string(lastBlock.Data), data)
}

// test for new proof of work
func TestNewProofOfWork(t *testing.T) {
	bc := NewBlockchain()
	pow := NewProofOfWork(bc.blocks[0])
	fmt.Println(len(pow.prepareData(0)))
	assert.Equal(t, pow.block, bc.blocks[0])
	assert.Equal(t, pow.target, pow.target)
}

// test for prepare data
func TestPrepareData(t *testing.T) {
	bc := NewBlockchain()
	pow := NewProofOfWork(bc.blocks[0])
	data := pow.prepareData(0)

	// Check that data is not nil
	assert.NotNil(t, data)

	// Check that data contains the expected fields
	assert.Contains(t, string(data), string(bc.blocks[0].PrevBlockHash))
}

// test for run
func TestRun(t *testing.T) {

	bc := NewBlockchain()
	fmt.Printf("After creating In TestRun: %v\n", bc.blocks[0])
	bc.AddBlock("Pay 1 coin for coffee")
	bc.AddBlock("Send 2 more BTC to mok")
	pow := NewProofOfWork(bc.blocks[0])
	fmt.Printf("After creating In TestRun: %v\n", pow)
	nonce, hash := pow.Run()
	assert.NotNil(t, nonce)
	assert.NotNil(t, hash)
}

// test for validate all of the blocks
func TestIsValid(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("Send 1 BTC to bak")
	bc.AddBlock("Send 2 more BTC to bak #2")
	assert.True(t, bc.IsValid())
}

func TestBCIterator_Next(t *testing.T) {
	// Create a new blockchain
	bc := NewBlockchain()
	// Add a block to the blockchain
	bc.AddBlock("Test Block")

	// Create a new blockchain iterator
	it := bc.Iterator()

	// Get the next block in the blockchain
	block := it.Next()

	// Check that the block is not nil
	assert.NotNil(t, block)

	// Check that the block's data is correct
	assert.Equal(t, "Test Block", string(block.Data))

	// Get the next block in the blockchain
	block = it.Next()

	// Check that the block is the genesis block
	assert.Equal(t, "Genesis Block", string(block.Data))

	// Try to get another block from the blockchain
	block = it.Next()

	// Check that there are no more blocks in the blockchain
	assert.Nil(t, block)
}

func TestPrintAllBlocks(t *testing.T) {

	// make a random seed

	// Create a new blockchain
	bc := NewBlockchain()

	// take time in seconds in a string
	str := strings.Split(time.Now().String(), " ")[0]

	// Add a block to the blockchain

	bc.AddBlock("Test Block " + str)

	// Print all blocks in the blockchain
	PrintAllBlocks(bc)
}
