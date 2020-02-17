package ethereum

func NewBody(method string, params ...string) map[string]interface{} {
	body := map[string]interface{}{
		"id":      "1",
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}

	return body
}
