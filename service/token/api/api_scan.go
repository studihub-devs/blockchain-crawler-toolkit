package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"new-token/pkg/log"
	"new-token/pkg/utils"
	"new-token/service/token/constant"
	"new-token/service/token/dao"
	"strings"
	"time"
)

func CallApiScanIsProxyContractVerified(crypto dao.Crypto) dao.Crypto {
	time.Sleep(230 * time.Millisecond)
	cryptoCallApiScan := dao.Crypto{}
	contractCode, err := CallApiContractCode(crypto)
	if err != nil {
		log.Println(log.LogLevelError, "CallApiContractCode(crypto)"+crypto.ChainId+" "+crypto.Address, err.Error())
		crawlFail := dao.CrawlFail{
			ChainId: crypto.ChainId,
			Address: strings.ToLower(crypto.Address),
			Type:    "CallApiContractCode",
			Errors:  err.Error(),
		}
		err := crawlFail.CallInsertFail()
		if err != nil {
			log.Println(log.LogLevelError, "CallApiScanIsProxyContractVerified crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
		}
	}

	var contractVerified bool
	if contractCode.ABI == "" {
		log.Println(log.LogLevelError, "CallApiContractCode(crypto)"+crypto.ChainId+" "+crypto.Address, "contractCode.ABI is empty")
	} else {
		if contractCode.ABI == "Contract source code not verified" {
			contractVerified = false
		} else {
			// verified
			contractVerified = true
		}
	}

	var isProxy bool

	if contractCode.Proxy == "1" || contractCode.Proxy == "0" {
		if contractCode.Proxy == "1" {
			isProxy = true
		}
		if contractCode.Proxy == "0" {
			isProxy = false
		}
	} else {
		log.Println(log.LogLevelError, "CallApiContractCode(crypto)"+crypto.ChainId+" "+crypto.Address, "contractCode.Isproxy is valid")
		crawlFail := dao.CrawlFail{
			ChainId: crypto.ChainId,
			Address: strings.ToLower(crypto.Address),
			Type:    "CallApiContractCode",
			Errors:  "contractCode.Isproxy is valid",
		}
		err := crawlFail.CallInsertFail()
		if err != nil {
			log.Println(log.LogLevelError, "CallApiScanIsProxyContractVerified crawlFail.CallInsertFail()"+crawlFail.ChainId+":"+crawlFail.Address, err.Error())
		}
	}
	cryptoCallApiScan.IsProxy = isProxy
	cryptoCallApiScan.ContractVerified = contractVerified

	return cryptoCallApiScan
}

func CallApiContractCode(crypto dao.Crypto) (ContractCode, error) {
	contractCode := ContractCode{}
	countFail := 20

	urlScan, apikey, _ := constant.GetExplorerByChainId(crypto.ChainId)

	for {
		url := fmt.Sprintf(urlScan, crypto.Address, apikey)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Println(log.LogLevelDebug, "CallApiContractCode http.NewRequest "+crypto.ChainId+" "+crypto.Address, err.Error())
			time.Sleep(5 * time.Second)
			countFail -= 1
			if countFail <= 0 {
				return contractCode, err
			}
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Println(log.LogLevelDebug, "CallApiContractCode client.Do "+crypto.ChainId+" "+crypto.Address, err.Error())
			time.Sleep(5 * time.Second)
			countFail -= 1
			if countFail <= 0 {
				return contractCode, err
			}
			continue
		}
		var body []byte
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return contractCode, err
		}
		defer res.Body.Close()

		responseContractCode := ResponseContractCode{}
		response := make(map[string]any)

		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println(log.LogLevelDebug, "CallApiContractCode json.Unmarshal "+crypto.ChainId+" "+crypto.Address, err.Error())
			time.Sleep(5 * time.Second)
			countFail -= 1
			if countFail <= 0 {
				return contractCode, err
			}
			continue
		}

		utils.Mapping(response, &responseContractCode)

		if response["status"] == "0" {

			log.Println(log.LogLevelError, "CallApiContractCode response status = 0 "+crypto.ChainId+" "+crypto.Address, response["result"].(string))
			time.Sleep(10 * time.Second)
			countFail -= 1
			if countFail <= 0 {
				return contractCode, errors.New(response["result"].(string))
			}
			continue
		}

		if len(responseContractCode.Result) > 0 {
			contractCode.ABI = responseContractCode.Result[0].ABI
			contractCode.SourceCode = responseContractCode.Result[0].SourceCode
			contractCode.Proxy = responseContractCode.Result[0].Proxy
		}

		// fmt.Println("alo", contractCode)

		// fmt.Println("abi", contractCode.ABI[:5])
		// fmt.Println("code", contractCode.SourceCode[:5])
		// fmt.Println("proxy", contractCode.Proxy)

		return contractCode, nil
	}

}
