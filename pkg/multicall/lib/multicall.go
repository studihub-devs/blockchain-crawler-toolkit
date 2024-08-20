package multicall

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"new-token/pkg/multicall/gen"
	"new-token/pkg/server"
	"new-token/service/token/constant"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Eth client Variable
var EthereumClient *ethclient.Client
var BinanceSmartChainClient *ethclient.Client
var FantomClient *ethclient.Client
var CeloClient *ethclient.Client
var PolygonClient *ethclient.Client
var AvalancheCChainClient *ethclient.Client
var OptimismClient *ethclient.Client
var ArbitrumClient *ethclient.Client
var MoonbeamClient *ethclient.Client
var KavaClient *ethclient.Client
var CronosClient *ethclient.Client

var ListURLEthereum string
var ListURLBinance string
var ListURLFantom string
var ListURLCelo string
var ListURLPolygon string
var ListURLAvalancheCChain string
var ListURLOptimism string
var ListURLArbitrum string
var ListURLMoonbeam string
var ListURLKava string
var ListURLCronos string

//todo add new client

func init() {

	ListURLEthereum = server.Config.GetString("ETHEREUM_NODE_URL")
	ListURLBinance = server.Config.GetString("BINANCE_SMART_CHAIN_NODE_URL")
	ListURLFantom = server.Config.GetString("FANTOM_NODE_URL")
	ListURLCelo = server.Config.GetString("CELO_NODE_URL")
	ListURLPolygon = server.Config.GetString("POLYGON_NODE_URL")
	ListURLAvalancheCChain := server.Config.GetString("AVALANCHE_C_CHAIN_NODE_URL")
	ListURLOptimism = server.Config.GetString("OPTIMISM_NODE_URL")
	ListURLArbitrum = server.Config.GetString("ARBITRUM_NODE_URL")
	ListURLMoonbeam = server.Config.GetString("MOONBEAM_NODE_URL")
	ListURLKava = server.Config.GetString("KAVA_NODE_URL")
	ListURLCronos = server.Config.GetString("CRONOS_NODE_URL")

	EthereumClient = SwapClient(ListURLEthereum, nil)
	BinanceSmartChainClient = SwapClient(ListURLBinance, nil)
	FantomClient = SwapClient(ListURLFantom, nil)
	CeloClient = SwapClient(ListURLCelo, nil)
	PolygonClient = SwapClient(ListURLPolygon, nil)
	AvalancheCChainClient = SwapClient(ListURLAvalancheCChain, nil)
	OptimismClient = SwapClient(ListURLOptimism, nil)
	ArbitrumClient = SwapClient(ListURLArbitrum, nil)
	MoonbeamClient = SwapClient(ListURLMoonbeam, nil)
	KavaClient = SwapClient(ListURLKava, nil)
	CronosClient = SwapClient(ListURLCronos, nil)

	//todo add new node url
}

// func connectClient(listUrlRPC string) *ethclient.Client {

// 	listAPIOrWS := strings.Fields(listUrlRPC)
// 	var client *ethclient.Client
// 	for i, link := range listAPIOrWS {
// 		clientSub, err := ethclient.Dial(link)
// 		if err != nil {
// 			log.Println(log.LogLevelError, fmt.Sprintf("NewClientBSC: change to URL index %v : ", i+1), err.Error())
// 		} else {
// 			client = clientSub
// 			break
// 		}
// 	}
// 	return client

// 	client, err := ethclient.Dial(clientUrl)
// 	if err != nil {
// 		log.Fatal("err ethclient.Dial: ", err)
// 	}
// 	return client
// }

type Call struct {
	contract common.Address
	method   string
	args     []any
	abi      *abi.ABI
}

func NewCall(contract common.Address, method string, args []any, abi *abi.ABI) *Call {
	return &Call{
		contract: contract,
		method:   method,
		args:     args,
		abi:      abi,
	}
}

var (
	uint256Type    abi.Type
	bytesSliceType abi.Type
	arguments      abi.Arguments
	parsedAbi      abi.ABI
)

func init() {
	var err error
	uint256Type, err = abi.NewType("uint256", "", nil)
	if err != nil {
		log.Fatalf("addressSliceType failed:%v", err)
	}
	bytesSliceType, err = abi.NewType("bytes[]", "", nil)
	if err != nil {
		log.Fatalf("bytesSliceType failed:%v", err)
	}
	arguments = abi.Arguments{
		{Type: uint256Type, Name: "Height"},
		{Type: bytesSliceType, Name: "ReturnDatas"},
	}

	parsedAbi, err = abi.JSON(strings.NewReader(gen.MultiCallABI))
	if err != nil {
		log.Fatalf("abi.JSON failed:%v", err)
	}
}

