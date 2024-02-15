package transactions

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}


func NewTXInput(txid []byte, vout int, scriptSig string) *TXInput {
	txin := &TXInput{txid, vout, scriptSig}
	return txin
}