package algorand

import (
	"context"
	"errors"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"strings"
)

func GetCreatedAssets(ctx context.Context, indexerClient *indexer.Client, accountID string, unitNamePrefixFilter ...string) ([]models.Asset, error) {
	unitNameFilter := ""
	if len(unitNamePrefixFilter) > 0 {
		unitNameFilter = unitNamePrefixFilter[0]
	}
	nextToken := ""
	assets := make([]models.Asset, 0)
	for {
		accountAssetsRes, err := indexerClient.LookupAccountCreatedAssets(accountID).
			Next(nextToken).
			Limit(1000).
			Do(ctx)
		if err != nil {
			if strings.HasPrefix(err.Error(), "HTTP 404") || strings.HasPrefix(err.Error(), "HTTP 400") {
				return nil, errors.New("account not found")
			}
			return nil, err
		}

		for _, a := range accountAssetsRes.Assets {
			if a.Deleted {
				continue
			}
			if unitNameFilter != "" && !strings.HasPrefix(a.Params.UnitName, unitNameFilter) {
				continue
			}
			assets = append(assets, a)
		}

		if accountAssetsRes.NextToken == "" {
			break
		}
		nextToken = accountAssetsRes.NextToken
	}

	return assets, nil
}

func GetOwnedAssets(ctx context.Context, indexerClient *indexer.Client, accountID string) ([]uint64, error) {
	nextToken := ""
	ids := make([]uint64, 0)
	for {
		accountAssetsRes, err := indexerClient.LookupAccountAssets(accountID).Next(nextToken).Do(ctx)
		if err != nil {
			if strings.HasPrefix(err.Error(), "HTTP 404") || strings.HasPrefix(err.Error(), "HTTP 400") {
				return nil, errors.New("account not found")
			}
			return nil, err
		}

		for _, a := range accountAssetsRes.Assets {
			if a.Amount < 1 || a.Deleted || a.IsFrozen {
				continue
			}
			ids = append(ids, a.AssetId)
		}

		if accountAssetsRes.NextToken == "" {
			break
		}
		nextToken = accountAssetsRes.NextToken
	}

	return ids, nil
}
