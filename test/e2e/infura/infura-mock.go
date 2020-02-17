package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

const validInfuraResponse = `{"jsonrpc":"2.0","id":"1","result":{"blockHash":"0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35","blockNumber":"0x5bad55","from":"0xc837f51a0efa33f8eca03570e3d01a4b2cf97ffd","gas":"0x15f90","gasPrice":"0x14b8d03a00","hash":"0x311be6a9b58748717ac0f70eb801d29973661aaf1365960d159e4ec4f4aa2d7f","input":"0x","nonce":"0x4241","r":"0xe9ef2f6fcff76e45fac6c2e8080094370082cfb47e8fde0709312f9aa3ec06ad","s":"0x421ebc4ebe187c173f13b1479986dcbff5c4997c0dfeb1fd149a982ad4bcdfe7","to":"0xf49bd0367d830850456d2259da366a054038dc46","transactionIndex":"0x1","v":"0x25","value":"0x1bafa9ee16e78000"}}`

type InfureRequest struct {
	Method string `json:"method"`
}

//nolint:errcheck
func main() {
	router := chi.NewRouter()

	router.Post("/v3/{projectId}", func(w http.ResponseWriter, r *http.Request) {
		request := InfureRequest{}
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		if request.Method == "eth_getTransactionByBlockNumberAndIndex" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(validInfuraResponse))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Request method not supported: %v", request.Method)))
	})

	err := http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatal(err)
	}
}
