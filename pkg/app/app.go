package app

import (
	"ethereum-api/pkg/app/config"
	env "ethereum-api/pkg/app/config/envconfig"
	"ethereum-api/pkg/app/ethereum"
	"ethereum-api/pkg/app/ethereum/infura"
	"ethereum-api/pkg/app/reporter"
	"ethereum-api/pkg/app/reporter/zerolog"
	"ethereum-api/pkg/app/transactions"
)

type App struct {
	Reporter            *reporter.Reporter
	TransactionsService *transactions.Service
	Config              config.Config
}

func CreateTransactionsService(infuraClient ethereum.Client, reporter *reporter.Reporter) *transactions.Service {
	transactionsClient := transactions.Client{
		EthereumClient: infuraClient,
	}
	return transactions.NewService(transactionsClient, reporter)
}

func New() *App {
	config := env.LoadConfig()
	reporter := reporter.New(zerolog.NewJSONLogger())
	infuraClient := infura.NewInfuraClient(config.InfuraBaseURL, config.InfuraProjectId)

	transactionService := CreateTransactionsService(infuraClient, reporter)
	return &App{
		Reporter:            reporter,
		TransactionsService: transactionService,
		Config:              config,
	}
}
