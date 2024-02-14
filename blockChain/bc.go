package blockchain

import (
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

const dbFile = "bc.db"
const blocksBucket = "bChains"

type blockchain struct {
	tip    []byte
	blocks []*Block
	db     *bolt.DB
}

// iterator pattern
// because we are using bolt db, we need to iterate through the blocks
type bcIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// create a new blockchain
func NewBlockchain() *blockchain {

	var tip []byte
	db, _ := bolt.Open(dbFile, 0600, nil)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, _ := tx.CreateBucket([]byte(blocksBucket))
			b.Put(genesis.Hash, genesis.Serialize())
			b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := blockchain{
		tip: tip,
		db:  db,
	}

	return &bc
}

// add a new block to the blockchain
func (bc *blockchain) AddBlock(data string) {

	fmt.Println("In AddBlock")

	var lastHash []byte

	fmt.Printf("LastHash: %v\n", lastHash)

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		fmt.Printf("LastHash: %v\n", lastHash)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})

}

// very expensive operation to check if the blockchain is valid
func (bc *blockchain) IsValid() bool {
	valFlag := true

	for _, block := range bc.blocks {

		pow := NewProofOfWork(block)
		valFlag = valFlag && pow.Validate()
		fmt.Printf("PoW: %s\n", strconv.FormatBool(valFlag))
		fmt.Println()

	}
	return valFlag
}

// create a new iterator
func (bc *blockchain) Iterator() *bcIterator {
	bci := &bcIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}

	return bci
}
func (i *bcIterator) Next() *Block {
    var block *Block

    err := i.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(blocksBucket))
        encodedBlock := b.Get(i.currentHash)

        if encodedBlock == nil {
            return nil
        }

        block = DeserializeBlock(encodedBlock)

        return nil
    })
    if err != nil {
        log.Panic(err)
    }

    if block != nil {
        i.currentHash = block.PrevBlockHash
    }

    return block
}
func PrintAllBlocks(bc *blockchain) {
    // Create a new blockchain iterator
    it := bc.Iterator()

    // Iterate over all blocks in the blockchain
    for {
        block := it.Next()

        // Break if there are no more blocks in the blockchain
        if block == nil {
            break
        }

        // Print the block's data
        fmt.Println(string(block.Data))
    }
}