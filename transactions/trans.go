package transactions

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

const subsidy = 1

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{
		Txid:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.SetID()

	return &tx
}

func (tx *Transaction) SetID() {
	var encoded []byte
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}
	encoded = txCopy.Serialize()
	hash = sha256.Sum256(encoded)
	txCopy.ID = hash[:]
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		fmt.Println("Error in Serialize Transaction:")
		fmt.Println(err)
	}

	return encoded.Bytes()
}
