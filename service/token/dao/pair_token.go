package dao

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"new-token/pkg/server"

	"github.com/google/uuid"
)

type PairToken struct {
	Id                     uuid.UUID `json:"id"`
	ChainId                string    `json:"chainId"`
	Defi                   string    `json:"defi"`
	PairAddress            string    `json:"pairAddress"`
	ChainName              string    `json:"chainName"`
	Token0                 string    `json:"token0"`
	Token1                 string    `json:"token1"`
	Reserve0               string    `json:"reserve0"`
	Reserve1               string    `json:"reserve1"`
	BlockTimestampLast     int       `json:"blockTimestampLast"`
	ReserveUSD             string    `json:"reserveUSD"`
	TxCount                string    `json:"txCount"`
	LiquidityProviderCount string    `json:"liquidityProviderCount"`

	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
}

type ListPairToken struct {
	PairTokens []PairToken `json:"tokens"`
}

func (pairToken *PairToken) CallInsert() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	fmt.Println(pairToken)

	jsonBody, err := json.Marshal(pairToken)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := urlConnectDatabase + "pairToken"
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
