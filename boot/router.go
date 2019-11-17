package boot

import (
	"autocreate/app/controller"
	"github.com/gogf/gf/frame/g"
)

func InitRouter() {
	urlPath := g.Config().GetString("url-path")
	s := g.Server()
	s.BindHandler(urlPath+"/", controller.Index)
	s.BindHandler(urlPath+"/main.html", controller.Main)
	s.BindHandler(urlPath+"/generator.html", controller.Gen)
	s.BindHandler(urlPath+"/generator/list", controller.List)
	s.BindHandler(urlPath+"/generator/code", controller.Code)
}
