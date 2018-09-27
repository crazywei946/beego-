package main

import (
	_ "classOne/routers"
	"github.com/astaxie/beego"
)

func main() {
	//执行映射关系
	beego.AddFuncMap("prePage", prePage)
	beego.AddFuncMap("nextPage", nextPage)
	beego.Run()
}

func prePage(currentPage int) int {
	return currentPage - 1
}

func nextPage(currentPage int) int {
	return currentPage + 1

}
