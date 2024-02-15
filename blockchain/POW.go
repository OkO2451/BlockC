package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

var (
	// it should be math.MaxInt64
	// but i put a smaller value for testing
	maxNonce = math.MaxInt64
)

// to make it harder to mine a block
// we can increase the targetBits
// the larger the targetBits, the harder it is to mine a block
const targetBits = 24

// ProofOfWork represents a proof-of-work
type pow struct {
	block  *Block
	target *big.Int
}

// create a new proof of work
func NewProofOfWork(b *Block) *pow {
	fmt.Printf("In NewProofOfWork\n")
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &pow{b, target}

	return pow

}

// return the string representation of the proof of work
func (p *pow) String() string {
	fmt.Printf("In String of pow\n")
	return "p: block" + string(p.block.Data) + "\n" +
		"target: " + fmt.Sprintf("%064x", p.target)
}

// run the proof of work
func (p *pow) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			p.block.PrevBlockHash,
			p.block.Data,
			[]byte(strconv.FormatInt(p.block.Timestamp, 10)),
			[]byte(strconv.FormatInt(int64(nonce), 10)),
			[]byte(strconv.FormatInt(int64(p.block.TargetBits), 10)),
		},
		[]byte{},
	)

	return data
}

func (p *pow) Run() (int, []byte) {
	fmt.Printf("In Run\n")
	var (
		hashInt big.Int
		hash    [32]byte
		nonce   = 0
	)

	fmt.Printf("Mining the block containing \"%s\"\n", p.block.Data)

	for nonce < maxNonce {
		data := p.prepareData(nonce)
		hash = sha256.Sum256(data)

		if nonce%100000 == 1 {
			fmt.Printf("Nounce: %d\n ", nonce)
			fmt.Printf("\r%x\t", hash)
		}

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(p.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	fmt.Printf("Finally Nonce: %d\n", nonce)
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate the pow
func (p *pow) Validate() bool {
	fmt.Printf("In Validate\n")
	var hashInt big.Int

	data := p.prepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(p.target) == -1

	return isValid
}
