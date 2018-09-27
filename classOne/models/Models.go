package models

import (
	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//存放1.表结构 2.连接数据库

type User struct {
	Id       int
	Name     string
	PassWord string
	Articles []*Article `orm:"rel(m2m)"`
}

//文章
type Article struct {
	Id          int          `orm:"pk;auto"'`
	ArtiName    string       `orm:"size(20)""`
	Content     string       `orm:"size(500)"`
	Time        time.Time    `orm:"auto_now"`
	Count       int          `orm:"default(0);null"`
	Img         string       `orm:"size(50)"`
	ArticleType *ArticleType `orm:"rel(fk);null;on_delete(set_null)"`
	Users       []*User      `orm:"reverse(many)"`
}

//文章结构类型表
type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}

func init() {
	//	连接数据库
	err := orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
	if err != nil {
		beego.Error("注册数据库失败*****", err)
		return
	}
	//映射表数据
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	//	生成表
	orm.RunSyncdb("default", false, true)
}
