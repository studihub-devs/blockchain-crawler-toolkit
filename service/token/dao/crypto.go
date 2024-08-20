package dao

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"new-token/pkg/server"
	"strings"
	"time"

	"github.com/google/uuid"
)

var MapTokenEthereum map[string]string
var MapTokenBinance map[string]string
var MapTokenFantom map[string]string
var MapTokenCelo map[string]string
var MapTokenPolygon map[string]string
var MapTokenAvalanche map[string]string
var MapTokenOptimism map[string]string
var MapTokenArbitrum map[string]string
var MapTokenMoonbeam map[string]string
var MapTokenKava map[string]string
var MapTokenCronos map[string]string

type Crypto struct {
	Id          uuid.UUID
	CryptoId    string `json:"cryptoid"`
	CryptoSrc   string `json:"cryptosrc"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimal     uint8  `json:"decimal"`
	Address     string `json:"address"`
	AddressShow string `json:"addressShow"` //for elasticsearch
	ChainId     string `json:"chainId"`
	ChainName   string `json:"chainName"`
	BigLogo     string `json:"bigLogo"`

	Socials          map[string]string `json:"socials"`
	Website          string            `json:"website"`
	Explorer         string            `json:"explorer"`
	ContractVerified bool              `json:"contractVerified"`
	Holders          int64             `json:"holders"`

	IsScam            bool `json:"isScam"`
	IsVerifiedByAdmin bool `json:"isVerifiedByAdmin"`
	IsShow            bool `json:"isShow"`
	IsWarning         bool `json:"isWarning"`
	IsProxy           bool `json:"isProxy"`

	Proof                          Proof    `json:"proof"`
	TotalSupplyBeforeDivideDecimal *big.Int `json:"totalSupplyBeforeDivideDecimal"`
	TotalSupply                    string   `json:"totalSupply"`
	SubCategory                    string   `json:"subcategory"`
	Multichain                     []Chain  `json:"multichain"`

	Type string `json:"type"`

	IsCoingecko     bool `json:"isCoingecko"`
	IsCoinmarketcap bool `json:"isCoinmarketcap"`
	IsBinance       bool `json:"isBinance"`
	IsCoinbase      bool `json:"isCoinbase"`
	IsPancakeswap   bool `json:"isPancakeswap"`
	IsUniswap       bool `json:"isUniswap"`
	IsDex           bool `json:"isDex"`

	Createddate string `json:"createdDate"`
	Updateddate string `json:"updatedDate"`
}

type Chain struct {
	Symbol    string `json:"symbol"`
	Address   string `json:"address"`
	Decimal   int64  `json:"decimal"`
	CryptoId  string `json:"cryptoId"`
	ChainName string `json:"chainName"`
	CryptoSrc string `json:"cryptoSrc"`
}

type ListCrypto struct {
	Cryptos []Crypto `json:"cryptos"`
}

type Proof struct {
	IsWarning *string `json:"isWarning"`
	IsScam    *string `json:"isScam"`
}

func CallAllCryptoEVM(listChain []string) error {
	limit := 100000
	offset := 0 //phaisua
	for {
		urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

		requestURL := urlConnectDatabase + "mapCrypto" + "?listChain=" + strings.Join(listChain, ",") + fmt.Sprintf("&limit=%d&offset=%d", limit, offset)

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

		responseMapCrypto := ResponseMapCrypto{}
		err = json.Unmarshal(resBody, &responseMapCrypto)
		if err != nil {
			return err
		}

		if len(responseMapCrypto.Data.MapTokenEthereum) == 0 &&
			len(responseMapCrypto.Data.MapTokenBinance) == 0 &&
			len(responseMapCrypto.Data.MapTokenFantom) == 0 &&
			len(responseMapCrypto.Data.MapTokenCelo) == 0 &&
			len(responseMapCrypto.Data.MapTokenPolygon) == 0 &&
			len(responseMapCrypto.Data.MapTokenAvalanche) == 0 &&
			len(responseMapCrypto.Data.MapTokenOptimism) == 0 &&
			len(responseMapCrypto.Data.MapTokenArbitrum) == 0 &&
			len(responseMapCrypto.Data.MapTokenMoonbeam) == 0 &&
			len(responseMapCrypto.Data.MapTokenKava) == 0 &&
			len(responseMapCrypto.Data.MapTokenCronos) == 0 {
			break
		}

		for address, exist := range responseMapCrypto.Data.MapTokenEthereum {
			MapTokenEthereum[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenBinance {
			MapTokenBinance[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenFantom {
			MapTokenFantom[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenCelo {
			MapTokenCelo[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenPolygon {
			MapTokenPolygon[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenAvalanche {
			MapTokenAvalanche[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenOptimism {
			MapTokenOptimism[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenArbitrum {
			MapTokenArbitrum[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenMoonbeam {
			MapTokenMoonbeam[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenKava {
			MapTokenKava[address] = exist
		}

		for address, exist := range responseMapCrypto.Data.MapTokenCronos {
			MapTokenCronos[address] = exist
		}
		offset += 100000 //phaisua
		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

type ResponseMapCrypto struct {
	Status  bool      `json:"status"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    MapCrypto `json:"data"`
}

type MapCrypto struct {
	MapTokenEthereum  map[string]string `json:"mapTokenEthereum"`
	MapTokenBinance   map[string]string `json:"mapTokenBinance"`
	MapTokenFantom    map[string]string `json:"mapTokenFantom"`
	MapTokenCelo      map[string]string `json:"mapTokenCelo"`
	MapTokenPolygon   map[string]string `json:"mapTokenPolygon"`
	MapTokenAvalanche map[string]string `json:"mapTokenAvalanche"`
	MapTokenOptimism  map[string]string `json:"mapTokenOptimism"`
	MapTokenArbitrum  map[string]string `json:"mapTokenArbitrum"`
	MapTokenMoonbeam  map[string]string `json:"mapTokenMoonbeam"`
	MapTokenKava      map[string]string `json:"mapTokenKava"`
	MapTokenCronos    map[string]string `json:"mapTokenCronos"`
}

func (crypto *Crypto) CallInsert() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	jsonBody, err := json.Marshal(crypto)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := urlConnectDatabase + "crypto"
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

	responseInsertCrypto := ResponseInsertCrypto{}
	err = json.Unmarshal(resBody, &responseInsertCrypto)
	if err != nil {
		return err
	}

	crypto.Id = responseInsertCrypto.Data.Id

	return nil
}

type ResponseInsertCrypto struct {
	Status  bool   `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    Crypto `json:"data"`
}
