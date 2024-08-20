package contract

import (
	"fmt"
	"math/big"
	"new-token/pkg/log"
	multicall "new-token/pkg/multicall/lib"
	"new-token/service/token/constant"
	"new-token/service/token/dao"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

/*
Input dex, list new pairAddress pool.
Get token0, token1, reserve0, reserve1, blockTimestampLast
Check exist token0, token1 -> return list new token
*/
func CrawlNewPairToken(dex dao.Dex, listPairAddress []string) ([]string, int, int, error) {
	pairTokenFail := 0

	pairTokenInfosBigCalls := [][]multicall.Call{}
	pairTokenInfosCalls := []multicall.Call{}
	for index, pairAddress := range listPairAddress {
		for _, method := range _pairTokenMethodList {
			call := multicall.NewCall(
				common.HexToAddress(pairAddress),
				method,
				[]interface{}{},
				pairABI,
			)
			pairTokenInfosCalls = append(pairTokenInfosCalls, *call)
		}
		if len(pairTokenInfosCalls) == _pairTokenCallPerMulticalls || index == int(len(listPairAddress)-1) {
			pairTokenInfosBigCalls = append(pairTokenInfosBigCalls, pairTokenInfosCalls)
			pairTokenInfosCalls = []multicall.Call{}
		}
	}

	listToken := []string{}

	for index, calls := range pairTokenInfosBigCalls {

		_, result, _ := multicall.Do(dex.ChainId, calls)

		for i := 0; i < len(result); {
			pairToken := &dao.PairToken{
				ChainId:     dex.ChainId,
				ChainName:   dao.MapChainList[dex.ChainId],
				Defi:        dex.Name,
				PairAddress: listPairAddress[index*_pairTokenCallPerMulticalls/3+i/3],
			}
			if result[i] != nil && result[i+1] != nil && result[i+2] != nil {
				pairToken.Token0 = strings.ToLower(result[i].(common.Address).String())
				pairToken.Token1 = strings.ToLower(result[i+1].(common.Address).String())

				if len(result[i+2].([]interface{})) > 2 {
					pairToken.Reserve0 = fmt.Sprintf("%d", result[i+2].([]interface{})[0].(*big.Int))
					pairToken.Reserve1 = fmt.Sprintf("%d", result[i+2].([]interface{})[1].(*big.Int))
					pairToken.BlockTimestampLast = int(result[i+2].([]interface{})[2].(uint32))
				}
			} else {
				pairTokenInfosCalls := []multicall.Call{}
				for _, method := range _pairTokenMethodList {
					call := multicall.NewCall(
						common.HexToAddress(pairToken.PairAddress),
						method,
						[]interface{}{},
						pairABI,
					)
					pairTokenInfosCalls = append(pairTokenInfosCalls, *call)
				}
				_, res, err := multicall.Do(dex.ChainId, pairTokenInfosCalls)
				if err != nil {
					//todo timesleep
					if err.Error()[:3] == "429" {
						_, r, err := multicall.Do(dex.ChainId, pairTokenInfosCalls)
						if err != nil {
							crawlFail := dao.CrawlFail{
								ChainId: dex.ChainId,
								Defi:    dex.Name,
								Address: pairToken.PairAddress,
								Type:    "CrawlNewPairToken",
								Errors:  err.Error(),
							}
							err := crawlFail.CallInsertFail()
							if err != nil {
								log.Println(log.LogLevelError, "CrawlNewPairToken crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
							}
							//todo fail pair token -> ++ count fail, address
							pairTokenFail += 1
							i += len(_pairTokenMethodList)
							continue
						} else {
							if r[0] != nil && r[1] != nil && r[2] != nil {
								pairToken.Token0 = strings.ToLower(r[0].(common.Address).String())
								pairToken.Token1 = strings.ToLower(r[1].(common.Address).String())

								if len(r[2].([]interface{})) > 2 {
									pairToken.Reserve0 = fmt.Sprintf("%d", r[2].([]interface{})[0].(*big.Int))
									pairToken.Reserve1 = fmt.Sprintf("%d", r[2].([]interface{})[1].(*big.Int))
									pairToken.BlockTimestampLast = int(r[2].([]interface{})[2].(uint32))
								}
							}
						}
					} else {
						crawlFail := dao.CrawlFail{
							ChainId: dex.ChainId,
							Defi:    dex.Name,
							Address: pairToken.PairAddress,
							Type:    "CrawlNewPairToken",
							Errors:  err.Error(),
						}
						err := crawlFail.CallInsertFail()
						if err != nil {
							log.Println(log.LogLevelError, "CrawlNewPairToken crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
						}
						//todo fail pair token -> ++ count fail, address
						pairTokenFail += 1
						i += len(_pairTokenMethodList)
						continue
					}
				} else {
					if res[0] != nil && res[1] != nil && res[2] != nil {
						pairToken.Token0 = strings.ToLower(res[0].(common.Address).String())
						pairToken.Token1 = strings.ToLower(res[1].(common.Address).String())

						if len(res[2].([]interface{})) > 2 {
							pairToken.Reserve0 = fmt.Sprintf("%d", res[2].([]interface{})[0].(*big.Int))
							pairToken.Reserve1 = fmt.Sprintf("%d", res[2].([]interface{})[1].(*big.Int))
							pairToken.BlockTimestampLast = int(res[2].([]interface{})[2].(uint32))
						}
					}
				}
			}

			err := pairToken.CallInsert()
			if err != nil {
				log.Println(log.LogLevelError, "CrawlNewPairToken pairToken.InsertPairToken()", err.Error())
				pairTokenFail += 1
				crawlFail := dao.CrawlFail{
					ChainId: dex.ChainId,
					Defi:    dex.Name,
					Address: pairToken.PairAddress,
					Type:    "CrawlNewPairToken",
					Errors:  err.Error(),
				}
				err := crawlFail.CallInsertFail()
				if err != nil {
					log.Println(log.LogLevelError, "CrawlNewPairToken crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
				}
			}

			exist := CheckExistToken(pairToken.Token0, dex.ChainId)
			if !exist {
				listToken = append(listToken, pairToken.Token0)
			}

			exist = CheckExistToken(pairToken.Token1, dex.ChainId)
			if !exist {
				listToken = append(listToken, pairToken.Token1)
			}

			i += len(_pairTokenMethodList)
		}
	}

	//remove duplicate token in listToken

	mapNewToken := make(map[string]string)
	listNewToken := []string{}
	for _, token := range listToken {
		_, exist := mapNewToken[token]
		if !exist {
			listNewToken = append(listNewToken, token)
			mapNewToken[token] = "exist"
		}
	}

	return listNewToken, len(listNewToken), pairTokenFail, nil
}

func CheckExistToken(address string, chainId string) bool {
	switch chainId {
	case constant.ETHEREUM_ID:
		_, exist := dao.MapTokenEthereum[address]
		// if !exist {
		// 	dao.MapTokenEthereum[address] = "exist"
		// }
		return exist

	case constant.BINANCE_SMART_CHAIN_ID:
		_, exist := dao.MapTokenBinance[address]
		// if !exist {
		// 	dao.MapTokenBinance[address] = "exist"
		// }
		return exist

	case constant.FANTOM_ID:
		_, exist := dao.MapTokenFantom[address]
		// if !exist {
		// 	dao.MapTokenFantom[address] = "exist"
		// }
		return exist

	case constant.CELO_ID:
		_, exist := dao.MapTokenCelo[address]
		// if !exist {
		// 	dao.MapTokenCelo[address] = "exist"
		// }
		return exist

	case constant.POLYGON_ID:
		_, exist := dao.MapTokenPolygon[address]
		// if !exist {
		// 	dao.MapTokenPolygon[address] = "exist"
		// }
		return exist

	case constant.AVALANCHE_C_CHAIN_ID:
		_, exist := dao.MapTokenAvalanche[address]
		// if !exist {
		// 	dao.MapTokenAvalanche[address] = "exist"
		// }
		return exist

	case constant.OPTIMISM_ID:
		_, exist := dao.MapTokenOptimism[address]
		// if !exist {
		// 	dao.MapTokenOptimism[address] = "exist"
		// }
		return exist

	case constant.ARBITRUM_ID:
		_, exist := dao.MapTokenArbitrum[address]
		// if !exist {
		// 	dao.MapTokenArbitrum[address] = "exist"
		// }
		return exist

	case constant.MOONBEAM_ID:
		_, exist := dao.MapTokenMoonbeam[address]
		// if !exist {
		// 	dao.MapTokenMoonbeam[address] = "exist"
		// }
		return exist

	case constant.KAVA_ID:
		_, exist := dao.MapTokenKava[address]
		// if !exist {
		// 	dao.MapTokenKava[address] = "exist"
		// }
		return exist

	case constant.CRONOS_ID:
		_, exist := dao.MapTokenCronos[address]
		// if !exist {
		// 	dao.MapTokenCronos[address] = "exist"
		// }
		return exist

	default:
		return true
	}
}
