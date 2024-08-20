package api

import (
	"encoding/json"
	"io"
	"net/http"
	"new-token/pkg/log"
	"new-token/service/token/dao"
)

func CallStaySafu(crypto dao.Crypto) (dao.Crypto, error) {

	cryptoStaySafu := dao.Crypto{
		Proof: crypto.Proof,
	}

	api := "https://api.staysafu.org/api/freescan?tokenAddress=" + crypto.Address
	res, err := http.Get(api)
	if err != nil {
		log.Println(log.LogLevelDebug, "CallStaySafu : http.Get(api)"+api, err.Error())
		return cryptoStaySafu, err
	}
	// fmt.Println("scam cux", *crypto.Proof.IsScam)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return cryptoStaySafu, err
	}
	defer res.Body.Close()

	staySafu := &StaySafu{}
	err = json.Unmarshal(body, staySafu)
	if err != nil {
		return cryptoStaySafu, err
	}

	if staySafu.Result.IsToken != nil {
		proofIsScam := ""
		if crypto.Proof.IsScam != nil {
			proofIsScam = *crypto.Proof.IsScam
			for _, warn := range staySafu.Result.Warnings {
				proofIsScam += ("\n" + warn.Message)
			}
		} else {
			for _, warn := range staySafu.Result.Warnings {
				proofIsScam += (warn.Message + "\n")
			}
		}
		proofIsScam += "\nContract Self Destruct"
		crypto.Proof.IsScam = &proofIsScam

		//update proof for cryptoStaySafu and return
		cryptoStaySafu.Proof = crypto.Proof

	} else {
		proofIsScam := ""
		proofIsWarning := ""

		//todo: check token is scam
		//if scam proof exist -> append warning proof to scam proof
		//else -> append warning proof to warning proof

		if crypto.Proof.IsScam != nil {
			proofIsScam = *crypto.Proof.IsScam
			for _, warn := range staySafu.Result.Warnings {
				proofIsScam += ("\n" + warn.Message)
			}
		} else {
			lenWarning := len(staySafu.Result.Warnings)
			for index := 0; index < lenWarning; index++ {
				if index != lenWarning-1 {
					proofIsWarning += (staySafu.Result.Warnings[index].Message + "\n")
				} else {
					proofIsWarning += (staySafu.Result.Warnings[index].Message)
				}
			}
		}

		if len(proofIsScam) != 0 {
			crypto.Proof.IsScam = &proofIsScam
		} else {
			crypto.Proof.IsScam = nil
		}

		if len(proofIsWarning) != 0 {
			crypto.Proof.IsWarning = &proofIsWarning
		} else {
			crypto.Proof.IsWarning = nil
		}

		//update proof for cryptoStaySafu and return
		cryptoStaySafu.Proof = crypto.Proof
	}
	return cryptoStaySafu, nil
}
