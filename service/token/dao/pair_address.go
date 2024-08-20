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

type PairAddress struct {
	Id        uuid.UUID `json:"id"`
	ChainId   string    `json:"chainId"`
	ChainName string    `json:"chainName"`
	Defi      string    `json:"defi"`
	Index     int       `json:"index"`
	Address   string    `json:"address"`
}

type ListPairAddress struct {
	PairAddresses []PairAddress `json:"pairAddresses"`
}

func (pairAddress *PairAddress) CallInsert() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	jsonBody, err := json.Marshal(pairAddress)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := urlConnectDatabase + "pairAddress"
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
