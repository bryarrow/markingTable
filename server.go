package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func main() {
	r := gin.Default()
	r.StaticFS("/src", http.Dir("./src"))
	r.StaticFile("/favicon.ico", "./src/pic/logo.png")
	r.LoadHTMLGlob("./tmp/*")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "每日评价",
			"navText": "每日评价",
		})
	})

	var province string
	var city string
	groupLogin := r.Group("/login")
	{

		groupLogin.GET("admin", func(ctx *gin.Context) {

			ctx.HTML(http.StatusOK, "login.html", gin.H{
				"title":   "每日评价-管理员",
				"navText": "管理员登录",
				"Type":    "admin",
			})
		})
		groupLogin.GET("Teacher", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", gin.H{
				"title":   "每日评价-助教",
				"navText": "助教登录",
				"Type":    "Teacher",
			})
		})
		groupLogin.GET("Parents", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", gin.H{
				"title":   "每日评价-学生/家长",
				"navText": "学生/家长登录",
				"Type":    "Parents",
			})
		})
		groupLogin.POST("admin", func(ctx *gin.Context) {
			//表单数据获取
			province = ctx.PostForm("province")
			city = ctx.PostForm("city")
			//管理员数据读取
			adminDataFile, err := os.Open("./data/admins.json")
			if err != nil {
				panic(err)
			}
			defer adminDataFile.Close()
			adminData, err := ioutil.ReadAll(adminDataFile)
			if err != nil {
				panic(err)
			}
			//数据验证并进入密码输入页面
			isFirst := gjson.Get(string(adminData), province+"."+city+".isFirst")
			if isFirst.Bool() {
				ctx.HTML(http.StatusOK, "login.html", gin.H{
					"title":   "每日评价-管理员",
					"navText": "管理员登录",
					"Type":    "adminNext",
					"IsFirst": "True",
				})
			} else {
				ctx.HTML(http.StatusOK, "login.html", gin.H{
					"title":   "每日评价-管理员",
					"navText": "管理员登录",
					"Type":    "adminNext",
					"IsFirst": "False",
				})
			}
		})

		groupLogin.POST("admin/next", func(ctx *gin.Context) {
			//表单数据获取
			passwordInForm := ctx.PostForm("password")
			//管理员数据读取
			adminDataFile, err := os.Open("./data/admins.json")
			if err != nil {
				panic(err)
			}
			defer adminDataFile.Close()
			adminData, err := ioutil.ReadAll(adminDataFile)
			if err != nil {
				panic(err)
			}
			//数据验证、存储并登录
			isFirst := gjson.Get(string(adminData), province+"."+city+".isFirst")
			if isFirst.Bool() {
				//计算设置的密码的MD5
				has := md5.Sum([]byte(passwordInForm))
				passwordMd5InForm := fmt.Sprintf("%x", has)
				//设置密码
				updatedJson, err := sjson.Set(string(adminData), province+"."+city+".passwordMD5", passwordMd5InForm)
				if err != nil {
					panic(err)
				}
				updatedJson, err = sjson.Set(updatedJson, province+"."+city+".isFirst", false)
				if err != nil {
					panic(err)
				}
				err = ioutil.WriteFile("./data/admins.json", []byte(updatedJson), 0644)
				if err != nil {
					panic(err)
				}
				//进入管理页面
				//TODO: 用js控制进入管理页面
				ctx.HTML(http.StatusOK, "admin.html", gin.H{
					"title":   "每日评价-管理员",
					"navText": "管理员控制台",
				})
			} else {
				//计算设置的密码的MD5
				has := md5.Sum([]byte(passwordInForm))
				passwordMd5InForm := fmt.Sprintf("%x", has)
				//验证密码
				passwordMD5 := gjson.Get(string(adminData), province+"."+city+".passwordMD5")
				if passwordMD5.String() == passwordMd5InForm {
					//进入管理页面
					ctx.HTML(http.StatusOK, "admin.html", gin.H{
						"title":   "每日评价-管理员",
						"navText": "管理员控制台",
					})
				} else {
					//提示密码错误
					//TODO: 更好的错误提示页面
					ctx.String(http.StatusOK, "您输入的密码有误")
				}
			}
		})
	}

	r.Run(":2333")
}
