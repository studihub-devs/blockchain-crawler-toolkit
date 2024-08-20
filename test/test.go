package main

import (
	"context"
	"fmt"
	"math/big"
	"new-token/pkg/log"
	"new-token/pkg/server"
	"new-token/service/token/crawler"
	"new-token/service/token/dao"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

// func main() {

// 	// totalSupply, ok := new(big.Int).SetString("1887", 10)

// 	// totalSupplyFloat := new(big.Float).SetInt(totalSupply)

// 	// fmt.Println("totalSupply", ok, totalSupply)

// 	// dec, _ := new(big.Float).SetString("10")

// 	// fmt.Println("chia", totalSupplyFloat.Quo(totalSupplyFloat, dec))
// 	// for i := 0; i < 5; i++ {
// 	// 	totalSupplyFloat = totalSupplyFloat.Quo(totalSupplyFloat, dec)
// 	// 	fmt.Println("chia", totalSupplyFloat)

// 	// }

// 	// fmt.Println("chia", totalSupplyFloat)

// 	// return

// 	// html.SplitHTML("0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270", "137")

// 	// crypto1 := dao.Crypto{
// 	// 	Address: "0x2170ed0880ac9a755fd29b2688956bd959f933f8",
// 	// 	ChainId: "56",
// 	// }

// 	// crypto2 := dao.Crypto{
// 	// 	Address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
// 	// 	ChainId: "1",
// 	// }

// 	// listCrypto := dao.ListCrypto{}

// 	// listCrypto.Cryptos = append(listCrypto.Cryptos, crypto1)
// 	// listCrypto.Cryptos = append(listCrypto.Cryptos, crypto2)

// 	// listNewPairToken := []string{"0xc001bbe2b87079294c63ece98bdd0a88d761434e"}
// 	// dex := dao.Dex{
// 	// 	ChainId: "56",
// 	// }
// 	// listCrypto, _, _, err := contract.CrawlInfoNewToken(dex, listNewPairToken)
// 	// if err != nil {
// 	// 	log.Println(log.LogLevelError, "CrawlerScheduleTokenDex contract.CrawlInfoNewToken(dex, listNewPairToken)", err.Error())
// 	// }

// 	// //name, symbol, totalSupply, Decimals

// 	// for _, crypto := range listCrypto.Cryptos {
// 	// 	cryptoCallApiScan := api.CallApiScanIsProxyContractVerified(crypto)

// 	// 	crypto.Proxy = cryptoCallApiScan.Proxy
// 	// 	crypto.ContractVerified = cryptoCallApiScan.ContractVerified

// 	// 	cryptoSplitHtmlScan, err := html.SplitHTML(crypto.Address, crypto.ChainId)
// 	// 	if err == nil {
// 	// 		crypto.Logo = cryptoSplitHtmlScan.Logo
// 	// 		crypto.Holders = cryptoSplitHtmlScan.Holders
// 	// 		crypto.Website = cryptoSplitHtmlScan.Website
// 	// 		crypto.Socials = cryptoSplitHtmlScan.Socials
// 	// 		crypto.Explorer = cryptoSplitHtmlScan.Explorer
// 	// 	}

// 	// 	if crypto.ChainId == "56" || crypto.ChainId == "1" {
// 	// 		//todo:Call Honeypot
// 	// 		proofScamHoneypot, err := api.CallHoneypot(crypto.Address, crypto.ChainId)
// 	// 		if err != nil {
// 	// 			log.Println(log.LogLevelError, "CrawlerScheduleTokenDex CallHoneypot "+crypto.ChainId+" "+crypto.Address, err.Error())
// 	// 		}

// 	// 		// fmt.Println("proofScamHoneypot", proofScamHoneypot)
// 	// 		if proofScamHoneypot != "" {
// 	// 			crypto.IsScam = true
// 	// 			err = json.Unmarshal([]byte(proofScamHoneypot), &crypto.Proof)
// 	// 			if err != nil {
// 	// 				log.Println(log.LogLevelError, "CrawlerScheduleTokenDex json.Unmarshal "+crypto.ChainId+" "+crypto.Address, err.Error())
// 	// 			}
// 	// 		}

// 	// 		// fmt.Println("proof", *crypto.Proof.IsScam)

// 	// 		if crypto.ChainId == "56" {
// 	// 			//todo:Call Staysafu
// 	// 			cryptoStaySafu, err := api.CallStaySafu(crypto)
// 	// 			if err != nil {
// 	// 				log.Println(log.LogLevelError, "CrawlerScheduleTokenDex  api.CallStaySafu(crypto) "+crypto.ChainId+" "+crypto.Address, err.Error())
// 	// 			} else {
// 	// 				crypto.Proof = cryptoStaySafu.Proof
// 	// 			}
// 	// 		}
// 	// 		if crypto.Proof.IsScam != nil {
// 	// 			fmt.Println("scam", *crypto.Proof.IsScam)
// 	// 			crypto.IsScam = true
// 	// 		}
// 	// 		if crypto.Proof.IsWarning != nil {
// 	// 			fmt.Println("warning", *crypto.Proof.IsWarning)
// 	// 			crypto.IsWarning = true
// 	// 		}
// 	// 	}

// 	// 	fmt.Println("address", crypto.Address, crypto.ChainId, crypto.Name, crypto.Symbol, crypto.Decimals, crypto.TotalSupply)
// 	// 	fmt.Println("holder", crypto.Holders, crypto.Website, crypto.Socials, crypto.Logo, crypto.Explorer)
// 	// 	fmt.Println("verified", crypto.ContractVerified, crypto.Proxy)
// 	// 	fmt.Println("scam", crypto.IsScam, crypto.IsWarning, crypto.Proof)

// 	// 	err = crypto.Insert()
// 	// 	if err != nil {
// 	// 		log.Println(log.LogLevelError, "CrawlerScheduleTokenDex crypto.Insert() "+crypto.ChainId+" "+crypto.Address, err.Error())
// 	// 	}

// 	// }

// 	// dex := dao.Dex{
// 	// 	Name:           "Pancakeswap bsc",
// 	// 	ChainId:        "56",
// 	// 	FactoryAddress: "0xca143ce32fe78f1f7019d7d551a6402fc5350c73",
// 	// 	PairLength:     1221428,
// 	// }

// 	// listNewPairAddress, _, err := contract.CrawlNewPairAddress(dex)
// 	// if err != nil {
// 	// 	log.Println(log.LogLevelError, "CrawlerScheduleTokenDex contract.CrawlNewPairAddress(dex)", err.Error())
// 	// }

// 	// fmt.Println("listNewPairAddress", listNewPairAddress)

// 	// listNewPairToken, _, err := contract.CrawlNewPairToken(dex, listNewPairAddress)
// 	// if err != nil {
// 	// 	log.Println(log.LogLevelError, "CrawlerScheduleTokenDex contract.CrawlNewPairToken(dex, listNewPairAddress)", err.Error())
// 	// }

// 	// fmt.Println("listNewPairToken", listNewPairToken)

// 	// listCrypto, _, err := contract.CrawlInfoNewToken(dex, listNewPairToken)
// 	// if err != nil {
// 	// 	log.Println(log.LogLevelError, "CrawlerScheduleTokenDex contract.CrawlInfoNewToken(dex, listNewPairToken)", err.Error())
// 	// }
// 	// //name, symbol, totalSupply, Decimals

// 	// fmt.Println("listCrypto", len(listCrypto.Cryptos))

// 	// for _, crypto := range listCrypto.Cryptos {
// 	// 	cryptoCallApiScan := api.CallApiScanIsProxyContractVerified(crypto)

// 	// 	crypto.Proxy = cryptoCallApiScan.Proxy
// 	// 	crypto.ContractVerified = cryptoCallApiScan.ContractVerified

// 	// 	cryptoSplitHtmlScan, err := html.SplitHTML(crypto.Address, crypto.ChainId)
// 	// 	if err == nil {
// 	// 		crypto.Logo = cryptoSplitHtmlScan.Logo
// 	// 		crypto.Holders = cryptoSplitHtmlScan.Holders
// 	// 		crypto.Website = cryptoSplitHtmlScan.Website
// 	// 		crypto.Socials = cryptoSplitHtmlScan.Socials
// 	// 	}

// 	// 	if crypto.ChainId == "56" || crypto.ChainId == "1" {
// 	// 		//todo:Call Honeypot
// 	// 		proofScamHoneypot, err := api.CallHoneypot(crypto.Address, crypto.ChainId)
// 	// 		if err != nil {
// 	// 			log.Println(log.LogLevelError, "CrawlerScheduleTokenDex CallHoneypot "+crypto.ChainId+" "+crypto.Address, err.Error())
// 	// 		}
// 	// 		if proofScamHoneypot != "" {
// 	// 			crypto.IsScam = true
// 	// 		}

// 	// 		if crypto.ChainId == "56" {
// 	// 			//todo:Call Staysafu

// 	// 			api.CallStaySafu(crypto)
// 	// 		}
// 	// 	}
// 	// 	fmt.Println("->", crypto.Address, crypto.ChainId, crypto.Name, crypto.Symbol, crypto.Decimals, crypto.TotalSupply, crypto.Holders, crypto.Website, crypto.Socials, crypto.Logo, crypto.Explorer, crypto.ContractVerified, crypto.Proxy, crypto.IsScam, crypto.IsWarning, crypto.Proof)
// 	// }

// 	// crypto, err := html.SplitHTML("0x31d8288a4849106d3d276de42496102ba54f843e", "1")

// 	// fmt.Println("err", err)

// 	// fmt.Println("-", crypto.Website, "webite")
// 	// fmt.Println("-", crypto.Holders, "Holders")
// 	// fmt.Println("-", crypto.Logo, "Logo")
// 	// fmt.Println("-", crypto.Socials, "Socials")

// 	dex := dao.Dex{
// 		ChainId:        "56",
// 		FactoryAddress: "0xca143ce32fe78f1f7019d7d551a6402fc5350c73",
// 	}
// 	a, err := contract.CallPairlength(dex)
// 	fmt.Println(a, err)
// }

var client *ethclient.Client
var startBlock uint64 = 15424185

func init() {
	fmt.Println("run")
	client = NewClient()
}

func NewClient() *ethclient.Client {

	nodeWssURL := server.Config.GetString("ETHEREUM_NODE_URL")
	client, err := ethclient.Dial(nodeWssURL)
	if err != nil {
		log.Println(log.LogLevelError, "NewClient: ethclient.Dial(constant.URLINFURA) ", err.Error())
	}
	return client
}

func CrawlFutureLogs() {
	latestBLock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Println(log.LogLevelError, "Client.BlockNumber(context.Background())", err)
	}
	if latestBLock <= startBlock {
		fmt.Println("No new block")
		return
	}

	fmt.Println("Latest block ", latestBLock)
	latestBlockBigInt := new(big.Int).SetUint64(latestBLock)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startBlock)), //block create factory uniswap v3 12369621
		ToBlock:   latestBlockBigInt,
		Addresses: []common.Address{
			common.HexToAddress("0x1f98431c8ad98523631ae4a59f267346ea31f984"),
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Println(log.LogLevelError, "client.FilterLogs(context.Background(), query)", err.Error())
	}

	fmt.Println("len log", len(logs))

	for _, log := range logs {
		if strings.ToLower(log.Topics[0].String()) == "0x783cca1c0412dd0d695e784568c96da2e9c22ff989357a2e8b1d9b2b4e6b7118" {

			fmt.Println(common.BytesToAddress(log.Topics[1].Bytes())) // token0
			fmt.Println(common.BytesToAddress(log.Topics[2].Bytes())) // token1
			fmt.Println("data")
			fmt.Println(common.BytesToAddress(log.Data)) // pai address
		}
	}

	startBlock = latestBLock
	fmt.Println("Latest filter block: ", latestBLock)
}

