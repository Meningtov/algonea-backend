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

	if err := generateArt(conf); err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to generate art")
	}

	if err := updateReserveAddress(conf); err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to mint nft")
	}

	log.Info().Msg("Done")
}

func generateArt(conf config.Config) error {
	artGenAsset, err := conf.AlgoNftsClient.Update(conf.AssetNumber, conf.NewTraits)
	if err != nil {
		return err
	}
	if err := writeAsset(conf, artGenAsset); err != nil {
		return err
	}
	log.Info().Msgf("Generated artwork for asset#%d", conf.AssetNumber)
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
	log.Info().Msgf("Setting reserve. creator=%s asset=#%d", conf.CreatorAccount, conf.AssetNumber)
	artGenMeta, mintMeta, err := readMetadata(conf.AssetNumber)

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

	log.Info().Msg("Successfully updated reserve address")
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
