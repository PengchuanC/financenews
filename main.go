package main

import (
	"./configs"
	"./database"
	"./eastmoney"
	"./hotnews"
	"./toutiao"
	"fmt"
	"github.com/robfig/cron"
	"time"
)

var c configs.Config

func init() {
	c.ReadConfig()
}

func main() {
	fmt.Println("爬虫任务调度器启动，当前执行周期为@hourly...")
	schedule := cron.New()
	_ = schedule.AddFunc("@hourly", func() {
		// "*/5 * * * * ?"
		go searchEastMoney()
		go searchTouTiao()
		go searchHotNews()
		info := fmt.Sprintf("\t%v 执行爬虫任务", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println(info)
	})
	schedule.Start()
	select {}
}

func searchEastMoney() {
	for _, keyword := range c.KeyWords.EastMoney {
		var n []eastmoney.News
		n = eastmoney.SimpleRunner(keyword)
		dbn := formatEastMoney(n)
		db := getDB()
		db.Connect()
		db.InsertAll(dbn)
		db.Close()
	}
}

func searchTouTiao() {
	cookie := toutiao.GetCookie()
	for _, keyword := range c.KeyWords.TouTiao {
		var n []toutiao.News
		n = toutiao.SimpleRunner(keyword, cookie)
		dbn := formatTouTiao(n)
		db := getDB()
		db.Connect()
		db.InsertAll(dbn)
		db.Close()
	}
}

func searchHotNews() {
	n := hotnews.SimpleRunner()
	dbn := formatHotNews(n)
	db := getDB()
	db.Connect()
	db.InsertAll(dbn)
	db.Close()
}

func getDB() database.DB {
	dbc := c.DataBase
	var db database.DB
	db = database.DB{User: dbc.User, Pass: dbc.Pass, Host: dbc.Host, Port: dbc.Port, Database: dbc.DBname, Instance: nil}
	return db
}

func formatEastMoney(news []eastmoney.News) []database.News {
	var dbNews []database.News
	for _, n := range news {
		var dbn database.News
		dbn.Title = n.Title
		dbn.Abstract = n.Abstract
		dbn.Url = n.Url
		dbn.Source = n.Source
		dbn.Savedate = n.Savedate
		dbn.Keyword = n.Keyword
		dbNews = append(dbNews, dbn)
	}
	return dbNews
}

func formatTouTiao(news []toutiao.News) []database.News {
	var dbNews []database.News
	for _, n := range news {
		var dbn database.News
		dbn.Title = n.Title
		dbn.Abstract = n.Abstract
		dbn.Url = n.Url
		dbn.Source = n.Source
		dbn.Savedate = n.Savedate
		dbn.Keyword = n.Keyword
		dbNews = append(dbNews, dbn)
	}
	return dbNews
}

func formatHotNews(news []hotnews.EastMoney) []database.News {
	var dbNews []database.News
	for _, n := range news {
		var dbn database.News
		dbn.Title = n.Title
		dbn.Abstract = n.Abstract
		dbn.Url = n.Url
		dbn.Source = n.Source
		dbn.Savedate = n.Datetime
		dbn.Keyword = n.Keyword
		dbNews = append(dbNews, dbn)
	}
	return dbNews
}
