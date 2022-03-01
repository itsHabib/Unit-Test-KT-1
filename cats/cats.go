package cats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiURL = "https://api.thecatapi.com/v1/images/search?breed_ids"
)

type Client struct {
	c      HttpDo
	apiKey string
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey is required")
	}

	return &Client{
		c:      new(http.Client),
		apiKey: apiKey,
	}, nil
}

func (c *Client) GetBreeds(breedId string) ([]Metadata, error) {
	// get cat image
	getBreeds, err := http.NewRequest(http.MethodGet, apiURL+"="+breedId, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create cat image request: %w", err)
	}
	getBreeds.Header.Set("x-api-key", c.apiKey)

	// grab cat metadata
	breedsResp, err := c.c.Do(getBreeds)
	if err != nil {
		return nil, fmt.Errorf("unable to get cat metadata: %v", err)
	}
	defer breedsResp.Body.Close()
	if breedsResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 status code: %d", breedsResp.StatusCode)
	}

	var breeds []Metadata
	if err := json.NewDecoder(breedsResp.Body).Decode(&breeds); err != nil {
		return nil, fmt.Errorf("unable to unmarshal cat metadata: %w", err)
	}

	return breeds, nil
}

func (c *Client) GetCatImage(url string) ([]byte, error) {
	// download image from API
	getImage, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create cat image request: %v", err)
	}
	getImage.Header.Set("x-api-key", c.apiKey)

	imageResp, err := c.c.Do(getImage)
	if err != nil {
		return nil, fmt.Errorf("unable to get cat image: %v", err)
	}
	defer imageResp.Body.Close()
	if imageResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 status code: %d", imageResp.StatusCode)
	}

	return ioutil.ReadAll(imageResp.Body)
}

type Metadata struct {
	Url string `json:"url"`
}

type HttpDo interface {
	Do(*http.Request) (*http.Response, error)
}
