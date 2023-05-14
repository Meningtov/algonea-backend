package config

import (
	"os"
)

type envVars struct {
	AlgodBaseURL    string
	IndexerBaseURL  string
	CreatorMnemonic string
	RekeyerMnemonic string
	NftStorageAPIKey string
	AlgoNftsAPIKey string
	UnitNamePrefix string
}

func loadConfig() envVars {
	envConfig := envVars{}
	envConfig.AlgodBaseURL = getOrDefault("ALGOD_URL", "https://testnet-api.algonode.cloud")
	envConfig.IndexerBaseURL = getOrDefault("INDEXER_URL", "https://testnet-idx.algonode.cloud")

	envConfig.CreatorMnemonic = os.Getenv("CREATOR_MNEMONIC")
	envConfig.RekeyerMnemonic = os.Getenv("REKEYER_MNEMONIC")

	envConfig.NftStorageAPIKey = os.Getenv("NFTSTORAGE_API_KEY")
	envConfig.AlgoNftsAPIKey = os.Getenv("ALGONFTS_API_KEY")
	envConfig.UnitNamePrefix = "ALG"

	return envConfig
}

func getOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
