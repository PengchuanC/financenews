package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

type News struct {
	Title string
	Abstract string
	Url string
	Source string
	Savedate string
	Keyword string
}


/*
 * 数据库相关内容
*/
// 数据库连接参数
type DB struct{
	User string
	Pass string
	Host string
	Port int
	Database string
	Instance *sqlx.DB
}

func (d *DB) Connect() {
	var dbUrl string = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Local&charset=utf8", d.User, d.Pass, d.Host, d.Port, d.Database)
	var db *sqlx.DB
	db, err := sqlx.Connect("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	d.Instance = db
}

func (d *DB) Insert(news News) {
	db := d.Instance
	if db == nil{
		panic("请初始化数据库连接")
	}
	tx := db.MustBegin()
	_, err := tx.NamedExec("INSERT INTO finance_news(title, abstract, url, savedate, source, keyword) VALUES(:title, :abstract, :url, :savedate, :source, :keyword)", news)
	if err != nil {
		_ = tx.Rollback()
	}else{
		_ = tx.Commit()
	}
}

func (d *DB) InsertAll(news []News) {
	for _, n := range news {
		d.Insert(n)
	}
}

func (d *DB) Close() {
	_ = d.Instance.Close()
}

// 从命令行获取数据库Host
func GetHost() string {
	args := os.Args
	if len(args) == 1{
		return "127.0.0.1"
	}else{
		return args[1]
	}
}