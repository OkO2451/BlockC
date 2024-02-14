package cryptoKeys

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go get -u github.com/stretchr/testify for the assert package

func TestGeneratePrivateKey(t *testing.T) {
	priv := GeneratePrivateKey()
	assert.Equal(t, len(priv.Bytes()), PrivateKeyLength)

	publi := priv.PublicKey()
	assert.Equal(t, len(publi.Bytes()), PublicKeyLength)
}

// testing the sign and verify functions
func TestSignAndVerify(t *testing.T) {
	priv := GeneratePrivateKey()
	publi := priv.PublicKey()
	data := []byte("test test 1 2 3")
	signa := priv.sign(data)

	assert.True(t, signa.Verify(data, publi))

	// test with wrong data
	assert.False(t, signa.Verify([]byte("wrong data"), publi))

	// test with wrong signature
	assert.False(t, (&Signature{Value: []byte("wrong signature")}).Verify(data, publi))

	// test with wrong public key
	priv2 := GeneratePrivateKey()
	publi2 := priv2.PublicKey()
	assert.False(t, signa.Verify(data, publi2))
}

func TestPublicKeyAddress(t *testing.T) {
	priv := GeneratePrivateKey()
	publi := priv.PublicKey()
	assert.Equal(t, publi.Address().Length(), AddressLength)
}

func TestPublicKeyToAddress(t *testing.T) {
	priv := GeneratePrivateKey()
	publi := priv.PublicKey()
	addr := publi.Address()
	assert.Equal(t, addr.Length(), AddressLength)
	fmt.Println(addr)
}

func TestPrivateKeyFromSeed(t *testing.T) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)
	fmt.Println(seed)
	fmt.Println("ok")
	fmt.Println(hex.EncodeToString(seed))
}

func TestPrivateKeyFromBString(t *testing.T) {
	var (
		addressStr = "11c1db7da340b949d8fe34954c18ffa6164529eb"
		seed = "2624f1d473599abe62fdeebe200fd221a3891f2dc098d4b4057610f2d58e737b"
		privKey = NewPrivateKeyFromBString(seed)
	)
	
	assert.Equal(t, PrivateKeyLength, len(privKey.Bytes()))
	adress := privKey.PublicKey().Address() 
	assert.Equal(t, adress.String() , addressStr)

}
