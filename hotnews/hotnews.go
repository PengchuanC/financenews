package hotnews

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
	"time"
)

const (
	URL = "http://finance.eastmoney.com/yaowen.html"
)


// 新闻结构
type EastMoney struct{
	Title string
	Abstract string
	Url string
	Datetime string
	Source string
	Keyword string
}
/*
 * 爬虫部分
*/
func visit(url string) []EastMoney{
	var news []EastMoney
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML("div[id=artitileList1]", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, li *colly.HTMLElement){
			var title string = li.ChildText("p[class=title]")
			var abstract string = li.ChildText("p[class=info]")
			var abstractCheck string = li.ChildAttr("p[class=info]", "title")
			// 东方财富提供两种模式的摘要，选取其中信息较多的一类
			if len(abstract) <= len(abstractCheck) {
				abstract = abstractCheck
			}
			var datetime string = li.ChildText("p[class=time]")
			var url string = li.ChildAttr("a", "href")
			if len(title) > 5 {
				if strings.Index(abstract, "】") != -1 {
					abstract = strings.Split(abstract, "】")[1]
				}
				var em EastMoney = EastMoney{title, abstract, url, datetime, "东方财富", ""}
				em.Strftime()
				news = append(news, em)
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {})

	_ = c.Visit(url)
	return news
}


/*为EastMoney结构体添加功能*/
// 时间格式化功能
func(e *EastMoney) Strftime(){
	var t string = e.Datetime
	t = strings.Replace(t, "月", "-", 1)
	t = strings.Replace(t, "日", "", 1)
	var now string = time.Now().String()
	var year string = now[0: 5]
	t = year + t
	e.Datetime = t
}



func SimpleRunner() []EastMoney{
	return visit(URL)
}