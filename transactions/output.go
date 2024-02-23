package transactions

import (
	"bytes"

	"github.com/OkO2451/BlockC/cryptoKeys"
	"github.com/mr-tron/base58"
)

type TXOutput struct {
	Value        int
	ScriptPubKey string
	PubKey 	 cryptoKeys.PublicKey
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{Value: value, ScriptPubKey: address}
	return txo
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func (out *TXOutput) Lock(address cryptoKeys.Address) {
	pubKeyhash := address.Bytes()
	pubk := base58.Encode(pubKeyhash)
	pubk = pubk[1 : len(pubk)-4]
	out.ScriptPubKey = string(pubk)
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	lockingHash := cryptoKeys.HashPubKey([]byte(out.ScriptPubKey))
	return bytes.Equal(lockingHash, pubKeyHash)
}


