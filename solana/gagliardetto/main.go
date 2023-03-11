package main

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"net"
	"net/http"
	"time"
	"web3Demo/solana/portto/httpProxy"
)

func NewHTTPTransport(
	timeout time.Duration,
	maxIdleConnsPerHost int,
	keepAlive time.Duration,
) *http.Transport {
	return &http.Transport{
		IdleConnTimeout:     timeout,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		Proxy:               http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepAlive,
		}).Dial,
	}
}

func NewHTTP(
	timeout time.Duration,
) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: httpProxy.LoggingRoundTripper{Proxied: http.DefaultTransport},
	}
}

// NewRPC creates a new Solana JSON RPC client.
func NewRPC(rpcEndpoint string) *rpc.Client {
	var (
		defaultTimeout = 25 * time.Second
	)
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient: NewHTTP(
			defaultTimeout,
		),
	}
	rpcClient := jsonrpc.NewClientWithOpts(rpcEndpoint, opts)
	return rpc.NewWithCustomRPCClient(rpcClient)
}

func main() {
	client := NewRPC("https://solana-devnet.g.alchemy.com/v2/On35d8LdFc1QGYD-wCporecGj359qian")
	resp, _ := client.GetVersion(context.TODO())
	fmt.Println("Solana-Core version: " + resp.SolanaCore)
}