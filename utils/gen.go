package utils

import (
	"math/big"
	"io"
)

type Generator struct {
	gen io.Reader
}

type PrimeGen interface {
	GeneratePrime(bytes int) (*big.Int, error)
}

func NewPrimeGen(rand io.Reader) PrimeGen {
	return &Generator{gen: rand}
}
// Generates arbitrarily long prime numbers, multiples of 8 (bytes)

func (g *Generator) GeneratePrime(bytes int) (*big.Int, error) {
	b := make([]byte, bytes)
	prime := new(big.Int)

	for {
		if n, err := io.ReadFull(g.gen, b); n != bytes || err != nil {
			return nil, err
		}

		// Force odd number
		b[len(b) - 1] |= 1

		prime.SetBytes(b)

		if prime.ProbablyPrime(20) {
			return prime, nil
		}
	}
}