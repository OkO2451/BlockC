package transactions

type TXOutput struct {
	Value        int
	ScriptPubKey string
}


func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, address}
	return txo
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}