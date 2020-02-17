package transactions

import (
	"encoding/json"
	"ethereum-api/pkg/app/errors"
	"ethereum-api/pkg/app/ethereum"
	"fmt"
	"strconv"

	"github.com/INFURA/go-ethlibs/eth"
)

type Client struct {
	EthereumClient ethereum.Client
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

type TransactionResponse struct {
	ID      string       `json:"id"`
	JSONRPC string       `json:"jsonrpc"`
	Result  *Transaction `json:"result,omitempty"`
}

func parseBlockValue(blockStringValue string) (string, error) {
	if blockStringValue == "latest" || blockStringValue == "earliest" || blockStringValue == "pending" {
		return blockStringValue, nil
	}

	blockValue, err := strconv.ParseInt(blockStringValue, 10, 64)
	if err != nil {
		return "", err
	}
	block := eth.QuantityFromInt64(int64(blockValue))

	return block.String(), nil
}

func (client *Client) GetTransactionByBlockNumberAndIndex(blockValue, indexValue string) (*Transaction, error) {
	block, err := parseBlockValue(blockValue)
	if err != nil {
		parsingError := fmt.Errorf("Error parsing block number value: %v: %w", err, errors.BadInputErr)
		return nil, parsingError
	}
	index, err := strconv.ParseInt(indexValue, 10, 64)
	if err != nil {
		parsingError := fmt.Errorf("Error parsing index number value: %v: %w", err, errors.BadInputErr)
		return nil, parsingError
	}

	transactionIndex := eth.QuantityFromInt64(index)

	payload, err := client.EthereumClient.Execute("eth_getTransactionByBlockNumberAndIndex", block, transactionIndex.String())
	if err != nil {
		retrievalError := fmt.Errorf("Error retrieving transation: %v", err)
		return nil, retrievalError
	}

	response := TransactionResponse{}
	err = json.Unmarshal(payload, &response)
	if err != nil {
		parsingError := fmt.Errorf("Error parsing response: %v", err)
		return nil, parsingError
	}

	return response.Result, nil
}
