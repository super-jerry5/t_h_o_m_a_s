package main

import (
	"github.com/astaxie/beego"
	"thomas/controller"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"thomas/entity"
	"runtime"
	"github.com/astaxie/beego/logs"
	config "github.com/astaxie/beego/config"
)


func init() {
	ini, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		panic(err)
	}
	mysql_string := ini.String("mysql_string")
	logs.SetLogger(logs.AdapterFile,`{"filename":"/data/logs/thomas/thomas.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	logs.SetLogFuncCall(true)
 	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysql_string)
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

