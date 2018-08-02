package api

import "github.com/kataras/iris"


func HelloHandler(ctx iris.Context) {
	ctx.JSON(iris.Map{"message": "Hello World"})
}