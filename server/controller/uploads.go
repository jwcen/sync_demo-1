package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

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
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		// writes the specified file into the body stream in an efficient way
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}
