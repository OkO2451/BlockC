package transactions

import (
	"bytes"

	"github.com/OkO2451/BlockC/cryptoKeys"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
	Signature cryptoKeys.Signature
	PubKey    cryptoKeys.PublicKey
}

func NewTXInput(txid []byte, vout int, scriptSig string, pubkey cryptoKeys.PublicKey ) *TXInput {
	txin := &TXInput{
		Txid:      txid,
		Vout:      vout,
		ScriptSig: scriptSig,
	}
	txin.SetPublicKey(pubkey)

	return txin
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
    return in.ScriptSig == unlockingData
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
    lockingHash := cryptoKeys.HashPubKey(in.PubKey.Bytes())
    return bytes.Equal(lockingHash, pubKeyHash)
}

func (in *TXInput) SetPublicKey(pubKey cryptoKeys.PublicKey) {
    in.PubKey = pubKey
}

func (in *TXInput) GetPublicKey() cryptoKeys.PublicKey {
    return in.PubKey
}

func (in *TXInput) SetSignature(sig cryptoKeys.Signature) {
	in.Signature = sig
}

