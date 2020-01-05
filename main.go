package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/astaxie/beego/session/mysql"
	"github.com/astaxie/beego/toolbox"
	_ "github.com/go-sql-driver/mysql"
	EmailScripts "go-core/middlewares/email/scripts"
	_ "go-core/routers"
)

func main() {
	sqlconn := beego.AppConfig.String("core_dns")
	err := orm.RegisterDataBase("default", "mysql", sqlconn)
	if err != nil {
		fmt.Println("Core database connect error!")
		return
	}
	orm.SetMaxIdleConns("default", 50)
	orm.SetMaxOpenConns("default", 100)

	tk := toolbox.NewTask("EmailTask", "0/10 * * * * *", func() error {
		return EmailScripts.EmailTask()
	})
	toolbox.AddTask("EmailTask", tk)

	// runs in all modes
	toolbox.StartTask()
	defer toolbox.StopTask()

	UploadDir := beego.AppConfig.String("UploadDir")
	beego.SetStaticPath("/Uploads", UploadDir)

	cdnUrls := beego.AppConfig.Strings("cdn_urls")
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowCredentials: true,
		AllowOrigins:     cdnUrls,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Authorization", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Host", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Expose-Headers", "Content-Type"},
	}))
	beego.Run()
}
