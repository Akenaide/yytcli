package lib

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const maxLength = 2
const yuyuteiURL = "https://yuyu-tei.jp"
const yuyuteiBase = "https://yuyu-tei.jp/game_ws"
const yuyuteiPart = "https://yuyu-tei.jp/game_ws/sell/sell_price.php?ver="

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

func fetchCards(url string, tmpCardChan chan Card) {
	fmt.Println(url)
	images, errCard := goquery.NewDocument(url)
	yytSetCode := images.Find("input[name='item[ver]']").AttrOr("value", "undefined")
	// fetch normal cards
	images.Find("li:not([class*=rarity_S-]).card_unit").Each(func(cardI int, cardLi *goquery.Selection) {
		v := buildMap(cardLi, yytSetCode)
		tmpCardChan <- v
	})
	// fetch EB foil cards
	images.Find("li[class*=rarity_S-]").Each(func(cardI int, cardLi *goquery.Selection) {
		v := buildMap(cardLi, yytSetCode)
		v.EBFoil = true
		tmpCardChan <- v
	})
	if errCard != nil {
		fmt.Println(errCard)
	}
}

// GetCards function
func GetCards(series []string) map[string]Card {
	// Get cards
	fmt.Println("getcards")
	fetchChannel := make(chan bool, maxLength)
	tmpCardChan := make(chan Card, 10)
	cardMap := map[string]Card{}
	go func() {
		for {
			select {
			case card := <-tmpCardChan:
				var cardID = ""
				if card.EBFoil {
					cardID = card.ID + "F"
				} else {
					cardID = card.ID
				}
				cardMap[cardID] = card
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
				go func(url string) {
					fetchCards(strings.Join([]string{yuyuteiURL, url}, ""), tmpCardChan)
					<-fetchChannel
				}(url)
			}
		})
	} else {
		for _, url := range series {
			fetchChannel <- true
			go func(url string) {
				fetchCards(strings.Join([]string{yuyuteiPart, url}, ""), tmpCardChan)
				<-fetchChannel
			}(url)
		}
	}

	return cardMap
}
