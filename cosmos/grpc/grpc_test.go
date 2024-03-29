package grpc

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	gov_v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmos_proto "github.com/cosmos/gogoproto/proto"
	"github.com/golang/protobuf/proto"
	"log"
	"own.cosmos.demo"
	"testing"
)

func Test_getNodeInfo(t *testing.T) {
	defer cosmos.Conn.Close()

	client := tmservice.NewServiceClient(cosmos.Conn)
	request := tmservice.GetNodeInfoRequest{}
	info, _ := client.GetNodeInfo(context.Background(), &request)
	log.Printf("Node info: Network: %v, ListenAddr: %v", info.DefaultNodeInfo.Network, info.DefaultNodeInfo.ListenAddr)
	log.Printf("Application info: goVer: %v, appName: %v, cosmosSdkVer: %v",
		info.ApplicationVersion.GoVersion, info.ApplicationVersion.AppName, info.ApplicationVersion.CosmosSdkVersion)
}

// Test_getLatestBlock
func Test_getLatestBlock(t *testing.T) {
	defer cosmos.Conn.Close()

	tmClient := tmservice.NewServiceClient(cosmos.Conn)
	request := tmservice.GetLatestBlockRequest{}
	res, _ := tmClient.GetLatestBlock(context.Background(), &request)
	log.Println("Latest block height:", res.Block.Header.Height)
}

// Test_getDenomMeta is for retrieving denom's meta data
func Test_getDenomMeta(t *testing.T) {
	defer cosmos.Conn.Close()

	queryClient := bank.NewQueryClient(cosmos.Conn)
	request := bank.QueryDenomMetadataRequest{Denom: "usdt"}
	metadata, err := queryClient.DenomMetadata(context.Background(), &request)
	if err != nil {
		log.Fatalf("Retrieve error %v", err)
	}

	fmt.Println(proto.MarshalTextString(metadata))
}

// Test_decodeBlock
func Test_decodeBlock(t *testing.T) {
	defer cosmos.Conn.Close()

	//height := 8592209 // msg transfer
	//height := 8708499 // delegate deposit
	height := 8389774 // delegate reward

	tmClient := tmservice.NewServiceClient(cosmos.Conn)
	request := tmservice.GetBlockByHeightRequest{Height: int64(height)}
	res, err := tmClient.GetBlockByHeight(context.Background(), &request)
	if err != nil {
		log.Fatalf("Failed to get block by height: %v", err)
	}

	blockHash := sha256.Sum256(res.BlockId.Hash)
	fmt.Printf("Block hash: %x\n", blockHash)

	// traverse & decode transaction
	for _, txBytes := range res.Block.Data.Txs {
		var txObj tx.Tx
		err := cosmos_proto.Unmarshal(txBytes, &txObj)
		if err != nil {
			log.Fatalf("Failed to unmarshal transaction: %v", err)
		}

		// construct txHash
		hash := sha256.Sum256(txBytes)
		fmt.Printf("Transaction hash: %x\n\n", hash)

		// decode messages
		for _, msg := range txObj.Body.Messages {

			switch msg.TypeUrl {

			// for transfer msg
			case "/cosmos.bank.v1beta1.MsgSend":
				var bankMsg bank.MsgSend
				if err := cosmos_proto.Unmarshal(msg.Value, &bankMsg); err != nil {
					log.Fatalf("Failed to unmarshal bankMsg")
				}
				fmt.Printf("From: %s\n", bankMsg.FromAddress)
				fmt.Printf("To: %s\n", bankMsg.ToAddress)
				for _, coin := range bankMsg.Amount {
					fmt.Printf("Denom: %s, Amount: %s\n", coin.Denom, coin.Amount)
				}

			// for governance
			case "/cosmos.gov.v1.msgDeposit":
				var govMsg gov_v1.MsgDeposit
				if err := cosmos_proto.Unmarshal(msg.Value, &govMsg); err != nil {
					t.Errorf("Failed to unmarshal msgDeposit")
				}
				fmt.Printf("Governance proposalId:%d, address:%s, deposit:%s\n",
					govMsg.ProposalId, govMsg.Depositor, govMsg.Amount)

			// for delegation deposit
			case "/cosmos.staking.v1beta1.MsgDelegate":
				var stakingMsg staking.MsgDelegate
				if err := cosmos_proto.Unmarshal(msg.Value, &stakingMsg); err != nil {
					t.Errorf("Failed to unmarshal msgDelegate")
				}
				fmt.Printf("Staking from:%s,\n to validator:%s,\n with amount:%s\n",
					stakingMsg.DelegatorAddress, stakingMsg.ValidatorAddress, stakingMsg.Amount)

			case "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":
				var distributionMsg distribution.MsgWithdrawDelegatorReward
				if err := cosmos_proto.Unmarshal(msg.Value, &distributionMsg); err != nil {
					t.Errorf("Failed to unmarshar MsgWithdrawDelegatorReward")
				}
				fmt.Printf("Distribution delegator:%s, validaton:%s\n",
					distributionMsg.DelegatorAddress, distributionMsg.ValidatorAddress)

			default:
				fmt.Printf("Skip for uncared msg type: %s\n", msg.TypeUrl)
			}

		}

	}
}
