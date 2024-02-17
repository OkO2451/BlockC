package cryptoKeys

import "crypto/ed25519"

type PublicKey struct {
	// 32 bytes long
	Key ed25519.PublicKey
}

func (p *PublicKey) Bytes() []byte {
	return p.Key
}

// PublicKey returns the PublicKey corresponding to the PrivateKey.
func (p PrivateKey) PublicKey() PublicKey {
	return PublicKey{Key: p.Key.Public().(ed25519.PublicKey)}
}

func (p PublicKey) Public() *PublicKey {
	return &PublicKey{
		Key: p.Key,
	}
}

// Address returns the Address corresponding to the PublicKey.
func (p PublicKey) Address() *Address {
	return &Address{
		value: p.Key[:AddressLength],
	}
}
