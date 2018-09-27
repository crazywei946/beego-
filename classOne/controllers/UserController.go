package controllers

import (
	_ "classOne/models"
	"github.com/astaxie/beego"

	"classOne/models"
	"github.com/astaxie/beego/orm"
	"path"

	"time"

	"math"
	"github.com/gomodule/redigo/redis"
	"bytes"
	"encoding/gob"
)

type RegisterController struct {
	beego.Controller
}

//登陆
type LoginController struct {
	beego.Controller
}

//文章相关
type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowArticleList() {
	//this.TplName = "index.html"

	//判断session是否存在
	//userName:=this.GetSession("userName")
	//if userName==nil {
	//	this.Redirect("/",302)
	//	return
	//}

	//查询数据
	o := orm.NewOrm()
	//article:=new(models.Article)

	qs := o.QueryTable("article")
	var articles []models.Article
	//qs.All(&articles)

	//查询出所有的文章类型
	articleTypes := make([]models.ArticleType, 0)
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		beego.Error("连接数据库失败")
		return
	}
	defer conn.Close()
	//如果是第一次不是第一次加载，就从redis数据库去数据库
	rel, _ := redis.Bytes(conn.Do("get", "articleTypes"))
	//将获取出的字节进行反序列化
	dec := gob.NewDecoder(bytes.NewReader(rel)) //解码器
	dec.Decode(&articleTypes)                   //解码


	if len(articleTypes) == 0 {
		_, err := o.QueryTable("ArticleType").All(&articleTypes)
		if err != nil {
			beego.Error("查寻文章类型表失败")
			return
		}
		beego.Error("mysql取出")
		//将文章类型保存到redis中
		//先将要存储的内容进行序列化
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(articleTypes)
		_, err = conn.Do("set", "articleTypes", &buffer)
		if err != nil {
			beego.Error("数据存储失败")
			return
		}

	}else {
		beego.Error("redis取出**")
	}


	this.Data["articleTypes"] = articleTypes

	//获取传递过来的当前页码
	currenPage, _ := this.GetInt("currenPage")
	if currenPage == 0 {
		currenPage = 1
	}

	//获得表但提交过来的typeame
	typename := this.GetString("select")
	//panduan 是否有选择类型
	var Count int64
	if typename == "" || typename == "所有分类" {
		Count, _ = qs.RelatedSel("ArticleType").Count()

	} else {
		Count, _ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typename).Count()
	}
	//获得总条数
	//默认当前页为第一页

	//每页显示条数
	pageSize := 2
	//计算总页数
	totalPage := math.Ceil(float64(Count) / float64(pageSize))

	//定义标志位用来判断是否到了第一页或者最后一页
	fristFlag := false
	endFlag := false
	if currenPage == 1 {
		//如果是首页，就将标志位置为ture
		fristFlag = true
	}

	if currenPage == int(totalPage) {
		//如果是首页，就将标志位置为ture
		endFlag = true
	}

	//查询每页显示的内容
	//起始查询位置
	start := pageSize * (currenPage - 1)

	this.Data["typename"] = typename
	beego.Error(typename)
	if typename == "" || typename == "所有分类" {
		//如果传递过来的typename为空，则表示这是无分类默认查询
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	} else {
		//不为空则表示带条件查询
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typename).All(&articles)
	}

	this.Data["fristFlag"] = fristFlag
	this.Data["endFlag"] = endFlag
	this.Data["Count"] = Count
	this.Data["totalPage"] = totalPage
	this.Data["currenPage"] = currenPage

	//beego.Info(len(articles))

	this.Data["articles"] = articles
	this.Layout = "layout.html"
	this.TplName = "index.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["title"] = "title.html"

}

