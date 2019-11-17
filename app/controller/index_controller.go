/**
 * @Author: Rocks
 * @Description: 
 * @File:  index_controller
 * @Version: 1.0.0
 * @Date: 2019-10-13 22:05
 */

package controller

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

/**
首页
 */
func Index(r *ghttp.Request) {
	err := r.Response.WriteTpl("index.html", g.Map{})
	if err != nil {
		glog.Error(err)
	}
}

func Main(r *ghttp.Request) {
	err := r.Response.WriteTpl("main.html", g.Map{})
	if err != nil {
		glog.Error(err)
	}
}
