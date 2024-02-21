// Package cryptoKeys provides functionality related to cryptographic keys and addresses.
package cryptoKeys

import (
	"crypto/ed25519" // Importing the ed25519 package for generating and using Ed25519 keys.
	"crypto/rand"    // Importing the rand package for generating random numbers.
	"crypto/sha256"
	"encoding/hex" // Importing the hex package for hexadecimal encoding/decoding.
)

const Version = byte(0x00) // version is the version byte of the address.
const addressChecksumLen = 4

// Address represents a cryptographic address,
// which is a hashed version of a public key.
type Address struct {
	// value is the byte slice representation of the address.
	value []byte
}

// String returns the hexadecimal string representation of the address.
func (a *Address) String() string {
	return hex.EncodeToString(a.value)
}

// Bytes returns the byte slice representation of the address.
func (a *Address) Bytes() []byte {
	return a.value
}

// Length returns the length of the byte slice representation of the address.
func (a *Address) Length() int {
	return len(a.value)
}

// NewAddressFromBString creates a new Address from a hexadecimal string.
func NewAddressFromBString(s string) Address {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return Address{value: b}
}

// PublicKey returns the PublicKey corresponding to the Address.
func (a Address) PublicKey() *PublicKey {
	return &PublicKey{
		Key: a.value,
	}
}

// Address returns the Address corresponding to the PublicKey.
func (p PublicKey) Address() *Address {
	return &Address{
		value: p.Key[:AddressLength],
	}
}

func HashPubKey(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)
	return hash[:]
}

func NewKeyPair() (Address, PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return Address{}, PrivateKey{}, err
	}

	// Hash the public key to generate the address
	hashedPubKey := HashPubKey(publicKey)

	return Address{value: hashedPubKey}, PrivateKey{Key: privateKey}, nil
}

// checksum creates a checksum of the payload.
func Checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