func (this *ArticleController) ShowAddArticle() {
	//查询出所有的文章内容显示到页面上
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	_, err := o.QueryTable("ArticleType").All(&articleTypes)
	if err != nil {
		beego.Error("查询文章类型失败")
		return
	}

	this.Data["articleTypes"] = articleTypes
	this.Layout = "layout.html"

	this.TplName = "add.html"
	this.Data["title"] = "文章内容"

}

//添加文章类型显示
func (this *ArticleController) ShowType() {
	//查询出所有的已经存在type并且显示
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	_, err := o.QueryTable("ArticleType").All(&articleTypes)
	if err != nil {
		beego.Error("查询数据失败")
	}

	this.Data["articleTypes"] = articleTypes
	this.Layout = "layout.html"
	this.TplName = "addType.html"
	this.Data["title"] = "添加分类"
}

//文章类型添加
func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		beego.Error("获得文章类型失败")
	}
	//将内容插入表中
	o := orm.NewOrm()
	articleType := models.ArticleType{}
	articleType.TypeName = typeName
	_, err := o.Insert(&articleType)
	if err != nil {
		beego.Error("插入文章到数据库中失败")
		return
	}
	//重定向到添加类型页面
	this.Redirect("/article/addType", 302)
}

//删除文章分类
func (this *ArticleController) DeleteArticleType() {
	//获得id
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取id失败")
		return
	}
	o := orm.NewOrm()
	articleType := models.ArticleType{Id: id}
	o.Delete(&articleType)

	this.Redirect("/article/addType", 302)
}

//文章详情
func (this *ArticleController) ContentArticle() {
	//this.Ctx.WriteString("******")
	//获得传递过来的参数

	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	article := models.Article{Id: id}
	err := o.Read(&article)
	if err != nil {
		beego.Error("查询数据失败")
		return
	}
	//每次点击一次就将阅读量加1
	article.Count += 1
	//多对多插入读者
	//1.获得当前文章为article
	//2.获得文章的关联的users属性
	m2m := o.QueryM2M(&article, "Users")
	//3.给users设置属性，通过属性查询出id
	userName := this.GetSession("userName")
	user := models.User{}
	user.Name = userName.(string)
	//查询出此user的id
	o.Read(&user, "Name")
	//插入
	_, err = m2m.Add(&user)
	if err != nil {
		beego.Error("插入失败")
		return
	}

	o.Update(&article)

	//查询阅读人
	//o.LoadRelated(&article,"Users")//可查询到却无法去重
	//多表查询方式2
	var users []*models.User
	o.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)

	this.Data["article"] = article
	this.Data["users"] = users
	this.Layout = "layout.html"

	this.TplName = "content.html"

}

//文章删除
func (this *ArticleController) DeleteArticle() {
	//获得需要删除的id
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	article := models.Article{Id: id}
	_, err := o.Delete(&article)
	if err != nil {
		beego.Error("删除文章失败", err)
		return
	}

	//重定向到文章列表页面
	this.Redirect("/article/showArticleList", 302)

}

