package utils

import (
	"errors"
	"fmt"
	"math/big"
)

func AlphaGen(p, q *big.Int) (*big.Int, error) {
	limiter := big.NewInt(4_294_967_295)
	ctr := big.NewInt(1)
	one := big.NewInt(1)
	mod := new(big.Int)

	if mod.Mod(p, q).Cmp(one) != 0 {
		fmt.Println(mod.Mod(p, q))
		return nil, errors.New("p != 1 mod q")
	}

	for mod.Exp(ctr, q, p).Cmp(one) != 0 {
		ctr.Add(ctr, one)
		if ctr.Cmp(limiter) == 0 {
			return nil, errors.New("could not find alpha")
		}
	}
	return ctr, nil
}

func SchnorrGen(p, q *big.Int) (*big.Int, error) {
	one := big.NewInt(1)
	mod := new(big.Int)
	if mod.Mod(p, q).Cmp(one) != 0 {
		fmt.Println(mod.Mod(p, q))
		return nil, errors.New("p != 1 mod q")
	}

	r := new(big.Int).Sub(p, one)
	r.Div(r, q)

	h := big.NewInt(2)

	g := new(big.Int).Exp(h, r, p)
	for g.Cmp(one) == 0 {
		h.Add(h, one)
		g.Exp(h, r, p)
	}

	return g, nil
}