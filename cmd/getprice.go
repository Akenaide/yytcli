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
	"fmt"

	mylib "github.com/Akenaide/yytcli/libs"
	"github.com/spf13/cobra"
)

var goodRarity = []string{"RR", "R", "U", "C", "CC", "CR"}

// getpriceCmd represents the getprice command
var getpriceCmd = &cobra.Command{
	Use:   "getprice",
	Short: "get playset price (RR, R, U, C, CC, CR)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getprice called")
		for _, serie := range Series {
			var price = 0
			cardMap := mylib.GetCards([]string{serie}, Kizu)
			for _, infos := range cardMap {
				for _, rarity := range goodRarity {
					if rarity == infos.Rarity {
						price = price + (infos.Price * 4)
					}
				}
			}
			fmt.Printf("%v: %v\n", serie, price)
		}
	},
}

func init() {
	rootCmd.AddCommand(getpriceCmd)
}
