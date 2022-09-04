package main

import (
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		r := gin.Default()
		
		r.GET("/", func(ctx *gin.Context) {
			ctx.Writer.Write([]byte("hello world"))
		})

		r.Run(":8080")
	}()

	chromePath := "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:8080/")
	cmd.Start()

	select {}
}