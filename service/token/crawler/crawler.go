package crawler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"new-token/pkg/log"
	"new-token/pkg/ms"
	"new-token/pkg/server"
	"new-token/pkg/utils"
	"new-token/service/token/constant"
	"new-token/service/token/contract"
	"new-token/service/token/dao"
	"time"
)

var UpdateCount int
var InsertCount int
var PostES int
var PushQueue int
var LastUpdated string
var DayStart string

type Tracking struct {
	LastUpdateTime string
	Update         int
	Insert         int
	DayStart       string
}

type CrawlerInfo struct {
	Chain          string
	Dex            string
	Factory        string
	LastestCrawler string
	PairLength     int
	Info           []InfoCrawl
	DailyToken     int
	DailyPostES    int
	DailyPushQueue int
}

type InfoCrawl struct {
	Time       string
	NewPair    int
	PairLength int

	PairAddressFail int
	// ListPairAddressFail []string
	PairTokenFail int
	// ListPairTokenFail   []string

	NewToken int
	// ListToken     []string
	TokenFail int
	// ListTokenFail []string
	PostES    int
	PushQueue int
}

type ListCrawlerInfo struct {
	Data []CrawlerInfo
}

var MapCrawlerInfo = map[string]map[string]any{}

func init() {
	fmt.Println("init")

	//todo: get chain list

	dao.MapChainList = make(map[string]string, 0)

	err := dao.CallChainList()
	if err != nil {
		//todo: fatal
		log.Println(log.LogLevelFatal, "dao.CallChainList", err.Error())
	}

	log.Println(log.LogLevelInfo, `Get chainlist succes`, len(dao.MapChainList))

	//Get all token in db to check exist

	//init map to check exist
	dao.MapTokenEthereum = make(map[string]string)
	dao.MapTokenBinance = make(map[string]string)
	dao.MapTokenFantom = make(map[string]string)
	dao.MapTokenCelo = make(map[string]string)
	dao.MapTokenPolygon = make(map[string]string)
	dao.MapTokenAvalanche = make(map[string]string)
	dao.MapTokenOptimism = make(map[string]string)
	dao.MapTokenArbitrum = make(map[string]string)
	dao.MapTokenMoonbeam = make(map[string]string)
	dao.MapTokenKava = make(map[string]string)
	dao.MapTokenCronos = make(map[string]string)

	listChain := []string{"binance"}

	//call api to get all address and assign into map
	err = dao.CallAllCryptoEVM(listChain)
	if err != nil {
		//todo: fatal
		log.Println(log.LogLevelFatal, "listCrypto.GetAllCryptoEVM()", err.Error())
	}

	log.Println(log.LogLevelInfo, "Map crypto ethereum: ", len(dao.MapTokenEthereum))
	log.Println(log.LogLevelInfo, "Map crypto binance: ", len(dao.MapTokenBinance))
	log.Println(log.LogLevelInfo, "Map crypto fantom: ", len(dao.MapTokenFantom))
	log.Println(log.LogLevelInfo, "Map crypto celo: ", len(dao.MapTokenCelo))
	log.Println(log.LogLevelInfo, "Map crypto polygon: ", len(dao.MapTokenPolygon))
	log.Println(log.LogLevelInfo, "Map crypto avalanche: ", len(dao.MapTokenAvalanche))
	log.Println(log.LogLevelInfo, "Map crypto optimism: ", len(dao.MapTokenOptimism))
	log.Println(log.LogLevelInfo, "Map crypto arbitrum: ", len(dao.MapTokenArbitrum))
	log.Println(log.LogLevelInfo, "Map crypto moonbeam: ", len(dao.MapTokenMoonbeam))
	log.Println(log.LogLevelInfo, "Map crypto kava: ", len(dao.MapTokenKava))
	log.Println(log.LogLevelInfo, "Map crypto cronos: ", len(dao.MapTokenCronos))

}

