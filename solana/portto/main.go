package main

import (
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"web3Demo/solana/httpProxy"
	"web3Demo/solana/portto/transfer"
)

var (
	// own accountAddress account
	accountAddress         = "AnayTW335MabjhtXTJeBit5jdLhNeUVBVPXeRKCid79D"
	assTokenAccountAddress = "Gd8nxWzbnJ2zwtn5TukvEMKKjjbFhdtqA1L67DgnRvXc"
	tokenMintAddress       = "Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr"
	ownEndpoint            = "https://solana-devnet.g.alchemy.com/v2/On35d8LdFc1QGYD-wCporecGj359qian"
	// custom connection would be:
	cli = client.New(rpc.WithEndpoint(ownEndpoint), rpc.WithHTTPClient(httpProxy.GetHttpClient()))
)

// main are following stuff from https://portto.github.io/solana-go-sdk/tour/token-transfer.html
func main() {

	transfer.TryTransferToken(cli, "", assTokenAccountAddress, tokenMintAddress)

}
