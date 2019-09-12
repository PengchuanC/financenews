package eastmoney

import (
	"bytes"
	"github.com/bitly/go-simplejson"
	"github.com/kirinlabs/HttpRequest"
	"net/url"
	"strconv"
	"strings"
)

const (
	REFERERURL = "http://so.eastmoney.com/CArticle/s?&"
	URL = "http://api.so.eastmoney.com/bussiness/Web/GetSearchList?"
)


func requestRefererUrl(keyword string, page int) string {
	param := url.Values{}
	param.Add("keyword", keyword)
	param.Add("pageindex", strconv.Itoa(page))
	param.Add("searchrange", "16384")
	param.Add("sortfield", "8")
	fmtParam := param.Encode()
	_url := REFERERURL + fmtParam
	return _url
}


func requestUrl(keyword string, page int) string{
	p := url.Values{}
	p.Add("type", "16392")
	p.Add("pageindex", strconv.Itoa(page))
	p.Add("pagesize", "10")
	p.Add("keyword", keyword)
	p.Add("name", "caifuhaowenzhang")
	params := p.Encode()
	_url := URL + params
	return _url
}

func GetHttpResponse(keyword string, page int) string{
	_url := requestUrl(keyword, page)
	refererUrl := requestRefererUrl(keyword, page)
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		"referer":     refererUrl,
	})
	res, _:= req.Get(_url)
	body, err := res.Export()
	if err != nil {
		panic(err)
	}
	return string(body)
}


type HtmlBody struct{
	More int
	Total int
	Keyword string
	Data *simplejson.Json
}

type News struct {
	Title string
	Abstract string
	Url string
	Source string
	Savedate string
	Keyword string
}

func (h *HtmlBody) FormatJson(json string) {
	buf := bytes.NewBuffer([]byte(json))
	data, _ := simplejson.NewFromReader(buf)
	h.More = data.Get("Code").MustInt()
	h.Total = data.Get("TotalPage").MustInt()
	h.Keyword = data.Get("Keyword").MustString()
	h.Data = data.Get("Data")
}

func (h *HtmlBody) Info() []News{
	var allNews []News
	for i :=1; i <= 10; i++ {
		var n News
		n.Title = h.Data.GetIndex(i).Get("Title").MustString()
		n.Abstract = h.Data.GetIndex(i).Get("Content").MustString()
		n.Url = h.Data.GetIndex(i).Get("ArticleUrl").MustString()
		n.Source = h.Data.GetIndex(i).Get("NickName").MustString()
		n.Savedate = h.Data.GetIndex(i).Get("ShowTime").MustString()
		n.Keyword = h.Keyword
		n.TitleFormat()
		allNews = append(allNews, n)
	}
	return allNews
}

func (n *News) TitleFormat() {
	for _, flag := range []string{"<em>", "</em>"} {
		index := strings.Index(n.Title, flag)
		if index != -1{
			n.Title = strings.ReplaceAll(n.Title, flag, "")
		}
	}
}


func SimpleRunner(keyword string) []News{
	var allNews []News
	for page := 1;page <= 5;page++{
		var hb HtmlBody
		text := GetHttpResponse(keyword, page)
		hb.FormatJson(text)
		n := hb.Info()
		allNews = append(allNews, n...)
	}
	return allNews
}