package utils

import (
	"errors"
	"io"
	"math/big"
)

type Generator struct {
	gen io.Reader
}

type PrimeGen interface {
	GeneratePrimeBytes(bytes int) (*big.Int, error)
	GenerateFromPrimeFactor(bytes int, q *big.Int) (*big.Int, error)
	GeneratePQ() (*big.Int, *big.Int, error)
}

func NewPrimeGen(rand io.Reader) PrimeGen {
	return &Generator{gen: rand}
}
// Generates arbitrarily long prime numbers, multiples of 8 (bytes)

func (g *Generator) GeneratePrimeBytes(bytes int) (*big.Int, error) {
	b := make([]byte, bytes)
	prime := new(big.Int)

	for {
		if _, err := io.ReadFull(g.gen, b); err != nil {
			return nil, err
		}

		// Force big number
		b[0] |= 1 << 7

		// Force odd number
		b[len(b) - 1] |= 1

		prime.SetBytes(b)

		if prime.ProbablyPrime(20) {
			return prime, nil
		}
	}
}

func (g *Generator) GenerateFromPrimeFactor(bytes int, q *big.Int) (*big.Int, error) {
	b := make([]byte, bytes)
	prime := new(big.Int)

	for {
		if _, err := io.ReadFull(g.gen, b); err != nil {
			return nil, err
		}

		// Force big number
		b[0] |= 1 << 7

		prime.SetBytes(b)
		twoq := new(big.Int).Mul(q, big.NewInt(2))
		mod := new(big.Int).Mod(prime, twoq)
		prime.Sub(prime, mod)
		prime.Add(prime, big.NewInt(1))

		if prime.ProbablyPrime(20) {
			return prime, nil
		}
	}
}

func (g *Generator) GeneratePQ() (*big.Int, *big.Int, error) {
	q, err := g.GeneratePrimeBytes(32)
	if err != nil {
		return nil, nil, err
	}

	p, err := g.GenerateFromPrimeFactor(256, q)
	if err != nil {
		return nil, nil, err
	}

	
	mod := new(big.Int)
	if mod.Mod(p, q).Cmp(big.NewInt(1)) != 0 {
		return nil, nil, errors.New("p != 1 mod q")
	}

	return p, q, nil
}