//文章更新
func (this *ArticleController) ShowArticle() {
	//获得需要更新的文章的id
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	article := models.Article{Id: id}
	err := o.Read(&article)
	if err != nil {
		beego.Error("查询失败", err)
		return
	}
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

func (this *ArticleController) UpdateArticle() {
	//获取原有内容
	id, _ := this.GetInt("id")

	articleName := this.GetString("articleName")
	content := this.GetString("content")
	f, h, err := this.GetFile("uploadname")
	if err != nil {
		beego.Error("图片文件获得失败")
		return
	}
	defer f.Close()

	//获得图片文件的后缀

	filename := h.Filename
	ext := path.Ext(filename)
	filesize := h.Size
	if ext != ".jpg" && ext != ".png" {
		beego.Error("文件格式不正确")
		return
	}
	if filesize > 1024*1024*5 {
		beego.Error("文件台大")
		return
	}
	//上传文件到本地
	filepath := "./static/img/" + filename + ext
	err = this.SaveToFile("uploadname", filepath)
	if err != nil {
		beego.Error("图片文件上传失败", err)
		return
	}
	//将此内容村初到数据库中
	o := orm.NewOrm()
	article := models.Article{Id: id}
	article.ArtiName = articleName
	article.Content = content
	article.Img = filepath

	o.Update(&article)

	//重定向到文章列表
	this.Redirect("/article/showArticleList", 302)

}

func (this *ArticleController) HandleAddArticle() {
	//获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	//获得文章内性
	typeName := this.GetString("select")
	if articleName == "" || content == "" || typeName == "" {
		beego.Error("获取文章内容失败")
	}
	var articletype models.ArticleType
	articletype.TypeName = typeName

	//beego.Info(articleName, content)

	//获得文件
	f, h, err := this.GetFile("uploadname")
	if err != nil {
		beego.Error("文件获取失败", err)
		return
	}

	defer f.Close()

	//uploadname:=h.Filename
	ext := path.Ext(h.Filename)

	//panduan文件是否为空

	if h.Filename == "" {
		ext = ".jpg"
	}

	//判断文件格式
	if ext != ".jpg" && ext != "png" {
		beego.Error("文件格式不正确")
		return
	}

	if h.Size > 1024*1024 {
		beego.Error("文件太大无法上传")
		return
	}

	//给文件起个名防止重名
	filename := time.Now().Format("2006-01-02 15:04:05") + ext

	//评介文件上传路径
	filepath := "./static/img/" + filename

	err = this.SaveToFile("uploadname", filepath)
	if err != nil {
		beego.Error("文件上传失败***", err)
		return
	}

	//将内容存储到数据库
	o := orm.NewOrm()

	//根据文章呢类型名称查出文章类型
	err = o.Read(&articletype, "TypeName")
	if err != nil {
		beego.Error("查询文章类型失败")
		return
	}
	article := models.Article{}
	article.Content = content
	article.ArtiName = articleName
	article.Img = filepath
	article.ArticleType = &articletype

	o.Insert(&article)

	//重定向到文章列表页面
	this.Redirect("/article/showArticleList", 302)

}

//登陆
func (this *LoginController) ShowLogin() {
	name := this.Ctx.GetCookie("name")
	if name != "" {
		//this.Data["check"]="checked"
		this.Data["name"] = name
	}

	this.TplName = "login.html"
}

func (this *LoginController) HandleLogin() {
	//判断得到的数据
	name := this.GetString("userName")
	password := this.GetString("password")
	if name == "" || password == "" {
		beego.Info("用户名和密码不能为空")
		this.TplName = "login.html"
		return
	}

	remember := this.GetString("remember")

	//判断是否需要保存用户名
	if remember == "on" {
		this.Ctx.SetCookie("name", name, 600*10)
	} else {
		this.Ctx.SetCookie("name", name, -1)
	}

	//用户密码进行比对
	o := orm.NewOrm()
	user := models.User{}
	user.Name = name
	user.PassWord = password
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Info("用户名错误", err)
		return
	}
	//判断密码
	if user.PassWord != password {
		beego.Info("密码错误", password)
		return
	}

	//登陆成功，创建session
	this.SetSession("userName", name)

	this.Redirect("/article/showArticleList", 302)

}

//注册
func (this *RegisterController) Register() {
	this.TplName = "register.html"
}

//注册
func (this *RegisterController) ShowRegister() {
	//	获得传递过来的参数
	name := this.GetString("userName")
	password := this.GetString("password")
	//beego.Info(name,password)
	if name == "" || password == "" {
		this.TplName = "register.html"
		return
	}

	//插入数据
	o := orm.NewOrm()
	user := models.User{}
	user.Name = name
	user.PassWord = password
	o.Insert(&user)

	this.TplName = "login.html"

}

//推出
func (this *ArticleController) Quit() {
	//清空session
	this.DelSession("userName")

	this.TplName = "login.html"

}
