package main

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/mvc"
	"publish_client_user/tool"
	"publish_client_user/web/controller"
	"publish_client_user/web/middleware"
)

func main() {
	app := iris.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	app.Use(iris.Gzip, logger.New(), crs)
	app.AllowMethods(iris.MethodOptions)
	app.Logger().SetLevel("debug")
	//用户路由 不需要token验证
	mvc.New(app.Party("/user")).Handle(&controller.UserController{})
	//app.Use(middleware.CheckSession)
	app.Use(middleware.GetJWT().Serve)
	//需要验证token
	mvc.New(app.Party("/users")).Handle(&controller.UsersController{})
	app.Run(iris.Addr(":" + tool.ConfJson.Port))
}
