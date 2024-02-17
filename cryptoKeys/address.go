package cryptoKeys

import "encoding/hex"

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

func NewAddressFromBString(s string) Address {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return Address{value: b}
}

// give public key from address
func (a Address) PublicKey() *PublicKey {
	return &PublicKey{
		Key: a.value,
	}
}