package toutiao

import (
	"../configs"
	"bytes"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/fedesog/webdriver"
	"github.com/kirinlabs/HttpRequest"
	"net/url"
	"strconv"
	"strings"
)

const (
	URL = "https://www.toutiao.com/api/search/content/?aid=24&app_name=web_search&count=20&format=json&autoload=true"
)

/*
* 获取头条cookie
 */
func GetCookie() string {
	cp := configs.ChromePath()
	chrome := webdriver.NewChromeDriver(cp)
	if err := chrome.Start(); err != nil {
		panic(err)
	}
	desired := webdriver.Capabilities{"Platform": configs.Platform()}
	required := webdriver.Capabilities{}
	session, err := chrome.NewSession(desired, required)
	if err != nil {
		panic(err)
	}

	err = session.Url("https://www.toutiao.com")
	if err != nil {
		panic(err)
	}

	cookies, err := session.GetCookies()
	if err != nil {
		panic(err)
	}

	var cookieString string
	for _, cookie := range cookies {
		cookieString += strings.Join([]string{cookie.Name, "=", cookie.Value, ";"}, "")
	}

	err = session.Delete()
	if err != nil {
		panic(err)
	}

	err = chrome.Stop()
	if err != nil {
		panic(err)
	}
	return cookieString
}

func RequestParams(keyword string, offset int) string {
	var _url string
	params := url.Values{}
	params.Add("keyword", keyword)
	params.Add("offset", strconv.Itoa(offset))
	body := params.Encode()
	_url = URL + "&" + body
	return _url
}

func GetHttpResponse(keyword string, offset int, cookie string) string {
	_url := RequestParams(keyword, offset)
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		"server":       "Tengine",
		"referer":      "https://www.toutiao.com/search/?keyword=%E5%85%AC%E5%8B%9F%E5%9F%BA%E9%87%91",
		"cookie":       cookie,
	})
	res, _ := req.Get(_url)
	body, err := res.Export()
	if err != nil {
		panic(err)
	}
	return string(body)
}

type HtmlBody struct {
	More   int
	Offset int
	Data   *simplejson.Json
}

type News struct {
	Title    string
	Abstract string
	Url      string
	Source   string
	Savedate string
	Keyword  string
}

func (h *HtmlBody) FormatJson(json string) {
	buf := bytes.NewBuffer([]byte(json))
	data, _ := simplejson.NewFromReader(buf)
	h.More = data.Get("has_more").MustInt()
	h.Offset = data.Get("offset").MustInt()
	h.Data = data.Get("data")
}

func (h *HtmlBody) Info() []News {
	if len(h.Data.MustArray()) == 0 {
		fmt.Println("\t今日头条cookie失效，可按Ctrl C退出程序，在conf.json文件中重新设置cookie")
	}
	var allNews []News
	for i := 0; i < 20; i++ {
		news := h.Data.GetIndex(i)
		date := news.Get("datetime").MustString()
		if len(date) == 19 {
			title := news.Get("display").Get("title").Get("text").MustString()
			abstract := news.Get("abstract").MustString()
			_url := news.Get("article_url").MustString()
			source := news.Get("media_name").MustString()
			keyword := news.Get("keyword").MustString()
			n := News{Title: title, Abstract: abstract, Url: _url, Source: source, Savedate: date, Keyword: keyword}
			allNews = append(allNews, n)
		}
	}
	return allNews
}

func SimpleRunner(keyword string, cookie string) []News {
	var allNews []News

	for offset := 0; offset < 180; offset += 20 {
		data := GetHttpResponse(keyword, offset, cookie)
		var body HtmlBody
		body.FormatJson(data)
		news := body.Info()
		allNews = append(allNews, news...)
	}
	return allNews
}
