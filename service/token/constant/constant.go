package constant

import (
	"new-token/pkg/server"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

var (
	RESP_SUCCESS_STATUS_CODE      = 200
	RESP_TOO_MANY_REQ_STATUS_CODE = 429
	RESP_NOT_FOUND_STATUS_CODE    = 404
	WAIT_DURATION_WHEN_RATE_LIMIT = 5 * time.Second

	ETHEREUM_ID            = "1"
	BINANCE_SMART_CHAIN_ID = "56"
	FANTOM_ID              = "250"
	CELO_ID                = "42220"
	POLYGON_ID             = "137"
	AVALANCHE_C_CHAIN_ID   = "43114"
	OPTIMISM_ID            = "10"
	ARBITRUM_ID            = "42161"
	MOONBEAM_ID            = "1284"
	KAVA_ID                = "2222"
	CRONOS_ID              = "25"

	NULL_ADDRESS = common.HexToAddress("0x0000000000000000000000000000000000000000")
	ETH_ADDRESS  = common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")

	Ethereum_Chain_Name          = "ethereum"
	BinanceSmartChain_Chain_Name = "binance"
	Fantom_Chain_Name            = "fantom"
	Celo_Chain_Name              = "celo"
	Polygon_Chain_Name           = "polygon"
	AvalancheCChain_Chain_Name   = "avalanche"
	Optimism_Chain_Name          = "optimism"
	Arbitrum_Chain_Name          = "arbitrum"
	Moonbeam_Chain_Name          = "moonbeam"
	Kava_Chain_Name              = "kava"
	Cronos_Chain_Name            = "cronos"

	Uniswap_V2_Subghaph     = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"
	Ubeswap_Subghaph        = "https://api.thegraph.com/subgraphs/name/ubeswap/ubeswap"
	TraderJoe_Subghaph      = "https://api.thegraph.com/subgraphs/name/ibrahimkecici/traderjoe"
	Velodrome_Subghaph      = "https://api.thegraph.com/subgraphs/name/dmihal/velodrome"
	Sushi_V2_Subghaph       = "https://api.thegraph.com/subgraphs/name/mariorz/sushiswap-arbitrum"
	Solarflare_Subghaph     = "https://api.thegraph.com/subgraphs/name/solarbeamio/solarflare-subgraph"
	Spiritswap_Subghaph     = "https://api.thegraph.com/subgraphs/name/chrisjess/spiritswap-fantom-exchange"
	Pancakeswap_V2_Subghaph = "https://api.thegraph.com/subgraphs/name/jieeevest/subtest2"
	Quickswap_Subghaph      = "https://api.thegraph.com/subgraphs/name/sameepsi/quickswap06"

	Ethereum_Chain_Call_Transfers    = "https://etherscan.io/token/generic-tokentxns2?m=normal&contractAddress=%s&a=&sid=567d9322b5b08f81a26d55090301c2f0&p=1"
	BinanceSmartChain_Call_Transfers = "https://bscscan.com/token/generic-tokentxns2?m=normal&contractAddress=%s&a=&sid=507d4057ff01310fc11de79bcb2be7cf&p=1"
	Fantom_Call_Transfers            = ""
	Celo_Call_Transfers              = "https://celoscan.io/token/generic-tokentxns2?m=normal&contractAddress=%s&a=&sid=c299886bb864569910846521b639af89&p=1"
	Polygon_Call_Transfers           = "https://polygonscan.com/token/generic-tokentxns2?m=normal&contractAddress=%s&a=&sid=9cff6e2986726253c82831f3465e1e29&p=1"
	AvalancheCChain_Call_Transfers   = "https://snowtrace.io/token/generic-tokentxns2?m=normal&contractAddress=%s&a=&sid=f8d44d8d3046eb24777179c3cda5ee59&p=1"
	Optimism_Call_Transfers          = ""
	Arbitrum_Call_Transfers          = "https://arbiscan.io/token/generic-tokentxns2?m=normal&contractAddress=%s&a=&sid=61dbb59a70de37f0bca8ff8b714cfb43&p=1"
	Moonbeam_Call_Transfers          = "https://moonscan.io/token/generic-tokentxns2?m=normal&contractAddress=%s&a=&sid=7e4ce09adc26b1846d55bee124b13378&p=1"
	Kava_Call_Transfers              = ""
	Cronos_Call_Transfers            = "https://cronoscan.com/token/generic-tokentxns2?m=normal&contractAddress=%s&a=&sid=b866220b09f813c1d13b72bf3be060d8&p=1"

	Etherscan     = "https://etherscan.io/token/%s"
	Bscscan       = "https://bscscan.com/token/%s"
	Ftmscan       = "https://ftmscan.com/token/%s"
	Celoscan      = "https://celoscan.io/token/%s"
	Polygonscan   = "https://polygonscan.com/token/%s"
	Avalanchescan = "https://snowtrace.io/token/%s"
	Optimismscan  = "https://optimistic.etherscan.io/token/%s"
	Arbitrumscan  = "https://arbiscan.io/token/%s"
	Moonbeamscan  = "https://moonscan.io/token/%s"
	Kavascan      = ""
	Cronosscan    = "https://cronoscan.com/token/%s"

	Ethereum_Chain_Call_ContractCode    = "https://api.etherscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	BinanceSmartChain_Call_ContractCode = "https://api.bscscan.com/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	Fantom_Call_ContractCode            = "https://api.ftmscan.com/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	Celo_Call_ContractCode              = "https://api.celoscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	Polygon_Call_ContractCode           = "https://api.polygonscan.com/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	AvalancheCChain_Call_ContractCode   = "https://api.snowtrace.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	Optimism_Call_ContractCode          = "https://api-optimistic.etherscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	Arbitrum_Call_ContractCode          = "https://api.arbiscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	Moonbeam_Call_ContractCode          = "https://api-moonbeam.moonscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s"
	Kava_Call_ContractCode              = ""
	Cronos_Call_ContractCode            = "https://api.cronoscan.com/api?module=contract&action=getsourcecode&address=%s&apikey=%s"

	TIME_SLEEP_CRAWL_NEW_PAIR_DEX = 24 * 60 * 60 * time.Second //24Hours

	API_KEY_SEARCH = "n89M3osmUowY&nQ$"
)

func GetExplorerByChainId(chainId string) (string, string, string) {
	var urlCallContractCode string
	var apikey string
	var urlExplorer string
	// chainList := []string{"1", "10", "1284", "137", "2222", "25", "250", "42161", "42220", "43114", "56"}

	if chainId == "1" { //ethereum
		urlCallContractCode = Ethereum_Chain_Call_ContractCode
		apikey = server.Config.GetString("API_ETHERSCAN")
		urlExplorer = Etherscan
	}
	if chainId == "10" { //optimism
		urlCallContractCode = Optimism_Call_ContractCode
		apikey = server.Config.GetString("API_OPTIMISMSCAN")
		urlExplorer = Optimismscan
	}
	if chainId == "1284" { //moonbeam
		urlCallContractCode = Moonbeam_Call_ContractCode
		apikey = server.Config.GetString("API_MOONBEAMSCAN")
		urlExplorer = Moonbeamscan
	}
	if chainId == "137" { //polygon
		urlCallContractCode = Polygon_Call_ContractCode
		apikey = server.Config.GetString("API_POLYGONSCAN")
		urlExplorer = Polygonscan
	}
	if chainId == "2222" { //kava
		urlCallContractCode = ""
		apikey = ""
		urlExplorer = ""
	}
	if chainId == "25" { //cronos
		urlCallContractCode = Cronos_Call_ContractCode
		apikey = server.Config.GetString("API_CRONOSSCAN")
		urlExplorer = Cronosscan
	}
	if chainId == "250" { //fantom
		urlCallContractCode = Fantom_Call_ContractCode
		apikey = server.Config.GetString("API_FTMSCAN")
		urlExplorer = Ftmscan
	}
	if chainId == "42161" { //arbitrum
		urlCallContractCode = Arbitrum_Call_ContractCode
		apikey = server.Config.GetString("API_ARBISCAN")
		urlExplorer = Arbitrumscan
	}
	if chainId == "42220" { //celo
		urlCallContractCode = Celo_Call_ContractCode
		apikey = server.Config.GetString("API_CELOSCAN")
		urlExplorer = Celoscan
	}
	if chainId == "43114" { //avalanche
		urlCallContractCode = AvalancheCChain_Call_ContractCode
		apikey = server.Config.GetString("API_AVALANCHESCAN")
		urlExplorer = Avalanchescan
	}
	if chainId == "56" {
		urlCallContractCode = BinanceSmartChain_Call_ContractCode
		apikey = server.Config.GetString("API_BSCSCAN")
		urlExplorer = Bscscan
	}

	return urlCallContractCode, apikey, urlExplorer
}

func GenSubcategoryByChainname(chainname string) string {
	chainname = strings.TrimSpace(chainname)

	subcategory := strings.ToUpper(chainname[:1]) + strings.ToLower(chainname[1:]) + " " + "Ecosystem"

	return subcategory
}
