package transactions

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/OkO2451/BlockC/cryptoKeys"
)

type Transaction struct {
	ID        []byte
	Vin       []TXInput
	Vout      []TXOutput
	Signature *cryptoKeys.Signature
	PubKey    *cryptoKeys.PublicKey
	Value     int
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
	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{*txout},

		Value: subsidy,
	}
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

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       ScriptSig: %s", input.ScriptSig))
	}
	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       ScriptPubKey: %s", output.ScriptPubKey))
	}

	return strings.Join(lines, "\n")
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = cryptoKeys.Signature{Value: nil}
		txCopy.Vin[inID].SetPublicKey(prevTx.Vout[vin.Vout].PubKey)
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = cryptoKeys.PublicKey{Key: nil}

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			panic(err)
		}
		signature := &cryptoKeys.Signature{Value: append(r.Bytes(), s.Bytes()...)}

		tx.Vin[inID].Signature = *signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs,
			TXInput{
				Txid:      vin.Txid,
				Vout:      vin.Vout,
				Signature: cryptoKeys.Signature{Value: nil},
				PubKey:    cryptoKeys.PublicKey{Key: nil},
				ScriptSig: vin.ScriptSig})
	}
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{Value: vout.Value, ScriptPubKey: vout.ScriptPubKey})
	}

	txCopy := Transaction{ID: tx.ID, Vin: inputs, Vout: outputs, Value: tx.Value}

	return txCopy
}

func (tx *Transaction) Hash() []byte {
	txCopy := *tx
	txCopy.ID = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}
