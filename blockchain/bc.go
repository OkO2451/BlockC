package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/OkO2451/BlockC/cryptoKeys"
	"github.com/OkO2451/BlockC/transactions"
	"github.com/boltdb/bolt"
)

const dbFile = "bc.db"
const blocksBucket = "bChains"
const genesisCoinbaseData = "banks are the modern day robber baron"

type bChain struct {
	Tip []byte

	Db *bolt.DB
}

// iterator pattern
// because we are using bolt Db, we need to iterate through the blocks
type bcIterator struct {
	currentHash []byte
	Db          *bolt.DB
}

// create a new bChain
func NewBlockchain(address string) *bChain {

	var Tip []byte
	Db, _ := bolt.Open(dbFile, 0600, nil)

	Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			coin := transactions.NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(coin)
			b, _ := tx.CreateBucket([]byte(blocksBucket))
			b.Put(genesis.Hash, genesis.Serialize())
			b.Put([]byte("l"), genesis.Hash)
			Tip = genesis.Hash
		} else {
			Tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := bChain{
		Tip: Tip,
		Db:  Db,
	}

	return &bc
}

// add a new block to the bChain
func (bc *bChain) AddBlock(data string, tr []*transactions.Transaction) {
	fmt.Println("In AddBlock")

	var lastHash []byte

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash, tr)

	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}
		bc.Tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// very expensive operation to check if the bChain is valid
func (bc *bChain) IsValid() bool {
	it := bc.Iterator()

	for {
		block := it.Next()

		if block == nil {
			break
		}

		if !block.IsValid() {
			return false
		}

	}

	return true
}

// create a new iterator
func (bc *bChain) Iterator() *bcIterator {
	bci := &bcIterator{
		currentHash: bc.Tip,
		Db:          bc.Db,
	}

	return bci
}
func (i *bcIterator) Next() *Block {
	var block *Block

	err := i.Db.View(func(tx *bolt.Tx) error {
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
func PrintAllBlocks(bc *bChain) {
	// Create a new bChain iterator
	it := bc.Iterator()

	// Iterate over all blocks in the bChain
	for {
		block := it.Next()

		// Break if there are no more blocks in the bChain
		if block == nil {
			break
		}

		// Print the block's data
		fmt.Println(string(block.Data))
	}
}

// add a private key to the bChain
func (bc *bChain) AddPrivateKey(privateKey string) {
	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put([]byte("privateKey"), []byte(privateKey))
		return err
	})
	if err != nil {
		log.Panic(err)
	}
}

func (bc *bChain) FindUnspentTransactions(address string) []transactions.Transaction {
	var unspentTXs []transactions.Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if !tx.IsCoinbase() {

				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}

		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *bChain) FindUTXO(address string) []transactions.TXOutput {
	var UTXOs []transactions.TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func NewUTXOTransaction(from, to string, amount int, bc *bChain) *transactions.Transaction {
	var inputs []transactions.TXInput
	var outputs []transactions.TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		fmt.Println("Error: Not enough funds")
		return nil
	}

	for txid, outs := range validOutputs {
		txID := []byte(txid)
		for _, out := range outs {
			input := transactions.TXInput{
				Txid:      txID,
				Vout:      out,
				ScriptSig: from,
			}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *transactions.NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *transactions.NewTXOutput(acc-amount, from))
	}

	tx := transactions.Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outputs,

		Value: amount,
	}
	tx.SetID()

	return &tx
}

func (bc *bChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (bc *bChain) FindTransaction(ID []byte) (transactions.Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return transactions.Transaction{}, nil
}

func (bc *bChain) SignTransaction(tx *transactions.Transaction, privKey cryptoKeys.PrivateKey) {
	prevTXs := make(map[string]transactions.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *bChain) VerifyTransaction(tx *transactions.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]transactions.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}