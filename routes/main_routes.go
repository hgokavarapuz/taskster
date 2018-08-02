package routes

import "github.com/kataras/iris"

func IndexHandler(ctx iris.Context) {
	//ctx.HTML("<h1>Welcome to Taskster</h1>")
	renderTemplate(ctx, "index",nil)
}

func IdHandler(ctx iris.Context) {
	renderTemplate(ctx, "index",nil)
}