package masterHandleFunc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jasondrogba/multi-client-cacheTest/master/ec2test"
	"jasondrogba/multi-client-cacheTest/master/handleLock"
	"jasondrogba/multi-client-cacheTest/master/loadAlluxio"
	"jasondrogba/multi-client-cacheTest/master/readyForEc2"
	"jasondrogba/multi-client-cacheTest/master/startTest"
	"jasondrogba/multi-client-cacheTest/master/userMasterInfo"
	"log"
	"net/http"
)

func StopAlluxioHandler(*gin.Context) {
	//启动一个worker

	stopResp, err := http.Get("http://localhost:8888/stop")
	if err != nil {
		log.Fatalf("http.Get err: %v", err)
	}
	defer stopResp.Body.Close()
	if stopResp.StatusCode != http.StatusOK {
		log.Fatalf("停止失败，状态码：%d", stopResp.StatusCode)
	}

	fmt.Println("已发送停止命令，服务器返回：", stopResp.Status)
}

func GetWorkerAddr() map[string]string {
	instanceMap := ec2test.Getec2Instance()
	return instanceMap
}

// StartAlluxioHandler 与master通信，执行free，stop，format，start，选择cache eviction policy
func StartAlluxioHandler(c *gin.Context) {
	instanceMap := readyForEc2.GetInstanceMap()
	select {
	case handleLock.Getrunning() <- struct{}{}: // 尝试获取互斥锁
		// 成功获取互斥锁，执行处理函数
		fmt.Println("后台处理开始")
		go startTest.StartTest(instanceMap["Ec2Cluster-default-masters-0"], "LRU")
		c.JSON(200, gin.H{
			"message": "启动Alluxio",
		})
	default:
		// 未获取到互斥锁，处理函数正在执行中，直接返回错误响应
		c.JSON(500, gin.H{
			"message": "后台处理中",
		})
	}

}

// LoadAlluxioHandler 接收一个设置的列表，可以设置不同的worker分别预热多少的数据
func LoadAlluxioHandler(c *gin.Context) {

	//向所有的alluxio worker发送预热数据的命令
	var workerListInfo userMasterInfo.WorkerInfoList
	err := c.BindJSON(&workerListInfo)
	if err != nil {
		log.Println("JSON err:", err)
		c.JSON(400, gin.H{"error": "解析JSON数据失败"})
		return
	}

	select {
	case handleLock.GetLoadRunning() <- struct{}{}: // 尝试获取互斥锁
		// 成功获取互斥锁，执行处理函数
		fmt.Println("后台处理开始")
		go loadAlluxio.TotalLoad(workerListInfo)

		c.JSON(200, gin.H{
			"message": "在worker中预热Alluxio",
		})
	default:
		// 未获取到互斥锁，处理函数正在执行中，直接返回错误响应
		c.JSON(500, gin.H{
			"message": "后台处理中",
		})
	}

}

func ReadAlluxioHandler(c *gin.Context) {
	//向所有的alluxio worker发送读取数据的命令
	var workerListInfo userMasterInfo.WorkerInfoList
	err := c.BindJSON(&workerListInfo)
	if err != nil {
		log.Println("JSON err:", err)
		c.JSON(400, gin.H{"error": "解析JSON数据失败"})
		return
	}

	select {
	case handleLock.GetReadRunning() <- struct{}{}: // 尝试获取互斥锁
		// 成功获取互斥锁，执行处理函数
		fmt.Println("后台处理开始")
		go loadAlluxio.TotalRead(workerListInfo)

		c.JSON(200, gin.H{
			"message": "在worker中预热Alluxio",
		})
	default:
		// 未获取到互斥锁，处理函数正在执行中，直接返回错误响应
		c.JSON(500, gin.H{
			"message": "后台处理中",
		})
	}
}

func StartTrainingHandler(c *gin.Context) {
	var workerListInfo userMasterInfo.WorkerInfoList
	err := c.BindJSON(&workerListInfo)
	if err != nil {
		log.Println("JSON err:", err)
		c.JSON(400, gin.H{"error": "解析JSON数据失败"})
		return
	}
	select {
	case handleLock.GetTrainRunning() <- struct{}{}: // 尝试获取互斥锁
		// 成功获取互斥锁，执行处理函数
		fmt.Println("后台处理开始")
		go startTest.StartTraining(workerListInfo)

		c.JSON(200, gin.H{
			"message": "开始训练",
		})
	default:
		// 未获取到互斥锁，处理函数正在执行中，直接返回错误响应
		c.JSON(500, gin.H{
			"message": "后台处理中，还没有训练结束",
		})
	}

}

func GetResultHandler(c *gin.Context) {
	//remote := []float64{0.1, 0.2, 0.3}
	//readUfs := []float64{0.1, 0.2, 0.3}
	remote, readUfs := startTest.GetResult()
	c.JSON(200, gin.H{
		"message": "获取结果",
		"remote":  remote,
		"readUfs": readUfs,
	})
}
