package lib

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const yuyuteiURL = "http://yuyu-tei.jp"
const yuyuteiBase = "http://yuyu-tei.jp/game_ws"
const yuyuteiPart = "https://yuyu-tei.jp/game_ws/sell/sell_price.php?ver="

type Card struct {
	ID          string
	Translation string
	Amount      int
	URL         string
	Price       int
	YytSetCode  string
	Rarity      string
}

func buildMap(cardLi *goquery.Selection, yytSetCode string) (string, Card) {
	var price string
	cardLi.ToggleClass("card_unit")
	rarity := strings.TrimLeft(cardLi.AttrOr("class", "rarity_Unknow"), "rarity_")
	cardID := strings.TrimSpace(cardLi.Find(".id").Text())
	price = cardLi.Find(".price .sale").Text()
	if price == "" {
		price = strings.TrimSpace(cardLi.Find(".price").Text())
	}
	cardPrice, errAtoi := strconv.Atoi(strings.TrimSuffix(price, "円"))
	if errAtoi != nil {
		fmt.Println(errAtoi)
	}
	cardURL, _ := cardLi.Find(".image img").Attr("src")
	cardURL = fmt.Sprintf("%v%v", yuyuteiURL, strings.Replace(cardURL, "90_126", "front", 1))
	yytInfo := Card{
		URL:        cardURL,
		Price:      cardPrice,
		ID:         cardID,
		YytSetCode: yytSetCode,
		Rarity:     rarity,
	}
	return cardID, yytInfo
}

func fetchCards(url string, cardMap map[string]Card) map[string]Card {
	fmt.Println(url)
	images, errCard := goquery.NewDocument(url)
	yytSetCode := images.Find("input[name='item[ver]']").AttrOr("value", "undefined")
	// fetch normal cards
	images.Find("li:not([class*=rarity_S-]).card_unit").Each(func(cardI int, cardLi *goquery.Selection) {
		k, v := buildMap(cardLi, yytSetCode)
		cardMap[k] = v
	})
	// fetch EB foil cards
	images.Find("li[class*=rarity_S-]").Each(func(cardI int, cardLi *goquery.Selection) {
		k, v := buildMap(cardLi, yytSetCode)
		cardMap[strings.Join([]string{k, "F"}, "")] = v
	})
	if errCard != nil {
		fmt.Println(errCard)
	}
	return cardMap
}

// GetCards function
func GetCards(series []string) map[string]Card {
	// Get cards
	fmt.Println("getcards")
	cardMap := map[string]Card{}
	if len(series) == 0 {
		filter := "ul[data-class=sell] .item_single_card .nav_list_second .nav_list_third a"
		doc, err := goquery.NewDocument(yuyuteiBase)

		if err != nil {
			fmt.Println("Error in get yyt urls")
		}
		doc.Find(filter).Each(func(i int, s *goquery.Selection) {
			url, has := s.Attr("href")
			if has {
				fetchCards(strings.Join([]string{yuyuteiURL, url}, ""), cardMap)
			}
		})
	} else {
		for _, serie := range series {
			fetchCards(strings.Join([]string{yuyuteiPart, serie}, ""), cardMap)
		}
	}

	return cardMap
}
