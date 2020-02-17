package infura

import (
	"bytes"
	"encoding/json"
	"ethereum-api/pkg/app/ethereum"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseUrl string
}

func NewInfuraClient(baseURL, projectId string) *Client {
	urlWithProjectId := fmt.Sprintf("%v/v3/%v", baseURL, projectId)
	return &Client{
		baseUrl: urlWithProjectId,
	}
}

func internalInfuraError(response *http.Response) error {
	return fmt.Errorf("Server error. Status code: %v, Status: %v", response.StatusCode, response.Status)
}

func getErrorFromResponse(response *http.Response) error {
	if response != nil && response.StatusCode >= 400 {
		return internalInfuraError(response)
	}

	return nil
}

func (client *Client) Execute(method string, params ...string) ([]byte, error) {
	body := ethereum.NewBody(method, params...)

	payload, _ := json.Marshal(body)
	response, err := http.Post(client.baseUrl, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	err = getErrorFromResponse(response)
	if err != nil {
		return nil, err
	}

	reponsePayload, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return reponsePayload, err
}
