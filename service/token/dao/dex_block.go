package dao

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"new-token/pkg/server"

	"github.com/google/uuid"
)

type DexBlock struct {
	Id             uuid.UUID
	Name           string
	FactoryAddress string
	ChainId        string
	Block          int
	Chainname      string
	UpdatedDate    string
}

func (dexBlock *DexBlock) CallDexBlock() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	requestURL := urlConnectDatabase + "dexBlock" + "?chainname=" + dexBlock.Chainname + "&factoryAddress=" + dexBlock.FactoryAddress
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
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

	responseDexBlock := ResponseDexBlock{}
	err = json.Unmarshal(resBody, &responseDexBlock)
	if err != nil {
		return err
	}

	dexBlock.Name = responseDexBlock.Data.Name
	dexBlock.FactoryAddress = responseDexBlock.Data.FactoryAddress
	dexBlock.ChainId = responseDexBlock.Data.ChainId
	dexBlock.Block = responseDexBlock.Data.Block
	dexBlock.Chainname = responseDexBlock.Data.Chainname

	return nil
}

type ResponseDexBlock struct {
	Status  bool     `json:"status"`
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Data    DexBlock `json:"data"`
}

func (dexBlock *DexBlock) CallUpdateBlock() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	jsonBody, err := json.Marshal(dexBlock)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := urlConnectDatabase + "dexBlock"
	req, err := http.NewRequest(http.MethodPatch, requestURL, bodyReader)
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
