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

type FailRepo struct {
	Fails []CrawlFail `json:"failRepo"`
}

type CrawlFail struct {
	Id          uuid.UUID `json:"id"`
	ChainId     string    `json:"chainId"`
	Defi        string    `json:"defi"`
	Address     string    `json:"address"`
	Index       int       `json:"index"`
	Type        string    `json:"type"`
	Errors      string    `json:"errors"`
	CreatedDate string    `json:"createdDate"`
}

func (crawlFail *CrawlFail) CallInsertFail() error {
	urlConnectDatabase := server.Config.GetString("URL_CONNECT_DATABASE")

	jsonBody, err := json.Marshal(crawlFail)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := urlConnectDatabase + "fail"
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
