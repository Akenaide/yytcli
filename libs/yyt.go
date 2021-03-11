package lib

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/Akenaide/biri"
	"github.com/PuerkitoBio/goquery"
)

const maxLength = 15
const yuyuteiURL = "https://yuyu-tei.jp"
const yuyuteiBase = "https://yuyu-tei.jp/game_ws"
const yuyuteiPart = "https://yuyu-tei.jp/game_ws/sell/sell_price.php?ver="
const yytMenu = "kana"

var wg sync.WaitGroup

type Card struct {
	Amount      int
	CardURL     string
	EBFoil      bool
	ID          string
	Price       int
	Rarity      string
	Stock       int
	Translation string
	URL         string
	YytSetCode  string
}

func getStock(stockString string) (int, error) {
	stock := strings.Split(stockString, "：")[1]
	switch stock {
	case "×":
		return 0, nil
	case "◯":
		return 99, nil
	default:
		return strconv.Atoi(stock)
	}
}

func buildMap(cardLi *goquery.Selection, yytSetCode string) Card {
	var price string
	cardLi.ToggleClass("card_unit")
	rarity := strings.TrimLeft(cardLi.AttrOr("class", "rarity_Unknow"), "rarity_")
	cardID := strings.TrimSpace(cardLi.Find(".id").Text())
	price = cardLi.Find(".price .sale").Text()
	cardURL := fmt.Sprintf("%v%s", yuyuteiURL, cardLi.Find(".image_box a").AttrOr("href", "undefined"))

	if price == "" {
		price = strings.TrimSpace(cardLi.Find(".price").Text())
	}
	cardPrice, errAtoi := strconv.Atoi(strings.TrimSuffix(price, "円"))

	if errAtoi != nil {
		fmt.Println(errAtoi)
	}

	stock, stockErr := getStock(cardLi.Find(".stock").Text())

	if stockErr != nil {
		fmt.Println(stockErr)
	}

	imageURL, _ := cardLi.Find(".image img").Attr("src")
	imageURL = strings.Replace(imageURL, "90_126", "front", 1)
	yytInfo := Card{
		URL:        imageURL,
		Price:      cardPrice,
		ID:         cardID,
		YytSetCode: yytSetCode,
		Rarity:     rarity,
		CardURL:    cardURL,
		Stock:      stock,
	}
	return yytInfo
}

func fetchCards(url string, tmpCardChan chan Card) {
	var images *goquery.Document
	var errCard error
	var yytSetCode string
	numberOfCard := 0

	for {
		log.Println(" -- Begin", url)
		proxy := biri.GetClient()
		resp, errHTTP := proxy.Client.Get(url)
		if errHTTP != nil {
			// log.Printf(errHTTP.Error())
			proxy.Ban()
			continue
		}
		images, errCard = goquery.NewDocumentFromResponse(resp)
		if errCard != nil {
			// log.Println(errCard)
			// log.Printf("RETRY for %v", url)
			proxy.Ban()
			continue
		}

		yytSetCode = images.Find("input[name='item[ver]']").AttrOr("value", "undefined")
		if yytSetCode == "undefined" {
			continue
		}

		proxy.Readd()
		break
	}

	// log.Println(url)
	// fetch normal cards
	images.Find("li:not([class*=rarity_S-]).card_unit").Each(func(cardI int, cardLi *goquery.Selection) {
		wg.Add(1)
		v := buildMap(cardLi, yytSetCode)
		tmpCardChan <- v
		numberOfCard = numberOfCard + 1
	})
	// fetch EB foil cards
	images.Find("li[class*=rarity_S-]").Each(func(cardI int, cardLi *goquery.Selection) {
		wg.Add(1)
		v := buildMap(cardLi, yytSetCode)
		v.EBFoil = true
		tmpCardChan <- v
		numberOfCard = numberOfCard + 1
	})
	log.Println(numberOfCard, " for set: ", yytSetCode)
}

// GetCards function
func GetCards(series []string, kizu bool) map[string]Card {
	// Get cards
	fmt.Println("getcards")
	fetchChannel := make(chan bool, maxLength)
	tmpCardChan := make(chan Card, 10)
	cardMap := map[string]Card{}

	biri.ProxyStart()

	go func() {
		for {
			select {
			case waitCardS := <-tmpCardChan:
				var cardID = ""
				if waitCardS.EBFoil {
					cardID = waitCardS.ID + "F"
				} else {
					cardID = waitCardS.ID
				}
				cardMap[cardID] = waitCardS
				// log.Println(waitCardS.ID)
				wg.Done()
			}
		}
	}()

	if len(series) == 0 {
		filter := "ul[data-class=sell] .item_single_card .nav_list_second a"
		var doc *goquery.Document

		numberOfSet := 0

		for {
			proxy := biri.GetClient()

			resp, errHTTP := proxy.Client.Get(yuyuteiBase)
			if errHTTP != nil {
				continue
			}

			doc, _ = goquery.NewDocumentFromResponse(resp)
			break
		}

		doc.Find(filter).Each(func(i int, s *goquery.Selection) {
			urlSet, has := s.Attr("href")
			if kizu {
				urlSet = urlSet + "&kizu=1"
			}

			if has {
				urlParsed, _ := url.Parse(urlSet)
				if urlParsed.Query().Get("menu") != yytMenu {
					return
				}

				numberOfSet = numberOfSet + 1
				log.Printf("Nb %v / %v", numberOfSet, urlSet)
				fetchChannel <- true
				wg.Add(1)

				go func(url string) {
					fetchCards(strings.Join([]string{yuyuteiURL, url}, ""), tmpCardChan)
					<-fetchChannel
					wg.Done()
				}(urlSet)
			}
		})
	} else {
		for _, url := range series {
			if kizu {
				url = url + "&kizu=1"
			}
			wg.Add(1)
			fetchChannel <- true
			go func(url string) {
				fetchCards(strings.Join([]string{yuyuteiPart, url}, ""), tmpCardChan)
				<-fetchChannel
				wg.Done()
			}(url)
		}
	}

	log.Printf("Wait")
	wg.Wait()
	log.Printf("Finish get")
	return cardMap
}
