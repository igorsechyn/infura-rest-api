package ethereum

type Client interface {
	Execute(method string, params ...string) ([]byte, error)
}
