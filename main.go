package main

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		r := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		r.StaticFS("/static", http.FS(staticFiles))

		// 上传文件
		r.POST("api/v1/texts", TextController)
		// 获取局域网地址
		r.GET("api/v1/addresses", AddressesController)
		// 文件下载
		r.GET("api/v1/")


		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// 如果路径path是/static开头，就返回首页
			if strings.HasPrefix(path, "/static") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer reader.Close()

				stat, err := reader.Stat()  // Statistics返回文件统计信息
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
			} else {
				// 否则404
				c.Status(http.StatusNotFound)
			}
		})

		r.Run(":8080")
	}()

	chromePath := "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:8080/static/index.html")
	cmd.Start()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)  // Ctrl C 触发中断, 信号写入channel

	// 阻塞，直到接收到信号
	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}

func GetUploadsDir() (dir, uploads string) {
	// Executable 返回启动当前进程的可执行文件的路径
	// /home/cenjw/synk/synk.exe
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// 返回可执行文件的目录
	// /home/cenjw/synk/
	dir = filepath.Dir(exe)
	// 拼接上传路径：/home/cenjw/synk/uploads
	uploads = filepath.Join(dir, "uploads")
	return
}

func UploadsController(c *gin.Context) {
	if path := c.Param("path"); path != "" {
		_, uploads := GetUploadsDir()
		target := filepath.Join(uploads, path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename=" + path)
		c.Header("Content-Type", "application/octet-stream")
		// writes the specified file into the body stream in an efficient way
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"addresses":result,
	})
}

func TextController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		dir, uploads := GetUploadsDir()
		// 生成一个随机文件名：haitaos-hsjdfhk-sfhsk
		filename := uuid.New().String()
		// 创建目录
		err = os.MkdirAll(uploads, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		// 拼接上传文件的全路径：/home/cenjw/synk/uploads/haitaos-hsjdfhk-sfhsk.txt
		fullpath := path.Join("uploads", filename + ".txt")
		// 将用户传来的json.Raw数据写入文件
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// 返回全路径给前端
		c.JSON(http.StatusOK, gin.H{
			"url": "/" + fullpath,
		})
	}
}