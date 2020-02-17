package transactions

import (
	"ethereum-api/pkg/app/reporter"
	"fmt"
)

func NewService(client Client, reporter *reporter.Reporter) *Service {
	return &Service{
		transactionsClient: client,
		reporter:           reporter,
	}
}

type Service struct {
	transactionsClient Client
	reporter           *reporter.Reporter
}

func (service *Service) GetTransactionByBlockNumberAndIndex(blockNumber, index string) (*Transaction, error) {
	transaction, err := service.transactionsClient.GetTransactionByBlockNumberAndIndex(blockNumber, index)

	if err != nil {
		service.reporter.Error("transactions.get.blocknumberandindex.error", err, map[string]interface{}{"block": blockNumber, "index": index})
		return nil, err
	}

	if transaction == nil {
		service.reporter.Info("transactions.get.blocknumberandindex.notfound", fmt.Sprintf("No transaction found in block %v at index %v", blockNumber, index), map[string]interface{}{"block": blockNumber, "index": index})
		return nil, nil
	}

	service.reporter.Info("transactions.get.blocknumberandindex", fmt.Sprintf("Found transaction in block %v at index %v", blockNumber, index), map[string]interface{}{"block": blockNumber, "index": index})
	return transaction, nil
}
