package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Meningtov/algonea-backend/config"
	"github.com/Meningtov/algonea-backend/domain"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conf, err := config.GetAppContext()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to load config")
	}
	log.Info().Msgf("Loaded config. Creator: %s, algod: %s", conf.CreatorAccount, conf.AlgodURL)
	log.Info().Msgf("Assets to swap : Asset #%d, Asset #%d", conf.AssetsToSwap[0].Asset, conf.AssetsToSwap[1].Asset)
	log.Info().Msgf("Layer to swap : %s", conf.LayerToSwap)

	loadedAssets, err := loadAssets(conf)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to load assets")
	}

	if err := generateArt(conf, loadedAssets); err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to generate art")
	}

	if err := updateReserveAddress(conf); err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to mint nft")
	}

	log.Info().Msg("Done")
}

func loadAssets(conf config.Config) ([]map[string]string, error) {
	var LayerToSwap = conf.LayerToSwap
	var assetsProperties []map[string]string
	var mapPropertiesToSwap []map[string]string

	for i := 0; i < len(conf.AssetsToSwap); i++ {
		assetsMeta, err := loadMeta(conf.AssetsToSwap[i].Asset)
		if err != nil {
			panic(err)
		}

		loadedProperties, PropertiesToSwap, err := loadProperties(assetsMeta, LayerToSwap)
		if err != nil {
			panic(err)
		}

		assetsProperties = append(assetsProperties, loadedProperties)
		mapPropertiesToSwap = append(mapPropertiesToSwap, PropertiesToSwap)
	}
	swappedProperties, err := swapLayer(assetsProperties, mapPropertiesToSwap)
	if err != nil {
		panic(err)
	}

	return swappedProperties, nil
}

func loadMeta(asset int) ([]domain.ArtGenAsset, error) {
	var meta []domain.ArtGenAsset
	file, err := os.ReadFile(fmt.Sprintf("assets/art_gen/%d.json", asset))

	if err != nil {
		panic(err)
	}
	var m domain.ArtGenAsset
	if err = json.Unmarshal(file, &m); err != nil {
		panic(err)
	}
	meta = append(meta, m)
	log.Info().Msgf("Loaded meta of asset #%d", asset)

	return meta, nil
}

func loadProperties(meta []domain.ArtGenAsset, layerToSwap string) (map[string]string, map[string]string, error) {
	p := map[string]string{}
	t := map[string]string{}

	for _, m := range meta {
		for traitLayer, traitName := range m.Metadata.Properties {
			if _, ok := p[traitLayer]; !ok {
				p[traitLayer] = traitName
			}
			if traitLayer == layerToSwap {
				t[traitLayer] = traitName
			}
		}
	}
	return p, t, nil
}

func swapLayer(properties []map[string]string, PropertiesToSwap []map[string]string) ([]map[string]string, error) {
	for traitLayer, _ := range properties[0] {
		for tL, tN := range PropertiesToSwap[1] {
			if traitLayer == tL {
				properties[0][traitLayer] = tN
			}
		}
	}

	for traitLayer, _ := range properties[1] {
		for tL, tN := range PropertiesToSwap[0] {
			if traitLayer == tL {
				properties[1][traitLayer] = tN
			}
		}
	}
	log.Info().Msgf("Swapped assets layer to swap")
	return properties, nil
}

func generateArt(conf config.Config, properties []map[string]string) error {
	for i := 0; i < len(conf.AssetsToSwap); i++ {
		artGenAsset, err := conf.AlgoNftsClient.SwapTraits(conf.AssetsToSwap[i].Asset, properties[i])
		if err != nil {
			return err
		}
		if err := writeAsset(conf, artGenAsset); err != nil {
			return err
		}
		log.Info().Msgf("Generated artwork for asset#%d", conf.AssetsToSwap[i].Asset)
	}
	return nil
}

