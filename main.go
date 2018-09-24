package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("%s [options] search query\n", os.Args[0])
		flag.PrintDefaults()
	}
	pages := flag.Int("pages", 4, "Number of pages (25 per page)")
	minPrice := flag.Int("min", 0, "Minimum price")
	maxPrice := flag.Int("max", 0, "Maximum price")
	matchScore := flag.Float64("score", 0.75, "Required Keyword Score")
	help := flag.Bool("help", false, "Show this message")
	flag.Parse()
	search := strings.Join(flag.Args(), " ")
	if flag.NArg() == 0 || *help {
		flag.Usage()
		return
	}
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Jar: cookieJar,
	}
	_, err = doRequest(client, "GET", "https://www.wish.com", nil)
	if err != nil {
		panic(err)
	}
	start := 0
	for {
		if start >= *pages {
			break
		}
		urlvalues := url.Values{}
		urlvalues.Set("count", "25")
		urlvalues.Set("only_wish_express", "false")
		urlvalues.Set("start", strconv.Itoa(25*start))
		urlvalues.Set("transform", "true")
		urlvalues.Set("query", search)
		resp, err := doRequest(client, "POST", "https://www.wish.com/api/search", urlvalues)
		if err != nil {
			panic(err)
		}
		if results, ok := resp["results"].([]interface{}); ok {
			for _, result := range results {
				if info, ok := result.(map[string]interface{})["commerce_product_info"].(map[string]interface{}); ok {
					if variations, ok := info["variations"].([]interface{}); ok {
						for _, variation := range variations {
							price := variation.(map[string]interface{})["localized_price"].(map[string]interface{})["localized_value"].(float64)
							shipping := variation.(map[string]interface{})["localized_shipping"].(map[string]interface{})["localized_value"].(float64)
							total := price + shipping
							if int(total) >= *minPrice && (int(total) <= *maxPrice || *maxPrice == 0) {
								v := url.Values{}
								v.Set("cid", info["id"].(string))
								v.Set("do_not_track", "true")
								v.Set("request_sizing_chart_info", "true")
								product, err := doRequest(client, "POST", "https://www.wish.com/api/product/get", v)
								if err != nil {
									panic(err)
								}
								product = product["contest"].(map[string]interface{})
								name := product["name"].(string)
								if checkScore(name, search, *matchScore) {
									fmt.Printf("%s\n\tURL:%s\n\tPrice:$%v+$%v=$%v\n", name, result.(map[string]interface{})["product_url"], price, shipping, total)
								}
							}
						}
					}
				}
			}
		}
		start++
	}
}

func checkScore(name, search string, required float64) bool {
	splt := strings.Split(strings.ToLower(search), " ")
	matched := 0
	for _, word := range splt {
		if strings.Contains(strings.ToLower(name), word) {
			matched++
		}
		if score := float64(matched) / float64(len(splt)); score >= required {
			return true
		}
	}
	return false
}

func doRequest(client *http.Client, method, url string, data url.Values) (map[string]interface{}, error) {
	var postData io.Reader
	if data != nil {
		postData = strings.NewReader(data.Encode())
	}
	req, err := http.NewRequest(method, url, postData)
	if err != nil {
		return nil, err
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for _, c := range client.Jar.Cookies(req.URL) {
		if c.Name == "_xsrf" {
			req.Header.Set("X-XSRFToken", c.Value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	m := map[string]interface{}{}
	if resp.Header.Get("Content-Type") != "application/json" {
		return m, nil
	}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	if int(m["code"].(float64)) != 0 {
		return nil, fmt.Errorf(m["msg"].(string))
	}
	return m["data"].(map[string]interface{}), nil
}
