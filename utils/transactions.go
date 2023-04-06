package utils

import (
	"crypto/ed25519"
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

func NewGroupTransactionBuilderWithSigner(signerPK ed25519.PrivateKey) *GroupTxnBuilder {
	return &GroupTxnBuilder{signerPK: &signerPK}
}

func NewGroupTransactionBuilder() *GroupTxnBuilder {
	return &GroupTxnBuilder{}
}

type GroupTxnBuilder struct {
	signerPK *ed25519.PrivateKey
	txns     []innerTxn
}

type ResponseTransaction struct {
	RequiresSigning bool   `json:"requires_signing"`
	ToSign          string `json:"to_sign"`
	Signed          string `json:"signed"`
}

type innerTxn struct {
	tnx      types.Transaction
	toSign   bool
	signerPK *ed25519.PrivateKey
}

func (b *GroupTxnBuilder) Add(txn types.Transaction, toSign bool) *GroupTxnBuilder {
	b.txns = append(b.txns, innerTxn{
		tnx:    txn,
		toSign: toSign,
	})
	return b
}

func (b *GroupTxnBuilder) AddWithSigner(signerPK ed25519.PrivateKey, txn types.Transaction, toSign bool) *GroupTxnBuilder {
	b.txns = append(b.txns, innerTxn{
		tnx:      txn,
		toSign:   toSign,
		signerPK: &signerPK,
	})
	return b
}

func (b *GroupTxnBuilder) Build() ([]ResponseTransaction, error) {
	var txns []ResponseTransaction

	var transactions []types.Transaction
	for _, txn := range b.txns {
		transactions = append(transactions, txn.tnx)
	}
	gid, err := crypto.ComputeGroupID(transactions)
	if err != nil {
		return nil, err
	}
	for i := range b.txns {
		b.txns[i].tnx.Group = gid
		var signedTxn []byte
		if b.txns[i].toSign {
			signer := b.signerPK
			if b.txns[i].signerPK != nil {
				signer = b.txns[i].signerPK
			}
			_, signedTxn, err = crypto.SignTransaction(*signer, b.txns[i].tnx)
			if err != nil {
				return nil, err
			}
		}
		txns = append(txns, ResponseTransaction{
			RequiresSigning: !b.txns[i].toSign,
			ToSign:          base64.StdEncoding.EncodeToString(msgpack.Encode(b.txns[i].tnx)),
			Signed:          base64.StdEncoding.EncodeToString(signedTxn),
		})
	}
	return txns, nil
}
