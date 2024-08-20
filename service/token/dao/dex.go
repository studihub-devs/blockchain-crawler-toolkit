package dao

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"new-token/pkg/server"
	"strings"

	"github.com/google/uuid"
)

type Dex struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	FactoryAddress string    `json:"factoryAddress"`
	ChainId        string    `json:"chainId"`
	PairLength     int       `json:"pairLength"`
	Chainname      string    `json:"chainname"`
	UpdatedDate    string    `json:"updatedDate"`
}

type ListDex struct {
	Dexs []Dex
}

func (dex *Dex) CallDex() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	requestURL := urlConnectDatabase + "dex" + "?chainname=" + dex.Chainname + "&factoryAddress=" + dex.FactoryAddress
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

	responseDex := ResponseDex{}
	err = json.Unmarshal(resBody, &responseDex)
	if err != nil {
		return err
	}

	dex.PairLength = responseDex.Data.PairLength

	return nil
}

type ResponseDex struct {
	Status  bool   `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    Dex    `json:"data"`
}

func (listDex *ListDex) CallListDex(listChain []string) error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	requestURL := urlConnectDatabase + "listDex" + "?listChain=" + strings.Join(listChain, ",")
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

	responseListDex := ResponseListDex{}
	err = json.Unmarshal(resBody, &responseListDex)
	if err != nil {
		return err
	}

	for _, dexDTO := range responseListDex.Data.Dexs {
		dex := Dex{
			Id:             dexDTO.Id,
			Name:           dexDTO.Name,
			FactoryAddress: dexDTO.FactoryAddress,
			ChainId:        dexDTO.ChainId,
			PairLength:     dexDTO.PairLength,
			Chainname:      dexDTO.Chainname,
			UpdatedDate:    dexDTO.UpdatedDate,
		}
		listDex.Dexs = append(listDex.Dexs, dex)
	}

	return nil
}

type ResponseListDex struct {
	Status  bool    `json:"status"`
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Data    ListDex `json:"data"`
}

func (dex *Dex) CallUpdatePairLength() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	jsonBody, err := json.Marshal(dex)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := urlConnectDatabase + "dex"
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
