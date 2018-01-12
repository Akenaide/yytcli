package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const yuyuteiURL = "http://yuyu-tei.jp/"
const yuyuteiBase = "http://yuyu-tei.jp/game_ws"

type card struct {
	ID          string
	Translation string
	Amount      int
	URL         string
	Price       int
}

func GetCards(out *os.File) {
	fmt.Println("getcards")
	var buffer bytes.Buffer
	cardMap := map[string]card{}
	filter := "ul[data-class=sell] .item_single_card .nav_list_second .nav_list_third a"
	doc, err := goquery.NewDocument(yuyuteiBase)

	if err != nil {
		fmt.Println("Error in get yyt urls")
	}

	doc.Find(filter).Each(func(i int, s *goquery.Selection) {
		url, has := s.Attr("href")
		fmt.Println(url)
		if has {
			images, errCard := goquery.NewDocument(yuyuteiURL + url)
			images.Find(".card_unit").Each(func(cardI int, cardS *goquery.Selection) {
				var price string
				price = cardS.Find(".price .sale").Text()
				if price == "" {
					price = strings.TrimSpace(cardS.Find(".price").Text())
				}
				cardPrice, errAtoi := strconv.Atoi(strings.TrimSuffix(price, "円"))
				if errAtoi != nil {
					fmt.Println(errAtoi)
				}
				cardURL, _ := cardS.Find(".image img").Attr("src")
				cardURL = strings.Replace(cardURL, "90_126", "front", 1)
				yytInfo := card{URL: cardURL, Price: cardPrice}
				cardMap[strings.TrimSpace(cardS.Find(".id").Text())] = yytInfo
			})
			if errCard != nil {
				fmt.Println(errCard)
			}
		}
	})
	b, errMarshal := json.Marshal(cardMap)
	if errMarshal != nil {
		fmt.Println(errMarshal)
	}
	json.Indent(&buffer, b, "", "\t")
	buffer.WriteTo(out)
	fmt.Println("finish")
}
