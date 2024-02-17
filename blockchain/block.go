package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/OkO2451/BlockC/transactions"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Data          []byte
	Nonce         int // for proof of work algorithm
	TargetBits    int
	Transactions  []*transactions.Transaction
}

func (b *Block) String() string {
	return "Block -\n" +
		"Timestamp: " + strconv.FormatInt(b.Timestamp, 10) + "\n" +
		"PrevBlockHash: " + hex.EncodeToString(b.PrevBlockHash) + "\n" +
		"Hash: " + hex.EncodeToString(b.Hash) + "\n" +
		"Data: " + hex.EncodeToString(b.Data) + "\n" +
		"Nonce: " + strconv.Itoa(b.Nonce)
}

// create a new block
func NewBlock(data string, prevBlockHash []byte, Tr []*transactions.Transaction) *Block {

	fmt.Printf("In NewBlock data is: %v\n", data)
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Data:          []byte(data),
		Nonce:         0,
		Transactions:  Tr,
	}
	p := NewProofOfWork(block)
	fmt.Printf("After creating NewProofOfWork\n")

	nonce, hash := p.Run()
	fmt.Printf("After creating Run\n")

	block.Nonce = nonce
	block.Hash = hash[:]

	return block

}

// create a new genesis block
func NewGenesisBlock(coinbase *transactions.Transaction) *Block {
	return NewBlock("Genesis Block", []byte{}, []*transactions.Transaction{coinbase})
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}

	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	if err != nil {
		panic(err)
	}

	return &block
}

func (b *Block) IsValid() bool {
	pow := NewProofOfWork(b)
	return pow.Validate()
}

// hashTransactions
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (b *Block) getTransactions() []*transactions.Transaction {
	return b.Transactions
}

/*
func (b *Block) Serialize() []byte {
	var result []byte

	nonceBytes := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(nonceBytes, uint64(b.Nonce))

	targetBitsBytes := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(targetBitsBytes, uint64(b.TargetBits))

	timeBytes := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(timeBytes, uint64(b.Timestamp))

	result = append(result, timeBytes...)
	result = append(result, b.PrevBlockHash...)
	result = append(result, b.Hash...)
	result = append(result, b.Data...)
	result = append(result, nonceBytes...)
	result = append(result, targetBitsBytes...)
	return result
}


func DeserializeBlock(d []byte) *Block {
	var block Block

	timeBytes, timeBytes := binary.Uvarint(d)
	block.Timestamp = int64(timeBytes)

	block.PrevBlockHash = d[timeBytes : timeBytes+32]
	block.Hash = d[timeBytes+32 : timeBytes+64]
	block.Data = d[timeBytes+64 : timeBytes+96]

	nonceBytes, nonceBytes := binary.Uvarint(d[timeBytes+96 : timeBytes+128])
	block.Nonce = int(nonceBytes)

	targetBitsBytes, targetBitsBytes := binary.Uvarint(d[timeBytes+128 : timeBytes+160])
	block.TargetBits = int(targetBitsBytes)

	return &block
}
*/
