package contract

import (
	"fmt"
	"math/big"
	"new-token/pkg/log"
	multicall "new-token/pkg/multicall/lib"
	"new-token/service/token/dao"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

/*
Call smartcontract get token info (name, symbol, decimal, totalSupply)
*/
func CrawlInfoNewToken(dex dao.Dex, listNewToken []string) (dao.ListCrypto, int, int, error) {

	tokenFail := 0

	listNewCrypto := dao.ListCrypto{}

	crawlFail := dao.CrawlFail{
		ChainId: dex.ChainId,
		Defi:    dex.Name,
		Type:    "CrawlInfoNewToken",
	}

	tokenInfosBigCalls := [][]multicall.Call{}
	tokenInfosCalls := []multicall.Call{}
	for index, address := range listNewToken {
		for _, method := range _tokenInfoMethodList {
			call := multicall.NewCall(
				common.HexToAddress(address),
				method,
				[]interface{}{},
				tokenInfoABI,
			)
			tokenInfosCalls = append(tokenInfosCalls, *call)
		}
		if len(tokenInfosCalls) == _tokenInfoCallPerMulticalls || index == int(len(listNewToken)-1) {
			tokenInfosBigCalls = append(tokenInfosBigCalls, tokenInfosCalls)
			tokenInfosCalls = []multicall.Call{}
		}
	}

	for index, calls := range tokenInfosBigCalls {
		// fmt.Println("index", index)
		_, result, _ := multicall.Do(dex.ChainId, calls)

		for i := 0; i < len(result); {
			// fmt.Println("i", i)
			crypto := &dao.Crypto{
				ChainId: dex.ChainId,
				Address: listNewToken[index*_tokenInfoCallPerMulticalls/4+i/4],
			}
			if result[i] != nil && result[i+1] != nil && result[i+2] != nil && result[i+3] != nil {
				crypto.Name = result[i].(string)
				crypto.Symbol = strings.ToUpper(result[i+1].(string))
				crypto.Decimal = result[i+2].(uint8)
				crypto.TotalSupplyBeforeDivideDecimal = result[i+3].(*big.Int)

				totalSupply := crypto.TotalSupplyBeforeDivideDecimal

				totalSupplyFloat := new(big.Float).SetInt(totalSupply)

				dec, ok := new(big.Float).SetString("10")
				if ok {
					for i := 0; i < int(crypto.Decimal); i++ {
						totalSupplyFloat = totalSupplyFloat.Quo(totalSupplyFloat, dec)
					}
				}

				crypto.TotalSupply = fmt.Sprintf("%f", totalSupplyFloat)

			} else {
				res, err := CallTokenInfo(crypto.Address, dex.ChainId)
				if err != nil {
					if err.Error() == "429" {
						time.Sleep(30 * time.Second)
						r, err := CallTokenInfo(crypto.Address, dex.ChainId)
						if err != nil {
							crawlFail.Address = crypto.Address
							crawlFail.Errors = err.Error()
							err := crawlFail.CallInsertFail()
							if err != nil {
								log.Println(log.LogLevelError, "CrawlInfoNewToken crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
							}
							tokenFail += 1
							// infoCrawl.ListTokenFail = append(infoCrawl.ListTokenFail, crypto.Address)
							i += len(_tokenInfoMethodList)
							continue
						} else {
							if r[0] != nil && r[1] != nil && r[2] != nil && r[3] != nil {
								crypto.Name = r[0].(string)
								crypto.Symbol = strings.ToUpper(r[1].(string))
								crypto.Decimal = r[2].(uint8)
								crypto.TotalSupplyBeforeDivideDecimal = r[3].(*big.Int)

								totalSupply := crypto.TotalSupplyBeforeDivideDecimal

								totalSupplyFloat := new(big.Float).SetInt(totalSupply)

								dec, ok := new(big.Float).SetString("10")
								if ok {
									for i := 0; i < int(crypto.Decimal); i++ {
										totalSupplyFloat = totalSupplyFloat.Quo(totalSupplyFloat, dec)
									}
								}

								crypto.TotalSupply = fmt.Sprintf("%f", totalSupplyFloat)
							}
						}
					} else {
						crawlFail.Address = crypto.Address
						crawlFail.Errors = err.Error()
						err := crawlFail.CallInsertFail()
						if err != nil {
							log.Println(log.LogLevelError, "CrawlInfoNewToken crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
						}
						tokenFail += 1
						// infoCrawl.ListTokenFail = append(infoCrawl.ListTokenFail, crypto.Address)
						i += len(_tokenInfoMethodList)
						continue
					}
				} else {
					if res[0] != nil && res[1] != nil && res[2] != nil && res[3] != nil {
						crypto.Name = res[0].(string)
						crypto.Symbol = strings.ToUpper(res[1].(string))
						crypto.Decimal = res[2].(uint8)
						crypto.TotalSupplyBeforeDivideDecimal = res[3].(*big.Int)

						totalSupply := crypto.TotalSupplyBeforeDivideDecimal

						totalSupplyFloat := new(big.Float).SetInt(totalSupply)

						dec, ok := new(big.Float).SetString("10")
						if ok {
							for i := 0; i < int(crypto.Decimal); i++ {
								totalSupplyFloat = totalSupplyFloat.Quo(totalSupplyFloat, dec)
							}
						}

						crypto.TotalSupply = fmt.Sprintf("%f", totalSupplyFloat)

					}
				}
			}

			listNewCrypto.Cryptos = append(listNewCrypto.Cryptos, *crypto)
			i += len(_tokenInfoMethodList)
		}
	}

	return listNewCrypto, len(listNewCrypto.Cryptos), tokenFail, nil
}

func CallTokenInfo(address string, chainId string) ([]any, error) {
	calls := []multicall.Call{}

	for _, method := range _tokenInfoMethodList {
		call := multicall.NewCall(
			common.HexToAddress(address),
			method,
			[]interface{}{},
			tokenInfoABI,
		)
		calls = append(calls, *call)

	}
	_, result, err := multicall.Do(chainId, calls)
	if err != nil {
		return result, err
	}

	return result, nil
}

// func GetChainName(chainId string) string {
// 	chainName := ""
// 	switch chainId {
// 	case constant.ETHEREUM_ID:
// 		chainName = constant.Ethereum_Chain_Name
// 	case constant.BINANCE_SMART_CHAIN_ID:
// 		chainName = constant.BinanceSmartChain_Chain_Name
// 	case constant.FANTOM_ID:
// 		chainName = constant.Fantom_Chain_Name
// 	case constant.CELO_ID:
// 		chainName = constant.Celo_Chain_Name
// 	case constant.POLYGON_ID:
// 		chainName = constant.Polygon_Chain_Name
// 	case constant.AVALANCHE_C_CHAIN_ID:
// 		chainName = constant.AvalancheCChain_Chain_Name
// 	case constant.OPTIMISM_ID:
// 		chainName = constant.Optimism_Chain_Name
// 	case constant.ARBITRUM_ID:
// 		chainName = constant.Arbitrum_Chain_Name
// 	case constant.MOONBEAM_ID:
// 		chainName = constant.Moonbeam_Chain_Name
// 	case constant.KAVA_ID:
// 		chainName = constant.Kava_Chain_Name
// 	case constant.CRONOS_ID:
// 		chainName = constant.Cronos_Chain_Name
// 	default:
// 		chainName = constant.Ethereum_Chain_Name
// 	}

// 	return chainName
// }
