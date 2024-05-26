package utils_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/nart4hire/goschnorr/utils"
)

func TestPrimeGen(t *testing.T) {
	primegen := utils.NewPrimeGen(rand.Reader)
	prime, err := primegen.GeneratePrimeBytes(2048 / 8)

	if err != nil {
		t.Errorf("Error: %s", err)
	
	}

	if prime.BitLen() != 2048 {
		t.Errorf("Error: %s", err)
	}

}

func TestPrimeGenFromPrimeFactor(t *testing.T) {
	primegen := utils.NewPrimeGen(rand.Reader)
	q, err := primegen.GeneratePrimeBytes(32)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	p, err := primegen.GenerateFromPrimeFactor(256, q)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	mod := new(big.Int)
	if mod.Mod(p, q).Cmp(big.NewInt(1)) != 0 {
		t.Error("Error: p != 1 mod q")
	}
}

func TestPrimeGenPQ(t *testing.T) {
	primegen := utils.NewPrimeGen(rand.Reader)
	p, q, err := primegen.GeneratePQ()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	mod := new(big.Int)
	if mod.Mod(p, q).Cmp(big.NewInt(1)) != 0 {
		t.Error("Error: p != 1 mod q")
	}
}

func TestAlphaGen(t *testing.T) {
	p, q, err := utils.NewPrimeGen(rand.Reader).GeneratePQ()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	a, err := utils.AlphaGen(p, q)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if a.Exp(a, q, p).Cmp(big.NewInt(1)) != 0 {
		t.Error("Error: Alpha is not a generator of the group")
	}
}

func TestSchnorrGen(t *testing.T) {
	p, q, err := utils.NewPrimeGen(rand.Reader).GeneratePQ()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	g, err := utils.SchnorrGen(p, q)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if g.Mod(g, p).Cmp(big.NewInt(1)) == 0 {
		t.Error("Error: g is 1")
	}
}