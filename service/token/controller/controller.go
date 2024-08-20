package controller

import (
	"net/http"
	"new-token/pkg/log"
	"new-token/pkg/router"
	"new-token/service/token/crawler"
	"new-token/service/token/dao"
)

func ScheduleCrawlTokenMonitor(w http.ResponseWriter, r *http.Request) {

	listCrawlerInfo := crawler.ListCrawlerInfo{}
	listDex := dao.ListDex{}

	listChain := []string{"binance"}
	err := listDex.CallListDex(listChain)
	if err != nil {
		log.Println(log.LogLevelError, "ScheduleCrawlTokenMonitor listDex.GetListDex()", err.Error())
	}

	for _, dex := range listDex.Dexs {
		crawlerInfo := crawler.CrawlerInfo{
			Chain:      dex.Chainname + " - " + dex.ChainId,
			Dex:        dex.Name,
			Factory:    dex.FactoryAddress,
			PairLength: dex.PairLength,
		}

		dailytoken := 0
		dailyPostES := 0
		dailyPushES := 0

		if crawler.MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["LastestCrawler"] != nil {
			crawlerInfo.LastestCrawler = crawler.MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["LastestCrawler"].(string)
		}

		if crawler.MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["LastestIndex"] != nil {
			crawlerInfo.PairLength = crawler.MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["LastestIndex"].(int)
		}

		if crawler.MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["Info"] != nil {
			crawlerInfo.Info = crawler.MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["Info"].([]crawler.InfoCrawl)

			for _, infoCrawl := range crawler.MapCrawlerInfo[dex.Chainname+" - "+dex.Name+" - "+dex.ChainId]["Info"].([]crawler.InfoCrawl) {
				dailytoken += infoCrawl.NewToken
				dailyPostES += infoCrawl.PostES
				dailyPushES += infoCrawl.PushQueue
			}
			crawlerInfo.DailyToken = dailytoken
			crawlerInfo.DailyPostES = dailyPostES
			crawlerInfo.DailyPushQueue = dailyPushES

		}
		listCrawlerInfo.Data = append(listCrawlerInfo.Data, crawlerInfo)
	}
	router.ResponseSuccessWithData(w, "200", "Get Info Succesfully - eth", listCrawlerInfo)
}
