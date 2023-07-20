package workerHandleFunc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jasondrogba/multi-client-cacheTest/worker/loadAlluxio"
	"jasondrogba/multi-client-cacheTest/worker/workerHandleLock"
	"jasondrogba/multi-client-cacheTest/worker/workerInfo"
	"log"
	"strconv"
)

func LoadAlluxioHandler(c *gin.Context) {
	//fmt.Println("收到预热请求")
	fmt.Println("收到预热请求", c.Query("loadFile"))
	fileCount, err := strconv.Atoi(c.Query("loadFile"))
	if err != nil {
		fmt.Println("预热文件数量转换失败")
		return
	}
	select {
	case workerHandleLock.GetLoadRunning() <- struct{}{}: // 尝试获取互斥锁
		// 成功获取互斥锁，执行处理函数
		fmt.Println("后台处理开始")
		go loadAlluxio.LoadAlluxio(fileCount)
		c.JSON(200, gin.H{
			"message": "在Alluxio worker中预热数据",
		})

	default:
		// 未获取到互斥锁，处理函数正在执行中，直接返回错误响应
		c.JSON(500, gin.H{
			"message": "后台处理中",
		})
	}

}

func ReadAlluxioHandler(c *gin.Context) {
	//fmt.Println("收到读取请求", c.Query("readRatio"))
	var masterWorkerInfo workerInfo.RequestData
	err := c.BindJSON(&masterWorkerInfo)
	if err != nil {
		fmt.Println("读取请求绑定失败")
		return
	}
	log.Println("收到读取请求", masterWorkerInfo.ReadRatio, masterWorkerInfo.HotFile,
		masterWorkerInfo.TotalFile, masterWorkerInfo.Count, masterWorkerInfo.MasterIP)
	masterIP := masterWorkerInfo.MasterIP
	readRatio, err := strconv.Atoi(masterWorkerInfo.ReadRatio)
	if err != nil {
		fmt.Println("读取请求绑定失败")
		return
	}
	hotFile, err := strconv.Atoi(masterWorkerInfo.HotFile)
	if err != nil {
		fmt.Println("读取请求绑定失败")
		return
	}
	totalFile, err := strconv.Atoi(masterWorkerInfo.TotalFile)
	if err != nil {
		fmt.Println("读取请求绑定失败")
		return
	}
	count, err := strconv.Atoi(masterWorkerInfo.Count)
	if err != nil {
		fmt.Println("读取请求绑定失败")
		return
	}
	select {
	case workerHandleLock.GetLoadRunning() <- struct{}{}: // 尝试获取互斥锁
		// 成功获取互斥锁，执行处理函数
		fmt.Println("后台处理开始")
		go loadAlluxio.ReadAlluxio(masterIP, count, readRatio, hotFile, totalFile)
		c.JSON(200, gin.H{
			"message": "在Alluxio worker中预热数据",
		})

	default:
		// 未获取到互斥锁，处理函数正在执行中，直接返回错误响应
		c.JSON(500, gin.H{
			"message": "后台处理中",
		})
	}

}
