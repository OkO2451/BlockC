package blockchain

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/OkO2451/BlockC/cryptoKeys"
	"github.com/OkO2451/BlockC/transactions"
	"golang.org/x/crypto/ripemd160"

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

const walletFile = "wallets.data"

func NewWallets() *Wallets {
	// load the wallets from the file using the LoadFromFile function
	ws := Wallets{}
	ws.LoadFromFile(walletFile)
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

func NewUTXOTransaction(from, to string, amount int, bc *bChain) *transactions.Transaction {
	var inputs []transactions.TXInput
	var outputs []transactions.TXOutput

	// Create a Wallets object
	wallets := NewWallets()

	// Get the wallet and public key
	wlt := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wlt.PublicKey)

	acc, validOutputs := bc.FindSpendableOutputs(hex.EncodeToString(pubKeyHash), amount)

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
				Signature: cryptoKeys.Signature{Value: nil},
				PubKey:    *wlt.PublicKey,
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
	}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, *wlt.PrivateKey)

	return &tx
}

func (ws *Wallets) GetWallet(address string) *wlt {
	return ws.Wallets[address]
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

// function to load the wallet from the file

func (ws *Wallets) LoadFromFile(walletFile string) error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// function to save the wallet to the file
func (ws *Wallets) SaveToFile(walletFile string) {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func HashPubKey(pubKey *cryptoKeys.PublicKey) []byte {
	publicSHA256 := sha256.Sum256(pubKey.Bytes())

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}
