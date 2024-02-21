package walet

import (
	"fmt"

	"github.com/OkO2451/BlockC/blockchain"
	"github.com/OkO2451/BlockC/cryptoKeys"

	// Bitcoin and many other cryptocurrencies use Base58
	// because it is designed to avoid visual ambiguity. 
	// o, O, 0, I, and l are not present in the character set,
	"github.com/mr-tron/base58"
)

type wlt struct {
	// address
	Address *cryptoKeys.Address

	// private key
	PrivateKey *cryptoKeys.PrivateKey

	// public key
	PublicKey *cryptoKeys.PublicKey

	// wallet file
	File string

	// wallet name
	Name string

	// wallet balance
	Balance int

	// wallet transactions
	Transactions []string
}

type Wallets struct {
	Wallets map[string]*wlt
}

func NewWallets() *Wallets {
	ws := Wallets{}
	ws.Wallets = make(map[string]*wlt)
	return &ws
}

func (ws *Wallets) CreateWallet() string {
	w := NewWallet()
	address := w.Address.String()
	ws.Wallets[address] = w
	return address
}

func NewWallet() *wlt {
	w := wlt{}
	var err error
	address, privateKey, err := cryptoKeys.NewKeyPair()
	if err != nil {
		fmt.Println("Error generating new key pair:", err)
		return nil
	}
	w.Address = &address
	w.PrivateKey = &privateKey
	w.PublicKey = w.Address.PublicKey()
	return &w
}

func (w *wlt) GetBalance() int {
	return w.Balance
}

func (w *wlt) GetTransactions() []string {
	// consult the bChain for transactions
	bc := blockchain.NewBlockchain(w.Address.String())
	defer bc.Db.Close()

	var transactions []string

	unspentTransactions := bc.FindUnspentTransactions(w.Address.String())

	for _, out := range unspentTransactions {
		transactions = append(transactions, out.String())
	}

	return transactions
}

// function to hashthepublic key
func (w *wlt) HashPubKey() []byte {
	return w.PublicKey.Hash()
}

func (w *wlt) GetAddress() []byte {
	pubKeyHash := cryptoKeys.HashPubKey(w.PublicKey.Bytes())

	versionedPayload := append([]byte{cryptoKeys.Version}, pubKeyHash...)
	checksum := cryptoKeys.Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)

	return []byte(address)
}
