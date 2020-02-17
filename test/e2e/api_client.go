// +build integration

package e2e

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type APIClient struct {
	baseURL string
}

type Transaction struct {
	BlockHash   string `json:"blockHash"`
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	Gas         string `json:"gas"`
	GasPrice    string `json:"gasPrice"`
	Hash        string `json:"hash"`
	Input       string `json:"input"`
	Nonce       string `json:"nonce"`
	To          string `json:"to"`
	Index       string `json:"transactionIndex"`
	Value       string `json:"value"`
}

func NewAPIClient() (*APIClient, error) {
	baseURL := os.Getenv("SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &APIClient{baseURL: baseURL}, nil
}

func (client *APIClient) Healthcheck() (int, error) {
	url := fmt.Sprintf("%v%v", client.baseURL, "/healthcheck")
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	return response.StatusCode, nil
}

func (client *APIClient) GetTransaction(blockNumber, index string, headers map[string]string) (*Transaction, int, error) {
	url := fmt.Sprintf("%v/transactions/blockNumber/%v/index/%v", client.baseURL, blockNumber, index)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, 0, err
	}

	if response != nil && response.StatusCode != 200 {
		return nil, response.StatusCode, nil
	}

	var transaction Transaction
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}

	err = json.Unmarshal(body, &transaction)
	if err != nil {
		return nil, 0, err
	}

	return &transaction, response.StatusCode, nil
}
