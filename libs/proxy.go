package lib

import (
	"fmt"
	"log"
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

var availableProxies = make(chan *Proxy, 200)
var banProxy = make(chan string)

// Proxy handle proxy things
type Proxy struct {
	Info   string
	Client *http.Client
}

// Readd good proxy
func (p *Proxy) Readd() {
	availableProxies <- p
}

// Ban proxy
func (p *Proxy) Ban() {
	// log.Printf("Ban %v", p.Info)
	banProxy <- p.Info
}

// ProxyStart start channels
func ProxyStart() {
	ticker := time.NewTicker(3 * time.Minute)
	go getProxy()

	go func() {
		for {
			select {
			case skip := <-banProxy:
				SkipProxies = append(SkipProxies, skip)
			case <-ticker.C:
				go getProxy()
			}
		}
	}()
}

// GetClient return client with proxy
func GetClient() *Proxy {
	return <-availableProxies
}

func getProxy() {
	log.Printf("Get new proxies")
	response, errGet := http.Get(freeProxy)
	if errGet != nil {
		fmt.Println("Error on get proxy")
		return
	}

	query, errParse := goquery.NewDocumentFromReader(response.Body)
	if errParse != nil {
		return
	}

	query.Find("table tr").Each(func(_ int, proxyLi *goquery.Selection) {
		if strings.Contains(proxyLi.Text(), "elite proxy") {
			ip := proxyLi.Children().First()
			res := fmt.Sprintf("%v:%v", ip.Text(), ip.Next().Text())

			for _, val := range SkipProxies {
				if res == val {
					return
				}
			}
			go basicTestProxy(res)
		}
	})
}

func basicTestProxy(p string) {
	proxy := Proxy{Info: p}
	proxyURL, err := url.Parse(fmt.Sprintf("http://%v", proxy.Info))
	if err != nil {
		log.Println("Error in parse url")
	}
	proxy.Client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	_, errHTTP := proxy.Client.Get(bin)
	if errHTTP != nil {
		proxy.Ban()
		return
	}
	// log.Println("Good proxy", proxy.Info)
	availableProxies <- &proxy
}
