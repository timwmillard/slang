package urbandictionary

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	BaseURL = "https://mashape-community-urban-dictionary.p.rapidapi.com"
)

type Definition struct {
	Definition  string    `json:"definition"`
	Permalink   string    `json:"permalink"`
	ThumbsUp    int       `json:"thumbs_up"`
	SoundUrls   []string  `json:"sound_urls"`
	Author      string    `json:"author"`
	Word        string    `json:"word"`
	Defid       int       `json:"defid"`
	CurrentVote string    `json:"current_vote"`
	WrittenOn   time.Time `json:"written_on"`
	Example     string    `json:"example"`
	ThumbsDown  int       `json:"thumbs_down"`
}

type Client struct {
	BaseURL    *url.URL
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	base, err := url.Parse(BaseURL)
	if err != nil {
		panic(err)
	}
	return &Client{
		BaseURL:    base,
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) Define(word string) ([]Definition, error) {
	endpoint := *c.BaseURL
	endpoint.Path = path.Join(endpoint.Path, "define")
	q := endpoint.Query()
	q.Set("term", word)
	endpoint.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Rapidapi-Host", c.BaseURL.Host)
	req.Header.Add("X-Rapidapi-Key", c.APIKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		errMsg := struct {
			Message string `json:"message"`
		}{}
		err = json.NewDecoder(res.Body).Decode(&errMsg)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errMsg.Message)
	}

	result := struct {
		List []Definition `json:"list"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.List, nil
}
