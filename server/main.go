package main

import (
	"github.com/astaxie/beego"
	"thomas/controller"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"thomas/entity"
	"runtime"
	"github.com/astaxie/beego/logs"
)


func init() {
	logs.SetLogger("console")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root123@/d_live_transcode?charset=utf8&loc=Local")
	orm.RegisterModel(&entity.JobInfo{})

}
func main()  {


	logs.Debug("Thomas Live Transcode System start")
	// set CPU
	runtime.GOMAXPROCS(runtime.NumCPU())
	orm.Debug = true
	// TODO Init jobList
	
	// set home  path
	beego.Router("/",&controller.IndexController{},"get:Index")

	// jobinfo
	beego.Router("/jobinfo/list",&controller.JobInfoManagerController{},"*:List")
	beego.Router("/jobinfo/add",&controller.JobInfoManagerController{},"get:ToAdd")
	beego.Router("/jobinfo/add",&controller.JobInfoManagerController{},"post:Add")
	beego.Router("/jobinfo/edit",&controller.JobInfoManagerController{},"get:ToEdit")
	beego.Router("/jobinfo/edit",&controller.JobInfoManagerController{},"post:Edit")
	beego.Router("/jobinfo/delete",&controller.JobInfoManagerController{},"post:Delete")
	beego.Router("/jobinfo/info",&controller.JobInfoManagerController{},"get:Info")
	beego.Router("/jobinfo/operate",&controller.JobInfoManagerController{},"*:Operate")
	
	//about
	beego.Router("/about",&controller.AboutController{},"*:Index")

	// set static resource
	beego.SetStaticPath("static","static")
	beego.SetStaticPath("public","static")


	// start web app
	beego.Run()

}

