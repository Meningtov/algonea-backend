package domain

import (
	"fmt"
)

type AlgoNftsClient interface {
	Generate(number int) (ArtGenAsset, error)
	Update(number int, traits []Traits) (ArtGenAsset, error)
	SwapTraits(number int, properties map[string]string) (ArtGenAsset, error)
}

type Traits struct {
	Layer string `json:"layer"`
	Trait string `json:"trait"`
}

type Asset struct {
	Asset int `json:"asset"`
}

type ArtGenAsset struct {
	Number      int      `json:"number"`
	ImageURL    string   `json:"image_url"`
	MetadataCID string   `json:"metadata_cid"`
	Metadata    Metadata `json:"metadata"`
}

type Metadata struct {
	Name          string            `json:"name"`
	Standard      string            `json:"standard"`
	Image         string            `json:"image"`
	ImageMimeType string            `json:"image_mime_type"`
	Properties    map[string]string `json:"properties"`
}

type InputAsset struct {
	AssetID        uint64
	Reserve        string
	Name           string
	UnitName       string
	MetadataBase64 string
}

type MintMetadata struct {
	Number   int    `json:"number"`
	AssetID  uint64 `json:"asset_id"`
	UnitName string `json:"unit_name"`
	Name     string `json:"name"`
}

func FormatAssetName(namePrefix string, number int) string {
	displayNumber := fmt.Sprintf("%s", fmt.Sprintf("%04d", number))
	return fmt.Sprintf("%s #%s", namePrefix, displayNumber)
}

func FormatAssetUnitName(unitNamePrefix string, number int) string {
	displayNumber := fmt.Sprintf("%s", fmt.Sprintf("%04d", number))
	return fmt.Sprintf("%s%s", unitNamePrefix, displayNumber)
}
