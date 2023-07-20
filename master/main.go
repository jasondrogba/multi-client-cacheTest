package main

import (
	"github.com/gin-gonic/gin"
	"jasondrogba/multi-client-cacheTest/master/masterHandleFunc"
	"jasondrogba/multi-client-cacheTest/master/readyForEc2"
)

func main() {
	//使用gin框架实现一个接口，接收启动指令，启动一个worker
	readyForEc2.Prepare()

	r := gin.Default()
	r.GET("/start", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "start",
		})
	})
	r.GET("/stop", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "stop",
		})
	})
	r.GET("/stopWorker", masterHandleFunc.StopAlluxioHandler)
	r.GET("/startAlluxio", masterHandleFunc.StartAlluxioHandler)
	r.POST("/loadAlluxio", masterHandleFunc.LoadAlluxioHandler)
	r.POST("/readAlluxio", masterHandleFunc.ReadAlluxioHandler)
	r.GET("/startTraining", masterHandleFunc.StartTrainingHandler)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
