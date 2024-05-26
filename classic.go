package schnorr

import (
	"errors"
	"hash"
	"io"
	"math/big"

	"github.com/nart4hire/goschnorr/utils"
)

type classicschnorr struct {
	p *big.Int
	q *big.Int
	alpha *big.Int
	hash hash.Hash
	rand io.Reader
}

type Schnorr interface {
	Verify(pubkey, sig, hash []byte, message string) bool
	Sign(privkey []byte, message string) ([]byte, []byte, error)
	GenKeyPair() ([]byte, []byte, error)
}

func NewClassicSchnorrFromParam(p, q, alpha *big.Int, rand io.Reader, hash hash.Hash) Schnorr {
	return &classicschnorr{p: p, q: q, alpha: alpha, hash: hash, rand: rand}
}

func NewClassicSchnorr(rand io.Reader, hash hash.Hash) (Schnorr, error) {
	gen := utils.NewPrimeGen(rand)
	p, q, err := gen.GeneratePQ()
	if err != nil {
		return nil, err
	}

	alpha, err := utils.AlphaGen(p, q)
	if err != nil {
		return nil, err
	}

	return &classicschnorr{p: p, q: q, alpha: alpha, hash: hash, rand: rand}, nil
}

func (s *classicschnorr) Verify(pubkey, sig, hash []byte, message string) bool {
	h := new(big.Int).SetBytes(hash)
	sigNum := new(big.Int).SetBytes(sig)
	pubkeyNum := new(big.Int).SetBytes(pubkey)

	rv := new(big.Int).Exp(s.alpha, sigNum, s.p)
	rv.Mul(rv, new(big.Int).Exp(pubkeyNum, h, s.p))
	rv.Mod(rv, s.p)

	ev := append(rv.Bytes(), []byte(message)...)
	s.hash.Reset()
	s.hash.Write(ev)
	ev = s.hash.Sum([]byte{})
	evnum := new(big.Int).SetBytes(ev)

	return h.Cmp(evnum) == 0
}

func (s *classicschnorr) Sign(privkey []byte, message string) ([]byte, []byte, error) {
	kbytes := make([]byte, 32)
	_, err := io.ReadFull(s.rand, kbytes)
	if err != nil {
		return nil, nil, err
	}
	k := new(big.Int).SetBytes(kbytes)
	r := new(big.Int).Exp(s.alpha, k, s.p)

	mbytes := append(r.Bytes(), []byte(message)...)
	s.hash.Reset()
	s.hash.Write(mbytes)
	h := s.hash.Sum([]byte{}) // The hash has to be 256 bits
	hnum := new(big.Int).SetBytes(h)

	key := new(big.Int).SetBytes(privkey)
	sig := new(big.Int).Add(k, new(big.Int).Mul(key, hnum))
	return sig.Bytes(), h, nil
}

func (s *classicschnorr) GenKeyPair() ([]byte, []byte, error) {
	buffer := make([]byte, 32)
	privkey := new(big.Int).Add(s.q, big.NewInt(1))
	for privkey.Cmp(s.q) >= 0 {
		_, err := io.ReadFull(s.rand, buffer)
		if err != nil {
			return nil, nil, err
		}
		privkey.SetBytes(buffer)
	}

	pubkey := new(big.Int).Exp(s.alpha, new(big.Int).Neg(privkey), s.p)
	if pubkey == nil {
		return nil, nil, errors.New("could not generate pubkey")
	}

	return privkey.Bytes(), pubkey.Bytes(), nil
}