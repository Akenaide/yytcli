package lib

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// freeProxy Url
const freeProxy = "https://free-proxy-list.net/"

// GetClient return client with proxy
func GetClient() {

	// proxyUrl, err := url.Parse("http://proxyIp:proxyPort")
	// myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

}

func getProxy() ([]string, error) {
	response, errGet := http.Get(freeProxy)
	var res []string
	if errGet != nil {
		fmt.Println("Error on get proxy")
		return nil, errGet
	}

	query, errParse := goquery.NewDocumentFromReader(response.Body)
	if errParse != nil {
		return nil, errParse
	}

	query.Find("table tr").Each(func(_ int, proxyLi *goquery.Selection) {
		if strings.Contains(proxyLi.Text(), "elite proxy") {
			ip := proxyLi.Children().First()
			res = append(res, fmt.Sprintf("%v:%v", ip.Text(), ip.Next().Text()))
		}
	})

	return res, nil
}
