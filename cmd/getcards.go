// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var Series []string
var Output string
var Rotate bool

const yuyuteiURL = "http://yuyu-tei.jp/"
const yuyuteiBase = "http://yuyu-tei.jp/game_ws"

type card struct {
	ID          string
	Translation string
	Amount      int
	URL         string
	Price       int
}

// getcardsCmd represents the getcards command
var getcardsCmd = &cobra.Command{
	Use:   "getcards",
	Short: "get card infos from yyt",
	Long: `get :
	* card Id
	* card price
	* card image `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getcards")
		if Rotate {
			var ext = filepath.Ext(Output)
			var datestamp = strings.Split(time.Now().Format(time.RFC3339), "T")[0]
			Output = fmt.Sprint(strings.TrimSuffix(Output, ext), "-", datestamp, ext)
		}
		out, err := os.Create(Output)
		var buffer bytes.Buffer
		defer out.Close()
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
	},
}

func init() {
	rootCmd.AddCommand(getcardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getcardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	var outputDefault = "yyt_infos.json"
	getcardsCmd.Flags().StringArrayVarP(&Series, "series", "s", []string{}, "Default fetch all series")
	getcardsCmd.Flags().StringVarP(&Output, "output", "o", outputDefault, "Specify output path")
	getcardsCmd.Flags().BoolVarP(&Rotate, "rotate", "r", false, "Do you want to rotate (add datetime) to output filename")
}
