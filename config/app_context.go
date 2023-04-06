package config

import (
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

// AppContext contains all the dependencies that are used in the application.
// If you need other dependencies like an ipfs client, add it here.
type AppContext struct {
	AlgodClient    *algod.Client
	IndexerClient  *indexer.Client
	CreatorAccount crypto.Account
}

func GetAppContext() (AppContext, error) {
	cfg := loadConfig()

	algodClient, err := algod.MakeClient(cfg.AlgodBaseURL, "")
	if err != nil {
		return AppContext{}, err
	}

	indexerClient, err := indexer.MakeClient(cfg.IndexerBaseURL, "")
	if err != nil {
		return AppContext{}, err
	}

	privateKey, err := mnemonic.ToPrivateKey(cfg.CreatorMnemonic)
	if err != nil {
		return AppContext{}, err
	}
	account, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		return AppContext{}, err
	}

	return AppContext{
		AlgodClient:    algodClient,
		IndexerClient:  indexerClient,
		CreatorAccount: account,
	}, nil
}
