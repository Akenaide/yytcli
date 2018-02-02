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

type card struct {
	ID          string
	Translation string
	Amount      int
	URL         string
	Price       int
}

func buildMap(cardS *goquery.Selection) (string, card) {
	var price string
	cardID := strings.TrimSpace(cardS.Find(".id").Text())
	price = cardS.Find(".price .sale").Text()
	if price == "" {
		price = strings.TrimSpace(cardS.Find(".price").Text())
	}
	cardPrice, errAtoi := strconv.Atoi(strings.TrimSuffix(price, "円"))
	if errAtoi != nil {
		fmt.Println(errAtoi)
	}
	cardURL, _ := cardS.Find(".image img").Attr("src")
	cardURL = fmt.Sprintf("%v%v", yuyuteiURL, strings.Replace(cardURL, "90_126", "front", 1))
	yytInfo := card{URL: cardURL, Price: cardPrice, ID: cardID}
	return cardID, yytInfo
}

func fetchCards(url string, cardMap map[string]card) map[string]card {
	fmt.Println(url)
	images, errCard := goquery.NewDocument(url)
	// fetch normal cards
	images.Find("li:not([class*=rarity_S-]).card_unit").Each(func(cardI int, cardS *goquery.Selection) {
		k, v := buildMap(cardS)
		cardMap[k] = v
	})
	// fetch EB foil cards
	images.Find("li[class*=rarity_S-]").Each(func(cardI int, cardS *goquery.Selection) {
		k, v := buildMap(cardS)
		cardMap[strings.Join([]string{k, "F"}, "")] = v
	})
	if errCard != nil {
		fmt.Println(errCard)
	}
	return cardMap
}

// GetCards function
func GetCards(series []string) map[string]card {
	// Get cards
	fmt.Println("getcards")
	cardMap := map[string]card{}
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
