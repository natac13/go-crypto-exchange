package main

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type User struct {
	ID         int64
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func NewUser(privateKey string, userId int64) (*User, error) {
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	return &User{
		ID:         userId,
		PrivateKey: pk,
		PublicKey:  publicKeyECDSA,
	}, nil
}
