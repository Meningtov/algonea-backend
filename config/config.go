package config

import (
	"os"
)

type envVars struct {
	AlgodBaseURL    string
	IndexerBaseURL  string
	CreatorMnemonic string
}

func loadConfig() envVars {
	envConfig := envVars{}
	envConfig.AlgodBaseURL = getOrDefault("ALGOD_URL", "https://testnet-api.algonode.cloud")
	envConfig.IndexerBaseURL = getOrDefault("INDEXER_URL", "https://testnet-idx.algonode.cloud")

	// FIXME added this account for demo purposes, replace with env var
	//   JA2RM7HHCRYM6UMSGYJSHEY3GODY4QSLLPPVVLTJAWHI2QKJSUQMPU6234
	envConfig.CreatorMnemonic = "fever wear obvious scissors galaxy you laundry fix public universe soft debris crystal rather illness announce rose point sentence glove wall random oil above load"

	return envConfig
}

func getOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
