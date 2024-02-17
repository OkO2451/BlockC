package cryptoKeys

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type PrivateKey struct {
	// 64 bytes long
	Key ed25519.PrivateKey
}

func (p PrivateKey) bytes() []byte {
	return p.Key
}

func (p PrivateKey) String() string {
	return hex.EncodeToString(p.Key)
}

// sign creates a new Signature from data.
// It uses the PrivateKey to sign the data.
func (p PrivateKey) sign(data []byte) *Signature {
	return &Signature{
		Value: ed25519.Sign(p.Key, data),
	}
}

// NewPrivateKeyFromSeed creates a new PrivateKey from a seed.
// The seed should be 32 bytes long.
func NewPrivateKeyFromSeed(seed []byte) PrivateKey {
	return PrivateKey{
		Key: ed25519.NewKeyFromSeed(seed),
	}
}

// NewPrivateKeyFromBString creates a new PrivateKey from a hexadecimal string.
func NewPrivateKeyFromBString(s string) PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(b)
}

// GeneratePrivateKey generates a new PrivateKey.
// It uses a random 32-byte seed.
func GeneratePrivateKey() PrivateKey {
	seed := make([]byte, SeedLength)
	_, err := io.ReadFull(rand.Reader, seed)
	// if there is an error, panic because if the reader is not working,
	// the program is not working
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		Key: ed25519.NewKeyFromSeed(seed),
	}
}

func (p *PrivateKey) Bytes() []byte {
	return p.Key
}
