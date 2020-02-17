# ethreum-api
> REST API for ethereum blockchain using INFURA

## Usage

### Configuration

ethereum-api expected following configuration values to be provided via environment variables:

- `EHTEREUM_API_INFURA_PROJECT_ID` - id of the project created in Infura
- `EHTEREUM_API_INFURA_BASE_URL` - base url of the network (e.g https://mainnet.infura.io/v3/)

### Building from sources

1. checkout git repository
2. verify all dev dependencies are installed per CONTRIBUTING guide
3. run `make verify` to run all the tests
4. run `make build` to create a server executable under `build` folder

Run

```
ETHEREUM_API_INFURA_BASE_URL=<INFURA_BASE_URL> ETHEREUM_API_INFURA_PROJECT_ID=<YOUR_PROJECT_ID> build/bin/darwin.amd64/ethereum-api
```

### Running in docker

To build a docker image with latest code run

```
make docker
```

This will create a docker image `igorsechyn/ethereum-api` tagged with the current git commit hash. Run

```
docker run -e ETHEREUM_API_INFURA_BASE_URL=<INFURA_BASE_URL> -e ETHEREUM_API_INFURA_PROJECT_ID=<YOUR_PROJECT_ID> -p 8080:8080 igorsechyn/ethereum-api:<HASH>
```

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md)

## Implementation notes
See [NOTES.md](NOTES.md)
