package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"sync_demo/server/controller"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run() {
	port := ":27149"
	// 启动Gin服务
	go func() {
		gin.SetMode(gin.DebugMode)
		r := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		r.StaticFS("/static", http.FS(staticFiles))

		// 上传文本
		r.POST("/api/v1/texts", controller.TextController)
		// 获取局域网地址
		r.GET("/api/v1/addresses", controller.AddressesController)
		// 文件下载
		r.GET("/uploads/:path", controller.UploadsController)
		// 生成二维码
		r.GET("/api/v1/qrcodes", controller.QrcodesController)
		// 文件上传
		r.POST("/api/v1/files", controller.FilesController)

		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// 如果路径path是/static开头，就返回首页
			if strings.HasPrefix(path, "/static") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer reader.Close()

				stat, err := reader.Stat() // Statistics返回文件统计信息
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
			} else {
				// 否则404
				c.Status(http.StatusNotFound)
			}
		})

		r.Run(port)
	}()
}
