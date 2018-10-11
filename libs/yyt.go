package lib

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const maxLength = 2
const yuyuteiURL = "https://yuyu-tei.jp"
const yuyuteiBase = "https://yuyu-tei.jp/game_ws"
const yuyuteiPart = "https://yuyu-tei.jp/game_ws/sell/sell_price.php?ver="

type waitCard struct {
	Card Card
	Wg   *sync.WaitGroup
}

type Card struct {
	ID          string
	Translation string
	Amount      int
	URL         string
	Price       int
	YytSetCode  string
	Rarity      string
	CardURL     string
	EBFoil      bool
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

	imageURL, _ := cardLi.Find(".image img").Attr("src")
	imageURL = fmt.Sprintf("%v%v", yuyuteiURL, strings.Replace(imageURL, "90_126", "front", 1))
	yytInfo := Card{
		URL:        imageURL,
		Price:      cardPrice,
		ID:         cardID,
		YytSetCode: yytSetCode,
		Rarity:     rarity,
		CardURL:    cardURL,
	}
	return yytInfo
}

func fetchCards(url string, tmpCardChan chan waitCard) {
	var wg sync.WaitGroup
	fmt.Println(url)
	images, errCard := goquery.NewDocument(url)
	yytSetCode := images.Find("input[name='item[ver]']").AttrOr("value", "undefined")
	// fetch normal cards
	images.Find("li:not([class*=rarity_S-]).card_unit").Each(func(cardI int, cardLi *goquery.Selection) {
		wg.Add(1)
		v := buildMap(cardLi, yytSetCode)
		tmpCardChan <- waitCard{Card: v, Wg: &wg}
	})
	// fetch EB foil cards
	images.Find("li[class*=rarity_S-]").Each(func(cardI int, cardLi *goquery.Selection) {
		wg.Add(1)
		v := buildMap(cardLi, yytSetCode)
		v.EBFoil = true
		tmpCardChan <- waitCard{Card: v, Wg: &wg}
	})
	if errCard != nil {
		fmt.Println(errCard)
	}
	wg.Wait()
}

// GetCards function
func GetCards(series []string) map[string]Card {
	// Get cards
	fmt.Println("getcards")
	fetchChannel := make(chan bool, maxLength)
	tmpCardChan := make(chan waitCard, 10)
	cardMap := map[string]Card{}
	var wg sync.WaitGroup

	go func() {
		for {
			select {
			case waitCardS := <-tmpCardChan:
				var cardID = ""
				if waitCardS.Card.EBFoil {
					cardID = waitCardS.Card.ID + "F"
				} else {
					cardID = waitCardS.Card.ID
				}
				cardMap[cardID] = waitCardS.Card
				waitCardS.Wg.Done()
			}
		}
	}()

	if len(series) == 0 {
		filter := "ul[data-class=sell] .item_single_card .nav_list_second .nav_list_third a"
		doc, err := goquery.NewDocument(yuyuteiBase)

		if err != nil {
			fmt.Println("Error in get yyt urls")
		}
		doc.Find(filter).Each(func(i int, s *goquery.Selection) {
			url, has := s.Attr("href")
			if has {
				fetchChannel <- true
				wg.Add(1)
				go func(url string) {
					fetchCards(strings.Join([]string{yuyuteiURL, url}, ""), tmpCardChan)
					<-fetchChannel
					wg.Done()
				}(url)
			}
		})
	} else {
		for _, url := range series {
			wg.Add(1)
			fetchChannel <- true
			go func(url string) {
				fetchCards(strings.Join([]string{yuyuteiPart, url}, ""), tmpCardChan)
				<-fetchChannel
				wg.Done()
			}(url)
		}
	}

	wg.Wait()
	return cardMap
}
