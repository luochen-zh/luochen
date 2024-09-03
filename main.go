package main

import (
	_ "fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 结构体，用于存储用户信息
type User struct {
	Stu_ID uint   `gorm:"primaryKey"`
	Name   string `json:"name"`
	Class  string `json:"class"`
	Phone  string `json:"phone"`
}

var db *gorm.DB

// 初始化数据库连接
func initDB() {
	var err error
	dsn := "user:password@tcp(127.0.0.1:3306)/qr_form_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// 自动迁移User表
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

// 生成二维码
func generateQRCode(c *gin.Context) {
	// 假设二维码指向填写表单的URL
	url := "http://117.159.39.81:8080/form"
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate QR code")
		return
	}
	c.Data(http.StatusOK, "image/png", png)
}

// 渲染表单页面
func showForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
}

// 处理表单提交
func submitForm(c *gin.Context) {
	name := c.PostForm("name")
	class := c.PostForm("class")
	phone := c.PostForm("phone")

	// 创建用户记录
	user := User{Name: name, Class: class, Phone: phone}
	result := db.Create(&user)
	if result.Error != nil {
		c.String(http.StatusInternalServerError, "Failed to save user info")
		return
	}

	c.String(http.StatusOK, "Form submitted successfully!")
}

func main() {
	// 初始化数据库
	initDB()

	// 创建Gin路由器
	r := gin.Default()

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 生成二维码
	r.GET("/generate", generateQRCode)

	// 显示表单
	r.GET("/form", showForm)

	// 处理表单提交
	r.POST("/submit", submitForm)

	// 启动服务器
	r.Run(":8080")
}
