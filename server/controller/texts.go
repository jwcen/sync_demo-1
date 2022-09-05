package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
		fullpath := path.Join("uploads", filename+".txt")
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
