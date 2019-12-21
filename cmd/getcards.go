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
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	mylib "github.com/Akenaide/yytcli/libs"
	"github.com/spf13/cobra"
)

var Output string
var Rotate bool

// getcardsCmd represents the getcards command
var getcardsCmd = &cobra.Command{
	Use:   "getcards",
	Short: "get card infos from yyt",
	Long: `get :
	* card Id
	* card price
	* card image `,
	Run: func(cmd *cobra.Command, args []string) {
		var buffer bytes.Buffer
		if Rotate {
			var ext = filepath.Ext(Output)
			var datestamp = strings.Split(time.Now().Format(time.RFC3339), "T")[0]
			Output = fmt.Sprint(strings.TrimSuffix(Output, ext), "-", datestamp, ext)
		}
		out, err := os.Create(Output)
		if err != nil {
			fmt.Printf(err.Error())
		}
		defer out.Close()
		cardMap := mylib.GetCards(Series, Kizu)
		b, errMarshal := json.Marshal(cardMap)
		if errMarshal != nil {
			fmt.Println(errMarshal)
		}
		json.Indent(&buffer, b, "", "\t")
		buffer.WriteTo(out)
		log.Println("finish")
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
	getcardsCmd.Flags().StringVarP(&Output, "output", "o", outputDefault, "Specify output path")
	getcardsCmd.Flags().BoolVarP(&Rotate, "rotate", "r", false, "Do you want to rotate (add datetime) to output filename")
}
