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

const (
	_allPairsLengthMethod = "allPairsLength"
	_allPairsMethod       = "allPairs"

	// _positionsCallMutilcalls = 40

	_allPairsCallPerMulticalls         = 200
	_pairTokenCallPerMulticalls        = 100
	_pairTokenCallReservesPerMulticals = 200
	_tokenInfoCallPerMulticalls        = 100

	_decimalsMethod    = "decimals"
	_totalSupplyMethod = "totalSupply"
	_token0Method      = "token0"
	_token1Method      = "token1"

	// _tokenByIndex    = "tokenByIndex"
	// _positionsMethod = "positions"

	_nameMethod   = "name"
	_symbolMethod = "symbol"

	_reservesMethod = "getReserves"
)

var CountFail int
var CountSuccess int
var PairLength int

var _pairTokenMethodList = []string{_token0Method, _token1Method, _reservesMethod}

var _tokenInfoMethodList = []string{_nameMethod, _symbolMethod, _decimalsMethod, _totalSupplyMethod}

var (
	CurrentDexCrawl     = ""
	CurrentChainIdCrawl = ""
	Stages              = ""

	CurrentTokenCount = 0
	Speed             = 0.0
	Start             = 0
)

func CallPairlength(dex dao.Dex) (int, error) {
	allPairsLengthCall := multicall.NewCall(
		common.HexToAddress(dex.FactoryAddress),
		_allPairsLengthMethod,
		[]interface{}{},
		factoryABI,
	)
	_, result, err := multicall.Do(dex.ChainId, []multicall.Call{*allPairsLengthCall})
	if err != nil {
		crawlFail := dao.CrawlFail{
			ChainId: dex.ChainId,
			Defi:    dex.Name,
			Address: dex.FactoryAddress,
			Type:    "CallPairlength",
			Errors:  err.Error(),
		}
		err := crawlFail.CallInsertFail()
		if err != nil {
			log.Println(log.LogLevelError, "CallPairlength crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
		}
		return 0, err
	}
	// fmt.Println(result...)
	pairsLength := result[0].(*big.Int).Int64()

	return int(pairsLength), nil
}

/*
Get last index of latest pair in db by dex(factoryAddress and chainId).
Call rpc get pairLength.
If last index < pairLength -> new pool -> get new pairAddress of pool by index(from last index in db to pairLength).

Return list new pairAddress pool, pairLength(use pairLength to update last index in db)
*/
func CrawlNewPairAddress(dex dao.Dex) ([]string, int, int, error) {

	listPairAddress := []string{}
	pairAddressFail := 0
	// err := dex.GetDex()
	// if err != nil {
	// 	log.Println(log.LogLevelError, dex.Name+" dex.GetDex()", err.Error())
	// 	return listPairAddress, dex.PairLength, err
	// }

	pairsLength, err := CallPairlength(dex)
	if err != nil {
		log.Println(log.LogLevelError, dex.Name+" CallPairlength(dex)", err.Error())
		return listPairAddress, pairsLength, pairAddressFail, err
	}

	if pairsLength == dex.PairLength {
		// log.Println(log.LogLevelInfo, dex.Name+" No new pool", pairsLength)
		return listPairAddress, pairsLength, pairAddressFail, err
	}
	startIndex := dex.PairLength
	endIndex := pairsLength

	pairAddressListBigCalls := [][]multicall.Call{}
	pairAddressListCalls := []multicall.Call{}
	for index := int64(startIndex); index < int64(endIndex); index++ {
		call := multicall.NewCall(
			common.HexToAddress(dex.FactoryAddress),
			_allPairsMethod,
			[]interface{}{big.NewInt(index)},
			factoryABI,
		)
		pairAddressListCalls = append(pairAddressListCalls, *call)
		if len(pairAddressListCalls) == _allPairsCallPerMulticalls || index == int64(endIndex)-1 {
			pairAddressListBigCalls = append(pairAddressListBigCalls, pairAddressListCalls)
			pairAddressListCalls = []multicall.Call{}
		}
	}

	newPair := endIndex - startIndex
	log.Println(log.LogLevelInfo, fmt.Sprintf("%s chain %s new pair index start: %d - end: %d", dex.Name, dex.ChainId, startIndex, endIndex), newPair)

	// infoCrawl := InfoCrawl{
	// 	NewPair:         newPair,
	// 	PairLength:      end,
	// 	NewToken:        0,
	// 	PairAddressFail: 0,
	// 	TokenFail:       0,
	// 	PairTokenFail:   0,
	// }

	for index, calls := range pairAddressListBigCalls {
		_, result, _ := multicall.Do(dex.ChainId, calls)

		for i, address := range result {
			pairAddress := dao.PairAddress{
				ChainId:   dex.ChainId,
				ChainName: dex.Chainname,
				Defi:      dex.Name,
				Index:     startIndex + index*_allPairsCallPerMulticalls + i,
			}
			if address != nil {
				pairAddress.Address = strings.ToLower(address.(common.Address).String())
			} else {
				call := multicall.NewCall(
					common.HexToAddress(dex.FactoryAddress),
					_allPairsMethod,
					[]interface{}{big.NewInt(int64(pairAddress.Index))},
					factoryABI,
				)
				_, res, err := multicall.Do(dex.ChainId, []multicall.Call{*call})
				if err != nil {
					//todo check error and time.sleep
					if err.Error()[:3] == "429" { //status code rate limit
						time.Sleep(30 * time.Second)
						_, r, err := multicall.Do(dex.ChainId, []multicall.Call{*call})
						if err != nil {
							crawlFail := dao.CrawlFail{
								ChainId: dex.ChainId,
								Defi:    dex.Name,
								Address: dex.FactoryAddress,
								Index:   pairAddress.Index,
								Type:    "PairAddress",
								Errors:  err.Error(),
							}
							err := crawlFail.CallInsertFail()
							if err != nil {
								log.Println(log.LogLevelError, "CrawlNewPairAddress crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
							}
							pairAddressFail += 1
							//todo fail pair -> ++ count fail, address
						} else {
							if r[0] != nil {
								pairAddress.Address = strings.ToLower(r[0].(common.Address).String())
							}
						}
					} else {
						crawlFail := dao.CrawlFail{
							ChainId: dex.ChainId,
							Defi:    dex.Name,
							Address: dex.FactoryAddress,
							Index:   pairAddress.Index,
							Type:    "PairAddress",
							Errors:  err.Error(),
						}
						err := crawlFail.CallInsertFail()
						if err != nil {
							log.Println(log.LogLevelError, "CrawlNewPairAddress crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
						}
						//todo fail pair -> ++ count fail, address
						pairAddressFail += 1
						// infoCrawl.ListPairAddressFail = append(infoCrawl.ListPairAddressFail, string(pairAddress.Index))
					}
				} else {
					if res[0] != nil {
						pairAddress.Address = strings.ToLower(res[0].(common.Address).String())
					}
				}
			}
			err := pairAddress.CallInsert()
			if err != nil {
				log.Println(log.LogLevelError, "CrawlNewPairAddress pairAddress.CallInsert()", err.Error())
				//todo fail pair -> ++ count fail, address
				crawlFail := dao.CrawlFail{
					ChainId: dex.ChainId,
					Defi:    dex.Name,
					Address: dex.FactoryAddress,
					Index:   pairAddress.Index,
					Type:    "PairAddress",
					Errors:  err.Error(),
				}
				err := crawlFail.CallInsertFail()
				if err != nil {
					log.Println(log.LogLevelError, "CrawlNewPairAddress crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
				}
				pairAddressFail += 1
			} else {
				listPairAddress = append(listPairAddress, pairAddress.Address)
			}

		}
	}

	return listPairAddress, pairsLength, pairAddressFail, nil
}
