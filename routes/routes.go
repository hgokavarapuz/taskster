package routes

import "github.com/kataras/iris"

func renderTemplate(ctx iris.Context, tmpl string, p interface{}) {
	ctx.StatusCode(iris.StatusOK)
	ctx.View(tmpl+".html", p)
}