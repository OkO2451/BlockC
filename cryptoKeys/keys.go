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

const (
	// Length of the private key
	PrivateKeyLength = 64
	// Length of the public key
	PublicKeyLength = 32

	// length of the seed
	SeedLength = 32

	// Length of the signature
	SignatureLength = 64

	// Length of the address
	AddressLength = 20
)

func (p PrivateKey) bytes() []byte {
	return p.Key
}

func (p PrivateKey) String() string {
	return hex.EncodeToString(p.Key)
}

// makig these Proto buffer types
func (p PrivateKey) sign(data []byte) *Signature {
	return &Signature{
		Value: ed25519.Sign(p.Key, data),
	}
}

func NewPrivateKeyFromSeed(seed []byte) PrivateKey {
	return PrivateKey{
		Key: ed25519.NewKeyFromSeed(seed),
	}
}

func NewPrivateKeyFromBString(s string) PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(b)
}

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

type PublicKey struct {
	// 32 bytes long
	Key ed25519.PublicKey
}

func (p *PublicKey) Bytes() []byte {
	return p.Key
}

func (p PrivateKey) PublicKey() PublicKey {
	return PublicKey{Key: p.Key.Public().(ed25519.PublicKey)}
}

func (p PublicKey) Public() *PublicKey {
	return &PublicKey{
		Key: p.Key,
	}
}

func (p PublicKey) Address() *Address {
	return &Address{
		value: p.Key[:AddressLength],
	}
}

type Signature struct {
	// 64 bytes long
	Value []byte
}

func (s *Signature) Bytes() []byte {
	return s.Value
}

func (s *Signature) Verify(data []byte, pub PublicKey) bool {
	return ed25519.Verify(pub.Key, data, s.Value)
}

type Address struct {
	// 32 bytes long
	value []byte
}

func (a *Address) String() string {
	return hex.EncodeToString(a.value)
}

func (a *Address) Bytes() []byte {
	return a.value
}

func (a *Address) Length() int {
	return len(a.value)
}