var Retry = 3
var BackoffInterval = time.Second * 2

func Do(chainId string, calls []Call) (height uint64, results []any, err error) {
	var client *ethclient.Client
	switch chainId {
	case constant.ETHEREUM_ID:
		client = SwapClient(ListURLEthereum, EthereumClient)
	case constant.BINANCE_SMART_CHAIN_ID:
		client = SwapClient(ListURLBinance, BinanceSmartChainClient)
	case constant.FANTOM_ID:
		client = SwapClient(ListURLFantom, FantomClient)
	case constant.CELO_ID:
		client = SwapClient(ListURLCelo, CeloClient)
	case constant.POLYGON_ID:
		client = SwapClient(ListURLPolygon, PolygonClient)
	case constant.AVALANCHE_C_CHAIN_ID:
		client = SwapClient(ListURLAvalancheCChain, AvalancheCChainClient)
	case constant.OPTIMISM_ID:
		client = SwapClient(ListURLOptimism, OptimismClient)
	case constant.ARBITRUM_ID:
		client = SwapClient(ListURLArbitrum, ArbitrumClient)
	case constant.MOONBEAM_ID:
		client = SwapClient(ListURLMoonbeam, MoonbeamClient)
	case constant.KAVA_ID:
		client = SwapClient(ListURLKava, KavaClient)
	case constant.CRONOS_ID:
		client = SwapClient(ListURLCronos, CronosClient)
		//todo add new client
	default:
		client = EthereumClient
	}

	results = make([]interface{}, len(calls))

	if len(calls) != len(results) {
		err = fmt.Errorf("#calls != #results")
		return
	}

	var (
		targets []common.Address
		inputs  [][]byte
	)
	for _, call := range calls {
		callABI := call.abi
		method, exist := callABI.Methods[call.method]
		if !exist {
			err = fmt.Errorf("method '%s' not found", call.method)
			return
		}

		var arguments []byte
		arguments, err = method.Inputs.Pack(call.args...)
		if err != nil {
			return
		}

		targets = append(targets, call.contract)
		inputs = append(inputs, append(method.ID, arguments...))
	}

	packed, err := parsedAbi.Pack("", targets, inputs)
	if err != nil {
		return
	}

	var resultBytes []byte
	for i := 0; i < Retry; i++ {
		resultBytes, err = client.CallContract(
			context.Background(), ethereum.CallMsg{
				Data: append(common.FromHex(gen.MultiCallBin), packed...),
			},
			nil,
		)
		if err != nil {
			time.Sleep(BackoffInterval)
			continue
		}
		break
	}
	if err != nil {
		return
	}

	var output struct {
		Height      *big.Int
		ReturnDatas [][]byte
	}
	resultInterface, err := arguments.Unpack(resultBytes)
	if err != nil {
		return
	}
	err = arguments.Copy(&output, resultInterface)
	if err != nil {
		return
	}

	if len(output.ReturnDatas) != len(calls) {
		err = fmt.Errorf("#ReturnDatas != #calls")
		return
	}
	for i, returnData := range output.ReturnDatas {
		if len(returnData) == 0 {
			continue
		}

		call := calls[i]
		callABI := call.abi

		method := callABI.Methods[call.method]
		var returnValue []any
		returnValue, err = method.Outputs.Unpack(returnData)
		if err != nil {
			return
		}

		err = method.Outputs.Copy(&results[i], returnValue)
		if err != nil {
			results[i] = returnValue
		}
	}
	height = output.Height.Uint64()

	return height, results, nil
}

func SwapClient(listNodeRPC string, currentClient *ethclient.Client) *ethclient.Client {
	if currentClient != nil {
		// call test currentClient if no err -> return current client
		_, err := currentClient.BlockNumber(context.Background())
		if err == nil {
			return currentClient
		}
	}

	listAPIOrWS := strings.Fields(listNodeRPC)
	client := currentClient
	for _, link := range listAPIOrWS {
		clientSub, err := ethclient.Dial(link)
		client = clientSub
		if err != nil {
			continue
		} else {
			//call test rpc
			_, err := clientSub.BlockNumber(context.Background())
			if err == nil {
				return clientSub
			}
		}
	}
	return client
}
