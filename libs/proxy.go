package lib

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// freeProxy Url
const freeProxy = "https://free-proxy-list.net/"
const bin = "https://www.google.com/"

// SkipProxies contains not working proxies rip
var SkipProxies = []string{}

// Proxy handle proxy things
type Proxy struct {
	Info   string
	Client *http.Client
}

func (p *Proxy) Ban() {
	log.Printf("Ban %v", p.Info)
	SkipProxies = append(SkipProxies, p.Info)
}

// GetClient return client with proxy
func GetClient() (*Proxy, error) {
	rand.Seed(time.Now().Unix())
	proxy := Proxy{}
	proxies, err := getProxy()

	if err != nil {
		log.Println("Error in getProxy")
	}

	for {
		blackListed := false
		proxy.Info = proxies[rand.Intn(len(proxies))]

		for _, val := range SkipProxies {
			if proxy.Info == val {
				blackListed = true
				break
			}
		}
		if blackListed {
			continue
		}
		proxyURL, err := url.Parse(fmt.Sprintf("http://%v", proxy.Info))
		if err != nil {
			log.Println("Error in parse")
		}
		proxy.Client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

		_, errHTTP := proxy.Client.Get(bin)
		if errHTTP != nil {
			proxy.Ban()
			continue
		}
		// log.Printf(proxyURL.String())
		break
	}

	return &proxy, nil
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
