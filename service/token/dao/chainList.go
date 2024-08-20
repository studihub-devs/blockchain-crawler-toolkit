package dao

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"new-token/pkg/server"
)

var MapChainList map[string]string

func CallChainList() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	requestURL := urlConnectDatabase + "chainList"
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 || res.StatusCode < 200 {
		return errors.New(res.Status + " " + string(resBody))
	}

	responseChainList := ResponseChainList{}
	err = json.Unmarshal(resBody, &responseChainList)
	if err != nil {
		return err
	}

	for chainId, chainName := range responseChainList.Data {
		MapChainList[chainId] = chainName
	}

	return nil
}

type ResponseChainList struct {
	Status  bool              `json:"status"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}
