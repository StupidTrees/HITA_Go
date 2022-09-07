package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type NsfwResult struct {
	Neutral float64 `json:"Neutral"`
	Drawing float64 `json:"Drawing"`
	Porn    float64 `json:"Porn"`
	Hentai  float64 `json:"Hentai"`
	Sexy    float64 `json:"Sexy"`
}

func CheckSensitive(imagePath string) (bool, error) {
	params := url.Values{}
	parseURL, err := url.Parse("http://localhost:3001/check/image")
	if err != nil {
		log.Println("err")
	}
	params.Set("ipath", imagePath)
	//如果参数中有中文参数,这个方法会进行URLEncode
	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	fmt.Println(string(b))
	r := NsfwResult{}
	err = json.Unmarshal(b, &r)
	return r.Porn > 0.5 || r.Hentai > 0.5 || r.Sexy > 0.5, err
}