func writeAsset(conf config.Config, asset domain.ArtGenAsset) error {
	response, err := http.Get(asset.ImageURL)
	if err != nil {
		return errors.WithStack(err)
	}
	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	// save image locally
	file, err := os.Create(fmt.Sprintf("assets/art_gen/%d.webp", asset.Number))
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(resBody))
	if err != nil {
		return errors.WithStack(err)
	}

	// store image on ipfs
	ipfsCID, err := conf.NftStorageClient.UploadFile(resBody, "image/webp")
	if err != nil {
		return errors.WithStack(err)
	}
	asset.Metadata.Image = "ipfs://" + ipfsCID

	// store metadata on ipfs
	metadataJson, err := json.Marshal(asset.Metadata)
	if err != nil {
		return errors.WithStack(err)
	}
	metadataCID, err := conf.NftStorageClient.UploadFile(metadataJson, "application/json")
	if err != nil {
		return errors.WithStack(err)
	}
	asset.MetadataCID = metadataCID

	// save metadata file locally
	data, err := json.MarshalIndent(asset, "", " ")
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.WriteFile(fmt.Sprintf("assets/art_gen/%d.json", asset.Number), data, 0666)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Switching to the new files
func updateReserveAddress(conf config.Config) error {
	for i := 0; i < len(conf.AssetsToSwap); i++ {
		log.Info().Msgf("Setting reserve. creator=%s asset=#%d", conf.CreatorAccount, conf.AssetsToSwap[i].Asset)
		artGenMeta, mintMeta, err := readMetadata(conf.AssetsToSwap[i].Asset)
		if err != nil {
			return errors.WithStack(err)
		}

		params, err := conf.AlgodClient.SuggestedParams().Do(context.TODO())
		if err != nil {
			return errors.WithStack(err)
		}

		reserveAddress, err := getReserveAddressFromCID(artGenMeta.MetadataCID)
		if err != nil {
			return errors.WithStack(err)
		}

		transaction, err := future.MakeAssetConfigTxn(conf.CreatorAccount, nil, params, mintMeta.AssetID, conf.CreatorAccount, reserveAddress, "", "", false)
		if err != nil {
			return errors.WithStack(err)
		}

		txID, signedTxn, err := crypto.SignTransaction(conf.RekeyPrivateKey, transaction)
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = conf.AlgodClient.SendRawTransaction(signedTxn).Do(context.Background())
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = future.WaitForConfirmation(conf.AlgodClient, txID, 4, context.Background())
		if err != nil {
			return errors.WithStack(err)
		}

		log.Info().Msgf("Successfully updated asset #%d reserve address", conf.AssetsToSwap[i].Asset)
	}
	return nil

}

func readMetadata(number int) (domain.ArtGenAsset, domain.MintMetadata, error) {
	artGenMetaJson, err := os.ReadFile(fmt.Sprintf("assets/art_gen/%d.json", number))
	if err != nil {
		return domain.ArtGenAsset{}, domain.MintMetadata{}, errors.WithStack(err)
	}
	var artGenAsset domain.ArtGenAsset
	err = json.Unmarshal(artGenMetaJson, &artGenAsset)
	if err != nil {
		return domain.ArtGenAsset{}, domain.MintMetadata{}, errors.WithStack(err)
	}

	mintMetaJson, err := os.ReadFile(fmt.Sprintf("assets/minted/%d.json", number))
	if err != nil {
		return domain.ArtGenAsset{}, domain.MintMetadata{}, errors.WithStack(err)
	}
	var mintMeta domain.MintMetadata
	err = json.Unmarshal(mintMetaJson, &mintMeta)
	if err != nil {
		return domain.ArtGenAsset{}, domain.MintMetadata{}, errors.WithStack(err)
	}

	return artGenAsset, mintMeta, nil
}

func getReserveAddressFromCID(cidToUpdate string) (string, error) {
	decodedCID, err := cid.Decode(cidToUpdate)
	if err != nil {
		return "", errors.WithStack(err)
	}
	reserve, err := reserveAddressFromCID(decodedCID)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return reserve, nil
}

func reserveAddressFromCID(cidToEncode cid.Cid) (string, error) {
	decodedMultiHash, err := multihash.Decode(cidToEncode.Hash())
	if err != nil {
		return "", errors.WithStack(err)
	}
	return types.EncodeAddress(decodedMultiHash.Digest)
}
