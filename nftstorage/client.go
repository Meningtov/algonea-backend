package nftstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	nftStorageApiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkaWQ6ZXRocjoweDA5NzMyYTdDQkFjYjlDRWJFOUJjMzgwMmNlNkYzODQ4ZjE1YWJhRDciLCJpc3MiOiJuZnQtc3RvcmFnZSIsImlhdCI6MTY1MzE2NzU4NjEwNCwibmFtZSI6Im1vc3RseWZyZW5zIn0.Udv1Uf0O324JEXvoBp1LokvMknopsVbP0sAgRLRghjY"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) UploadFile(data []byte, mediaType string) (string, error) {
	req, _ := http.NewRequest("POST", "https://api.nft.storage/upload", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+nftStorageApiKey)
	req.Header.Set("Content-Type", mediaType)
	req.Header.Set("Accept", "application/json")

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var resBody responseBody
	if err := json.NewDecoder(response.Body).Decode(&resBody); err != nil {
		return "", err
	}

	if !resBody.Ok || response.StatusCode != 200 {
		return "", fmt.Errorf("non 'Ok' response from nftstorage. code=%d error=%v", response.StatusCode, resBody.Error)
	}

	return resBody.Value.Cid, nil
}

func (c Client) PinHash(hash string) error {
	body := map[string]string{"cid": hash}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.nft.storage/pins", bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+nftStorageApiKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("received not OK status code. code=%d", resp.StatusCode)
	}
	return nil
}

type responseBody struct {
	Ok    bool `json:"ok"`
	Value struct {
		Cid string `json:"cid"`
	} `json:"value"`
	Error struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
}
