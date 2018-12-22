package validator

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	pvm "github.com/tendermint/tendermint/privval"
)

func GenFilePV(path string, prv crypto.PrivKey) *pvm.FilePV {
	privValidator := pvm.GenFilePV(path)
	privValidator.PrivKey = prv
	privValidator.PubKey = prv.PubKey()
	privValidator.Address = prv.PubKey().Address()
	privValidator.Save()
	return privValidator
}

func GenFilePVWithECDSA(path string, prv *ecdsa.PrivateKey) *pvm.FilePV {
	pb := gcrypto.FromECDSA(prv)
	var p secp256k1.PrivKeySecp256k1
	copy(p[:], pb)
	return GenFilePV(path, p)
}

func NewTransactorFromPV(pv *pvm.FilePV) *bind.TransactOpts {
	addr := common.BytesToAddress(pv.Address)
	return &bind.TransactOpts{
		From: addr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != addr {
				return nil, errors.New("not authorized to sign this account")
			}
			prv := pv.PrivKey.(secp256k1.PrivKeySecp256k1)
			key := bytesToECDSAPrvKey(prv[:])
			signature, err := gcrypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}

func bytesToECDSAPrvKey(b []byte) *ecdsa.PrivateKey {
	pv, err := gcrypto.ToECDSA(b)
	if err != nil {
		panic(err)
	}
	return pv
}