func CrawlerScheduleTokenDex(timeSleep time.Duration) {
	listDex := dao.ListDex{}

	// fmt.Println("sleep 10h")
	// time.Sleep(10 * time.Hour)

	UpdateCount = 0
	InsertCount = 0

	nextDay := false
	DayStart = utils.TimeNowStringVietNam()
	go func() {
		for {

			nextDay = false
			LastUpdated = utils.TimeNowStringVietNam()
			MapCrawlerInfo = make(map[string]map[string]any)

			//get list dex in db to start crawl
			listChain := []string{"binance"}
			err := listDex.CallListDex(listChain)

			if err != nil {
				log.Println(log.LogLevelError, "CrawlerScheduleTokenDex listDefi.GetListDefi()", err.Error())
				time.Sleep(5 * time.Minute)
				if nextDay {
					nextDay = false
					break
				}
				continue
			}

			for _, dex := range listDex.Dexs {
				// crawlerInfo := CrawlerInfo{
				// 	Chain:   dex.ChainId,
				// 	Dex:     GetChainName(dex.ChainId),
				// 	Factory: dex.FactoryAddress,
				// }
				// CrawlerData.Data = append(CrawlerData.Data, crawlerInfo)
				MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId] = make(map[string]any)
				MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["Factory"] = dex.FactoryAddress
				MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["LastestIndex"] = dex.PairLength
			}
			for {
				//loop through list dex

				LastUpdated = utils.TimeNowStringVietNam()

				for _, dex := range listDex.Dexs {
					PostES = 0
					PushQueue = 0

					err := dex.CallDex()
					if err != nil {
						log.Println(log.LogLevelError, "CrawlerScheduleTokenDex dex.GetDex() "+dex.Chainname+" "+dex.FactoryAddress, err.Error())
						continue
					}

					indexYesterday := dex.PairLength
					infoCrawl := InfoCrawl{
						Time: utils.TimeNowStringVietNam(),
					}

					MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["LastestCrawler"] = utils.TimeNowStringVietNam()
					MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["LastestIndex"] = dex.PairLength

					listNewPairAddress, pairLength, pairAddressFail, err := contract.CrawlNewPairAddress(dex)
					if err != nil {
						log.Println(log.LogLevelError, "CrawlerScheduleTokenDex contract.CrawlNewPairAddress(dex) "+dex.ChainId+" "+dex.FactoryAddress, err.Error())
						continue
					}

					dex.PairLength = pairLength
					//todo: update pair length
					err = dex.CallUpdatePairLength()
					if err != nil {
						log.Println(log.LogLevelError, "CrawlerScheduleTokenDex dex.UpdatePairLength() "+dex.ChainId+" "+dex.FactoryAddress, err.Error())
						//todo: log new file
						log.CrawlFailWrite(log.LogLevelError, fmt.Sprintf("Error: Update pair length dex:%v chain:%v pairLength: %v", dex.Name, dex.ChainId, dex.PairLength), err)
						continue
					}

					infoCrawl.NewPair = pairLength - indexYesterday
					infoCrawl.PairLength = pairLength
					infoCrawl.PairAddressFail = pairAddressFail

					listNewPairToken, newToken, pairTokenFail, err := contract.CrawlNewPairToken(dex, listNewPairAddress)
					if err != nil {
						log.Println(log.LogLevelError, "CrawlerScheduleTokenDex contract.CrawlNewPairToken(dex, listNewPairAddress) "+dex.ChainId+" "+dex.FactoryAddress, err.Error())
						continue
					}

					infoCrawl.NewToken = newToken
					infoCrawl.PairTokenFail = pairTokenFail

					listCrypto, newToken, tokenFail, err := contract.CrawlInfoNewToken(dex, listNewPairToken)
					if err != nil {
						log.Println(log.LogLevelError, "CrawlerScheduleTokenDex contract.CrawlInfoNewToken(dex, listNewPairToken) "+dex.ChainId+" "+dex.FactoryAddress, err.Error())
						continue
					}

					infoCrawl.NewToken = newToken
					infoCrawl.PairTokenFail = tokenFail

					//name, symbol, totalSupply, Decimals

					for _, crypto := range listCrypto.Cryptos {
						if dex.Name == "Pancakeswap bsc" {
							crypto.IsPancakeswap = true
						}

						crypto.Type = "token"
						crypto.CryptoSrc = "dex"
						crypto.CryptoId = "gear5_token_" + dao.MapChainList[crypto.ChainId] + "_" + crypto.Address
						crypto.ChainName = dao.MapChainList[crypto.ChainId]
						crypto.SubCategory = constant.GenSubcategoryByChainname(crypto.ChainName)
						crypto.Createddate = utils.TimeNowStringVietNam()
						crypto.Updateddate = utils.TimeNowStringVietNam()

						crypto.IsCoinbase = false
						crypto.IsCoinmarketcap = false
						crypto.IsCoinbase = false
						crypto.IsBinance = false
						crypto.IsDex = true
						crypto.IsVerifiedByAdmin = true
						crypto.IsShow = true

						//todo: insert db
						err := crypto.CallInsert()
						if err != nil {
							log.Println(log.LogLevelError, "CrawlerScheduleTokenDex crypto.Insert() "+crypto.ChainId+" "+crypto.Address, err.Error())
							fail := dao.CrawlFail{
								Defi:    dex.Name,
								ChainId: crypto.ChainId,
								Address: crypto.Address,
								Type:    "InsertCrypto",
								Errors:  err.Error(),
							}
							err := fail.CallInsertFail()
							if err != nil {
								log.Println(log.LogLevelError, "CrawlerScheduleTokenDex fail.CallInsertFail() "+crypto.ChainId+" "+crypto.Address, err.Error())
							}
							continue
						} else {
							AppendMapCrypto(crypto)
						}

						//todo: append redis queue to update
						err = AppendCryptoToMessageQueueRedis(crypto)
						if err != nil {
							log.Println(log.LogLevelError, "CrawlerScheduleTokenDex AppendCryptoToMessageQueueRedis(crypto) "+crypto.ChainId+" "+crypto.Address, err.Error())
							fail := dao.CrawlFail{
								Defi:    dex.Name,
								ChainId: crypto.ChainId,
								Address: crypto.Address,
								Type:    "AppendCryptoToMessageQueueRedis",
								Errors:  err.Error(),
							}
							err := fail.CallInsertFail()
							if err != nil {
								log.Println(log.LogLevelError, "CrawlerScheduleTokenDex fail.CallInsertFail() "+crypto.ChainId+" "+crypto.Address, err.Error())
							}
						} else {
							PushQueue += 1
						}

						// fmt.Println(crypto.Id, crypto.Address)

						//todo: post api elasticsearch
						crypto.AddressShow = crypto.Address
						err = PostCryptoToElasticsearch(crypto)
						if err != nil {
							log.Println(log.LogLevelError, "CrawlerScheduleTokenDex PostCryptoToElasticsearch(crypto) "+crypto.ChainId+" "+crypto.Address, err.Error())
							fail := dao.CrawlFail{
								Defi:    dex.Name,
								ChainId: crypto.ChainId,
								Address: crypto.Address,
								Type:    "PostCryptoToElasticsearch",
								Errors:  err.Error(),
							}
							err := fail.CallInsertFail()
							if err != nil {
								log.Println(log.LogLevelError, "CrawlerScheduleTokenDex fail.CallInsertFail() "+crypto.ChainId+" "+crypto.Address, err.Error())
							}
						} else {
							PostES += 1
						}
					}

					InsertCount += infoCrawl.NewToken
					infoCrawl.PostES = PostES
					infoCrawl.PushQueue = PushQueue

					if MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["Info"] != nil {
						listInfoCrawl := []InfoCrawl{}

						listInfoCrawl = append(listInfoCrawl, MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["Info"].([]InfoCrawl)...)
						if infoCrawl.NewPair != 0 {
							listInfoCrawl = append(listInfoCrawl, infoCrawl)

						}
						MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["Info"] = listInfoCrawl

					} else {
						listInfoCrawl := []InfoCrawl{}

						listInfoCrawl = append(listInfoCrawl, infoCrawl)

						MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["Info"] = listInfoCrawl
						// fmt.Sprintf("%s,%d,%d,%d,%s", utils.TimeNowString(), newPair, pairsLength, newToken, strings.Join(listToken, ","))
					}

				}

				time.Sleep(20 * time.Second)
				if nextDay {
					nextDay = false
					break
				}
			}
		}
	}()

	go func() {
		for {
			time.Sleep(timeSleep)
			nextDay = true
			InsertCount = 0
			UpdateCount = 0
			PostES = 0
			PushQueue = 0

			DayStart = utils.TimeNowStringVietNam()
		}
	}()
}

