package html

import (
	"fmt"
	"io"
	"net/http"
	"new-token/pkg/log"
	"new-token/service/token/constant"
	"new-token/service/token/dao"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
Split dom html scan website
Return holders, logo, website, socials
*/
func SplitHTML(address string, chainId string) (dao.Crypto, error) {
	crypto := dao.Crypto{}
	mapSocials := make(map[string]string)
	failCount := 10

	for {
		client := &http.Client{}

		_, _, url := constant.GetExplorerByChainId(chainId)
		//todo: switch url explorer by chainIn

		// req, err := http.NewRequest("GET", fmt.Sprintf("https://etherscan.io/token/%s", address), nil)
		req, err := http.NewRequest("GET", fmt.Sprintf(url, address), nil)
		req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
		req.Header.Add("Content-Type", "application/json")

		crypto.Explorer = fmt.Sprintf(url, address)

		if err != nil {
			log.Println(log.LogLevelError, "SplitHTML : http.NewRequestl "+chainId+" "+address, err)
			failCount--
			time.Sleep(3 * time.Second)
			if failCount == 0 {
				return crypto, err
			}
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			log.Println(log.LogLevelError, "SplitHTML : client.Do(req) "+chainId+" "+address, err)
			failCount--
			time.Sleep(3 * time.Second)
			if failCount == 0 {
				return crypto, err
			}
			continue
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(log.LogLevelError, "SplitHTML : io.ReadAll(res.Body) "+chainId+" "+address, err)
			failCount--
			time.Sleep(3 * time.Second)
			if failCount == 0 {
				return crypto, err
			}
			continue
		}

		domString := string(body)

		// if chainId == "250" {

		// 	indexHolders := strings.Index(domString, "Holders:</div>")

		// 	// fmt.Println(domString[indexHolders : indexHolders+100])

		// 	if indexHolders == -1 {
		// 		log.Println(log.LogLevelError, "SplitHTML : not found index of "+chainId+" "+address, "")
		// 		failCount--
		// 		time.Sleep(5 * time.Second)
		// 		if failCount == 0 {
		// 			return crypto, errors.New("not found index of")
		// 		}
		// 		continue
		// 	}

		// 	endIndexHolders := 0
		// 	for i := indexHolders + len("Holders:</div>"); i < len(domString); i++ {
		// 		if string(domString[i]) == "/" {
		// 			endIndexHolders = i
		// 			break
		// 		}
		// 	}

		// 	holderHTML := domString[indexHolders+len("Holders:</div>") : endIndexHolders]
		// 	tmp := strings.Split(holderHTML, ">")
		// 	if len(tmp) > 1 {
		// 		holderHTML = strings.Split(holderHTML, ">")[1]
		// 	} else {
		// 		log.Println(log.LogLevelError, "SplitHTML : valid html "+chainId+" "+address, "")
		// 		failCount--
		// 		time.Sleep(5 * time.Second)
		// 		if failCount == 0 {
		// 			return crypto, errors.New("valid html")
		// 		}
		// 		continue
		// 	}

		// 	tmp = strings.Split(holderHTML, "<")
		// 	if len(tmp) > 0 {
		// 		holderHTML = strings.Split(holderHTML, "<")[0]
		// 	} else {
		// 		log.Println(log.LogLevelError, "SplitHTML : valid html "+chainId+" "+address, "")
		// 		failCount--
		// 		time.Sleep(5 * time.Second)
		// 		if failCount == 0 {
		// 			return crypto, errors.New("valid html")
		// 		}
		// 		continue
		// 	}
		// 	holderHTML = regexp.MustCompile(`[^0-9.]+`).ReplaceAllString(holderHTML, "")
		// 	crypto.Holders = holderHTML

		// } else {

		indexHolders := strings.Index(domString, `number of holders`)

		if len(domString) > indexHolders+len("number of holders")+16 {
			holderHTML := domString[indexHolders+len("number of holders") : indexHolders+len("number of holders")+15]

			holderHTML = regexp.MustCompile(`[^0-9.]+`).ReplaceAllString(holderHTML, "")

			intVar, err := strconv.ParseInt(holderHTML, 0, 64)
			if err != nil {
				return crypto, err
			}
			crypto.Holders = intVar
		}

		// }
		//************************* image ***********************************

		startIndexOfImage := strings.Index(domString, "litAssetLogo")
		// fmt.Println("indexOfImage", domString[startIndexOfImage:startIndexOfImage+100])
		endIndexOfImage := startIndexOfImage
		if startIndexOfImage > 0 {
			for i := startIndexOfImage; i < len(domString); i++ {
				if string(domString[i]) == ";" {
					endIndexOfImage = i
					break
				}
			}
			imgHTML := domString[startIndexOfImage:endIndexOfImage]

			tmp := strings.Split(imgHTML, "\"")
			if len(tmp) > 2 {
				imgHTML = tmp[1]
			}

			crypto.BigLogo = imgHTML
		}

		//********************** office site ********************************

		indexOfOfficeSite := strings.Index(domString, "fa fa-external-link-alt small ml-1")
		if indexOfOfficeSite > 170 {
			officeSiteHTML := domString[indexOfOfficeSite-170 : indexOfOfficeSite]

			tmp := strings.Split(officeSiteHTML, "href='")
			if len(tmp) > 1 {
				officeSiteHTML = tmp[1]
			}
			tmp = strings.Split(officeSiteHTML, "'")
			if len(tmp) > 0 {
				officeSiteHTML = tmp[0]
			}
			crypto.Website = officeSiteHTML
		}

		//********************** socials ********************************

		indexOfReddit := strings.Index(domString, "'Reddit: http")
		if indexOfReddit > -1 {
			redditHTML := domString[indexOfReddit : indexOfReddit+100]

			tmp := strings.Split(redditHTML, "'Reddit: ")
			if len(tmp) > 1 {
				redditHTML = tmp[1]
			}

			tmp = strings.Split(redditHTML, "'")
			if len(tmp) > 0 {
				redditHTML = tmp[0]
			}
			mapSocials["reddit"] = redditHTML
		}

		indexOfTwitter := strings.Index(domString, "'Twitter: http")
		if indexOfTwitter > -1 {
			twitterHTML := domString[indexOfTwitter : indexOfTwitter+100]

			tmp := strings.Split(twitterHTML, "'Twitter: ")
			if len(tmp) > 1 {
				twitterHTML = tmp[1]
			}

			tmp = strings.Split(twitterHTML, "'")
			if len(tmp) > 0 {
				twitterHTML = tmp[0]
			}
			mapSocials["twitter"] = twitterHTML
		}

		indexOfGithub := strings.Index(domString, "'Github: http")
		if indexOfGithub > -1 {
			githubHTML := domString[indexOfGithub : indexOfGithub+100]

			tmp := strings.Split(githubHTML, "'Github: ")
			if len(tmp) > 1 {
				githubHTML = tmp[1]
			}

			tmp = strings.Split(githubHTML, "'")
			if len(tmp) > 0 {
				githubHTML = tmp[0]
			}
			mapSocials["github"] = githubHTML
		}

		indexOfTelegram := strings.Index(domString, "'Telegram: http")
		if indexOfTelegram > -1 {
			telegramHTML := domString[indexOfTelegram : indexOfTelegram+100]

			tmp := strings.Split(telegramHTML, "'Telegram: ")
			if len(tmp) > 1 {
				telegramHTML = tmp[1]
			}

			tmp = strings.Split(telegramHTML, "'")
			if len(tmp) > 0 {
				telegramHTML = tmp[0]
			}
			mapSocials["telegram"] = telegramHTML
		}

		indexOfDiscord := strings.Index(domString, "'Discord: http")
		if indexOfDiscord > -1 {
			discordHTML := domString[indexOfDiscord : indexOfDiscord+100]

			tmp := strings.Split(discordHTML, "'Discord: ")
			if len(tmp) > 1 {
				discordHTML = tmp[1]
			}

			tmp = strings.Split(discordHTML, "'")
			if len(tmp) > 0 {
				discordHTML = tmp[0]
			}
			mapSocials["discord"] = discordHTML
		}

		indexOfCoinmarketcap := strings.Index(domString, "https://coinmarketcap.com/currencies")
		if indexOfCoinmarketcap > -1 {
			coinmarketcapHTML := domString[indexOfCoinmarketcap : indexOfCoinmarketcap+100]

			tmp := strings.Split(coinmarketcapHTML, "'")
			if len(tmp) > 0 {
				coinmarketcapHTML = tmp[0]
			}
			mapSocials["coinmarketcap"] = coinmarketcapHTML
		}

		indexOfCoingecko := strings.Index(domString, "https://www.coingecko.com/en/coins")
		if indexOfCoingecko > -1 {
			coingeckoHTML := domString[indexOfCoingecko : indexOfCoingecko+100]

			tmp := strings.Split(coingeckoHTML, "'")
			if len(tmp) > 0 {
				coingeckoHTML = tmp[0]
			}
			mapSocials["coingecko"] = coingeckoHTML
		}

		indexOfBlog := strings.Index(domString, "Blog: http")
		if indexOfBlog > -1 {
			blogHTML := domString[indexOfBlog : indexOfBlog+100]

			tmp := strings.Split(blogHTML, "Blog: ")
			if len(tmp) > 1 {
				blogHTML = tmp[1]
			}

			tmp = strings.Split(blogHTML, "'")
			if len(tmp) > 0 {
				blogHTML = tmp[0]
			}
			mapSocials["blog"] = blogHTML
		}

		indexOfWhitepaper := strings.Index(domString, "Whitepaper: ")
		if indexOfWhitepaper > -1 {
			whitepaperHTML := domString[indexOfWhitepaper : indexOfWhitepaper+100]

			tmp := strings.Split(whitepaperHTML, "Whitepaper: ")
			if len(tmp) > 1 {
				whitepaperHTML = tmp[1]
			}

			tmp = strings.Split(whitepaperHTML, "'")
			if len(tmp) > 0 {
				whitepaperHTML = tmp[0]
			}
			mapSocials["whitepaper"] = whitepaperHTML
		}

		crypto.Socials = mapSocials

		return crypto, nil
	}
}
