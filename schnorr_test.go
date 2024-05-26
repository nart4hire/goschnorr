package schnorr_test

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/nart4hire/goschnorr"
)

func TestSchnorr(t *testing.T) {
	s, err := schnorr.NewSchnorr(rand.Reader, sha256.New())
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	privkey, pubkey, err := s.GenKeyPair()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	message := "Hello, World!"

	sig, hash, err := s.Sign(privkey, message)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if !s.Verify(pubkey, sig, hash, message) {
		t.Error("Error: Signature is invalid")
	}
}

func TestClassicSchnorr(t *testing.T) {
	s, err := schnorr.NewClassicSchnorr(rand.Reader, sha256.New())
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	privkey, pubkey, err := s.GenKeyPair()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	message := "Hello, World!"

	sig, hash, err := s.Sign(privkey, message)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if !s.Verify(pubkey, sig, hash, message) {
		t.Error("Error: Signature is invalid")
	}
	
}

func TestGetParams(t *testing.T) {
	s, err := schnorr.NewSchnorr(rand.Reader, sha256.New())
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	p, q, g := s.GetParams()

	if p == nil || q == nil || g == nil {
		t.Error("Error: Could not get params")
	}
}