func AppendCryptoToMessageQueueRedis(crypto dao.Crypto) error {
	byteValue, err := json.Marshal(crypto)
	if err != nil {
		return err
	}
	err = ms.Redis.Store.RPush("new-crypto", byteValue).Err()
	if err != nil {
		return err
	}

	return nil
}

type BodyInsertCrypto struct {
	ApiKey     string       `json:"apiKey"`
	ListCrypto []dao.Crypto `json:"listCrypto"`
}

func PostCryptoToElasticsearch(crypto dao.Crypto) error {
	// fmt.Println(crypto.Address)

	bodyInsertCrypto := BodyInsertCrypto{
		ApiKey:     constant.API_KEY_SEARCH,
		ListCrypto: []dao.Crypto{crypto},
	}

	jsonBody, err := json.Marshal(bodyInsertCrypto)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	urlSearch := server.Config.GetString("URL_SEARCH")

	requestURL := urlSearch + "insert/crypto"
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// fmt.Printf("client: got response!\n")
	// fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 || res.StatusCode < 200 {
		return errors.New(res.Status + " " + string(resBody))
	}

	return nil
}
func AppendMapCrypto(crypto dao.Crypto) {
	if crypto.ChainId == "1" { //ethereum
		dao.MapTokenEthereum[crypto.Address] = "exist"
	}
	if crypto.ChainId == "10" { //optimism
		dao.MapTokenOptimism[crypto.Address] = "exist"
	}
	if crypto.ChainId == "1284" { //moonbeam
		dao.MapTokenMoonbeam[crypto.Address] = "exist"
	}
	if crypto.ChainId == "137" { //polygon
		dao.MapTokenPolygon[crypto.Address] = "exist"
	}
	if crypto.ChainId == "2222" { //kava
		dao.MapTokenKava[crypto.Address] = "exist"
	}
	if crypto.ChainId == "25" { //cronos
		dao.MapTokenCronos[crypto.Address] = "exist"
	}
	if crypto.ChainId == "250" { //fantom
		dao.MapTokenFantom[crypto.Address] = "exist"
	}
	if crypto.ChainId == "42161" { //arbitrum
		dao.MapTokenArbitrum[crypto.Address] = "exist"
	}
	if crypto.ChainId == "42220" { //celo
		dao.MapTokenCelo[crypto.Address] = "exist"
	}
	if crypto.ChainId == "43114" { //avalanche
		dao.MapTokenAvalanche[crypto.Address] = "exist"
	}
	if crypto.ChainId == "56" { //binance
		dao.MapTokenBinance[crypto.Address] = "exist"
	}
}
