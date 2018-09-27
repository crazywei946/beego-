package routers

import (
	"classOne/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {

	beego.InsertFilter("/article/*",beego.BeforeRouter, func(ctx *context.Context) {
		userName:=ctx.Input.Session("userName")
		if userName==nil {
			ctx.Redirect(302,"/")
		}
	})

	//注册
	beego.Router("/register", &controllers.RegisterController{}, "get:Register")
	beego.Router("/register", &controllers.RegisterController{}, "post:ShowRegister")
	//登陆
	beego.Router("/", &controllers.LoginController{}, "get:ShowLogin")
	beego.Router("/login", &controllers.LoginController{}, "post:HandleLogin")
	//添加文章

	beego.Router("/article/showArticleList", &controllers.ArticleController{}, "get:ShowArticleList")
	beego.Router("/article/showAddArticle", &controllers.ArticleController{}, "get:ShowAddArticle")
	beego.Router("/article/showAddArticle", &controllers.ArticleController{}, "post:HandleAddArticle")

	//文章详情
	beego.Router("/article/content", &controllers.ArticleController{}, "get:ContentArticle")
	//删除文章
	beego.Router("/article/deleteArticle", &controllers.ArticleController{}, "get:DeleteArticle")

	//文章更新
	beego.Router("/article/updateArticle", &controllers.ArticleController{}, "get:ShowArticle;post:UpdateArticle")

	//添加分类
	beego.Router("/article/addType", &controllers.ArticleController{},"get:ShowType;post:HandleAddType")
	//删除文章类型
	beego.Router("/article/deleteArticleType", &controllers.ArticleController{},"get:DeleteArticleType")
	//推出功能
	beego.Router("article/quit", &controllers.ArticleController{},"get:Quit")

}