func main() {
	// listDex := dao.ListDex{}

	// listChain := []string{"1"}
	// err := listDex.GetListDex(listChain)
	// fmt.Println(err)

	// for _, ele := range listDex.Dexs {
	// 	fmt.Println(ele.Name, ele.ChainId)
	// }

	crypto := dao.Crypto{
		Id:          uuid.MustParse("f000aa01-0451-4000-b000-000000000000"),
		Name:        "testhoi",
		Symbol:      "TEST",
		Address:     "0xe00000000000000000000000000000000",
		AddressShow: "0xe00000000000000000000000000000000",
		ChainName:   "ethereum",
		CryptoId:    "gear5__token_ethereum_0xe00000000000000000000000000000000",
		Createddate: "123",
		Updateddate: "2023-02-17T18:16:39Z",
	}

	err := crawler.PostCryptoToElasticsearch(crypto)
	fmt.Println("err", err)

	// crypto, err := html.SplitHTML("0xe1146b9ac456fcbb60644c36fd3f868a9072fc6e", "250")
	// fmt.Println(err)
	// fmt.Println(crypto)

	// go func() {
	// 	count := 0
	// 	for {
	// 		crypto := dao.Crypto{
	// 			Name: fmt.Sprintf("%d", count),
	// 		}

	// 		byteValue, err := json.Marshal(crypto)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		err = ms.Redis.Store.RPush("crypto", byteValue).Err()
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		count += 1
	// 		time.Sleep(1 * time.Second)
	// 		if count == 10 {
	// 			break
	// 		}
	// 	}
	// }()

	// go func() {
	// 	time.Sleep(15 * time.Second)
	// 	crypto := dao.Crypto{
	// 		Name: fmt.Sprintf("%d", 100),
	// 	}

	// 	byteValue, err := json.Marshal(crypto)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	err = ms.Redis.Store.RPush("crypto", byteValue).Err()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()

	// // idEnd := ""
	// // for idEnd != "END" {
	// for {
	// 	// time.Sleep(5 * time.Second)
	// 	// fmt.Println(">", ms.Redis.Store.LLen("crypto"), "<")

	// 	if ms.Redis.Store.LLen("crypto").Val() != 0 {
	// 		fmt.Println(">", ms.Redis.Store.LLen("crypto"), "<")
	// 		result, err := ms.Redis.Store.BLPop(1*time.Second, "crypto").Result()
	// 		// fmt.Println(result)
	// 		if err != nil {
	// 			fmt.Println("err", err)
	// 		}
	// 		// idEnd = result[1]
	// 		crypto1 := dao.Crypto{}
	// 		if result != nil {
	// 			err = json.Unmarshal([]byte(result[1]), &crypto1)
	// 			if err != nil {
	// 				fmt.Println(err)
	// 			}
	// 			fmt.Println(crypto1.Name)
	// 		}
	// 	}
	// }
}
