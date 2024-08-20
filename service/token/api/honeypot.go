package api

import (
	"encoding/json"
	"io"
	"net/http"
	"new-token/pkg/log"
	"new-token/service/token/dao"
	"strings"
)

func CallHoneypot(address string, chainId string) (string, error) {
	url := ""
	if chainId == "1" {
		url = "https://aywt3wreda.execute-api.eu-west-1.amazonaws.com/default/IsHoneypot?chain=eth&token="
	}
	if chainId == "56" {
		url = "https://aywt3wreda.execute-api.eu-west-1.amazonaws.com/default/IsHoneypot?chain=bsc2&token="
	}

	api := url + address
	res, err := http.Get(api)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	isScam := &IsScam{}
	err = json.Unmarshal(body, isScam)
	if err != nil {
		return "", err
	}

	if isScam.IsHoneypot != nil {
		if *isScam.IsHoneypot {
			proof := &dao.Proof{
				IsScam:    isScam.Error,
				IsWarning: nil,
			}
			proofByte, err := json.Marshal(proof)
			if err != nil {
				log.Println(log.LogLevelError, "CallHoneypot : json.Marshal(proof)", err)
			}
			proofByteString := strings.ReplaceAll(string(proofByte), "{'code': -32000, 'message': '", "")
			proofByteString = strings.ReplaceAll(proofByteString, "'}", "")

			return proofByteString, nil

			// query := `update crypto set proof = $$` + proofByteString + `$$ where address = $$` + address + `$$ and chainId = $$` + chainId + `$$;`

			// fmt.Println(query)
		}
	}
	return "", nil
}
