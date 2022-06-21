package svm

import (
	"github.com/ethereum/go-ethereum/crypto"
	"hash"
	"math/big"
)

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func GetSelectorFromName(funcName string) *big.Int {
	kec := crypto.Keccak256([]byte(funcName))

	maskedKec := maskBits(250, 8, kec)

	return new(big.Int).SetBytes(maskedKec)
}

// mask excess bits
func maskBits(mask, wordSize int, slice []byte) (ret []byte) {
	excess := len(slice)*wordSize - mask
	for _, by := range slice {
		if excess > 0 {
			if excess > wordSize {
				excess = excess - wordSize
				continue
			}
			by <<= excess
			by >>= excess
			excess = 0
		}
		ret = append(ret, by)
	}
	return ret
}
