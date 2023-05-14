package algonfts

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Meningtov/algonea-backend/domain"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

func NewClient(apiKey, creatorAddress, assetNamePrefix string) domain.AlgoNftsClient {
	return &client{
		apiKey:          apiKey,
		httpClient:      retryablehttp.NewClient().StandardClient(),
		creatorAddress:  creatorAddress,
		assetNamePrefix: assetNamePrefix,
	}
}

type client struct {
	apiKey          string
	httpClient      *http.Client
	creatorAddress  string
	assetNamePrefix string
}

func (c *client) Generate(number int) (domain.ArtGenAsset, error) {
	req, err := http.NewRequest(http.MethodGet, "https://algonfts.art/v1/generate1/"+c.creatorAddress, nil)
	if err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}
	req.Header.Set("NFTGEN-KEY", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return domain.ArtGenAsset{}, errors.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var resBody responseBody
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}

	if resBody.Error != "" {
		return domain.ArtGenAsset{}, errors.Errorf("unexpected response error: %s", resBody.Error)
	}

	return domain.ArtGenAsset{
		Number:   number,
		ImageURL: resBody.Imageurl,
		Metadata: domain.Metadata{
			Name:          domain.FormatAssetName(c.assetNamePrefix, number),
			Standard:      "arc3",
			ImageMimeType: "image/webp",
			Properties:    resBody.Arc69.Properties,
			Image:         resBody.Arc69.MediaUrl,
		},
	}, nil
}

func (c *client) Update(number int, traits []domain.Traits) (domain.ArtGenAsset, error) {
	requestURL, err := url.Parse("https://algonfts.art/v1/generate1/" + c.creatorAddress)
	if err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}

	queryBuilder := requestURL.Query()
	for _, trait := range traits {
		queryBuilder.Add(trait.Layer, trait.Trait)
	}
	requestURL.RawQuery = queryBuilder.Encode()

	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}
	req.Header.Set("NFTGEN-KEY", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return domain.ArtGenAsset{}, errors.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var resBody responseBody
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}

	if resBody.Error != "" {
		return domain.ArtGenAsset{}, errors.Errorf("unexpected response error: %s", resBody.Error)
	}

	return domain.ArtGenAsset{
		Number:   number,
		ImageURL: resBody.Imageurl,
		Metadata: domain.Metadata{
			Name:          domain.FormatAssetName(c.assetNamePrefix, number),
			Standard:      "arc3",
			ImageMimeType: "image/webp",
			Properties:    resBody.Arc69.Properties,
			Image:         resBody.Arc69.MediaUrl,
		},
	}, nil
}

func (c *client) SwapTraits(number int, properties map[string]string) (domain.ArtGenAsset, error) {
	requestURL, err := url.Parse("https://algonfts.art/v1/generate1/" + c.creatorAddress)
	if err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}

	queryBuilder := requestURL.Query()
	for layer, trait := range properties {
		queryBuilder.Add(layer, trait)
	}
	requestURL.RawQuery = queryBuilder.Encode()

	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}
	req.Header.Set("NFTGEN-KEY", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return domain.ArtGenAsset{}, errors.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var resBody responseBody
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return domain.ArtGenAsset{}, errors.WithStack(err)
	}

	if resBody.Error != "" {
		return domain.ArtGenAsset{}, errors.Errorf("unexpected response error: %s", resBody.Error)
	}

	return domain.ArtGenAsset{
		Number:   number,
		ImageURL: resBody.Imageurl,
		Metadata: domain.Metadata{
			Name:          domain.FormatAssetName(c.assetNamePrefix, number),
			Standard:      "arc3",
			ImageMimeType: "image/webp",
			Properties:    resBody.Arc69.Properties,
			Image:         resBody.Arc69.MediaUrl,
		},
	}, nil
}

type responseBody struct {
	Filename string `json:"filename"`
	Cid      string `json:"cid"`
	Imageurl string `json:"imageurl"`
	Arc69    struct {
		Standard   string            `json:"standard"`
		MimeType   string            `json:"mime_type"`
		MediaUrl   string            `json:"media_url"`
		Properties map[string]string `json:"properties"`
		Arc19      struct {
			Url     string `json:"url"`
			Reserve string `json:"reserve"`
		} `json:"arc19"`
	}
	Error string `json:"error"`
}
