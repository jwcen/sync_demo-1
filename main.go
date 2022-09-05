package main

import (
	"os"
	"os/exec"
	"os/signal"
	"sync_demo/server"
)

func main() {
	go server.Run()

	cmd := startBrower()

	// 监听中断信号
	chSignal := listenToInterrupt()

	// 阻塞，直到接收到信号
	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}

func startBrower() *exec.Cmd {
	chromePath := "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:27149/static/index.html")
	cmd.Start()
	return cmd
}

func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt) // Ctrl C 触发中断, 信号写入channel
	return chSignal
}
