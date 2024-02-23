package cryptoKeys

import (
	"crypto/ed25519"
	"crypto/elliptic"
)

type Signature struct {
	// 64 bytes long
	Value []byte
}

func (s *Signature) Bytes() []byte {
	return s.Value
}

// Verify checks whether the Signature is a valid signature of data.
// It uses the PublicKey to verify the signature.
func (s *Signature) Verify(data []byte, pub PublicKey) bool {
	return ed25519.Verify(pub.Key, data, s.Value)
}

func Curve() elliptic.Curve {
	return elliptic.P256()
}
