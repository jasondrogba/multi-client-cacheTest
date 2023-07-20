package main

import (
	"github.com/gin-gonic/gin"
	"jasondrogba/multi-client-cacheTest/worker/workerHandleFunc"
)

func main() {
	//使用gin框架实现一个接口，接收启动指令，启动一个worker
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
	r.GET("/loadAlluxio", workerHandleFunc.LoadAlluxioHandler)
	r.POST("/readAlluxio", workerHandleFunc.ReadAlluxioHandler)

	r.Run(":8888") // 监听并在
}
