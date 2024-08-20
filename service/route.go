package service

import (
	"fmt"
	"new-token/pkg/router"
	"new-token/service/index"
	"new-token/service/token/constant"
	"new-token/service/token/controller"
	"new-token/service/token/crawler"
)

// LoadRoutes to Load Routes to Router
func LoadRoutes() {
	fmt.Println("starting routes")
	go crawler.CrawlerScheduleTokenDex(constant.TIME_SLEEP_CRAWL_NEW_PAIR_DEX)

	// Set Endpoint for admin
	router.Router.Get(router.RouterBasePath+"/", index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)
	router.Router.Get(router.RouterBasePath+"/info", index.GetInfo)
	router.Router.Get(router.RouterBasePath+"/detail", controller.ScheduleCrawlTokenMonitor)
}
