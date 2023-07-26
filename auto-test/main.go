package main

import (
	"jasondrogba/multi-client-cacheTest/auto-test/ec2test"
	"log"
	"time"
)

type ServerResponse struct {
	ReadUfs    []float64 `json:"readUfs"`
	ReadRemote []float64 `json:"remote"`
}

type CheckResponse struct {
	Load   string `json:"load"`
	Policy string `json:"policy"`
	Read   string `json:"read"`
	Train  string `json:"train"`
}

func main() {
	instanceMap := ec2test.Getec2Instance()
	MasterIp := instanceMap["Ec2Cluster-default-masters-0"]
	startUrl := "http://" + MasterIp + ":8080/startTraining"
	readUrl := "http://" + MasterIp + ":8080/readAlluxio"
	resultUrl := "http://" + MasterIp + ":8080/getMapResult"
	setPolicyUrl := "http://" + MasterIp + ":8080/setPolicy"
	checkStatusUrl := "http://" + MasterIp + ":8080/checkStatus"
	log.Println("startUrl:", startUrl)
	log.Println("readUrl:", readUrl)
	log.Println("resultUrl:", resultUrl)
	log.Println("setPolicyUrl:", setPolicyUrl)
	//设置JSON文件，包括每个worker的配置信息，和缓存策略
	//jsonBytes := StartSetJSON("LRU", "90", "20", "100", "200")

	log.Println("执行热门概率90%的测试，LRU和REPLICA分别执行3组，每组读取文件200次。")
	for i := 0; i < 4; i++ {
		allTest(startUrl, "LRU", "90", "200", checkStatusUrl)
		time.Sleep(2 * time.Second)
		allTest(startUrl, "REPLICA", "90", "200", checkStatusUrl)
		time.Sleep(2 * time.Second)
		dynamicTest(startUrl, readUrl, setPolicyUrl, "90", checkStatusUrl)
		time.Sleep(2 * time.Second)
	}

	log.Println("执行热门概率80%的测试，LRU和REPLICA分别执行3组，每组读取文件200次。")

	for i := 0; i < 4; i++ {
		allTest(startUrl, "LRU", "80", "200", checkStatusUrl)
		time.Sleep(2 * time.Second)
		allTest(startUrl, "REPLICA", "80", "200", checkStatusUrl)
		time.Sleep(2 * time.Second)
		dynamicTest(startUrl, readUrl, setPolicyUrl, "80", checkStatusUrl)
		time.Sleep(2 * time.Second)
	}

	log.Println("执行热门概率50%的测试，LRU和REPLICA分别执行3组，每组读取文件200次。")

	for i := 0; i < 4; i++ {
		allTest(startUrl, "LRU", "50", "200", checkStatusUrl)
		time.Sleep(2 * time.Second)
		allTest(startUrl, "REPLICA", "50", "200", checkStatusUrl)
		time.Sleep(2 * time.Second)
		dynamicTest(startUrl, readUrl, setPolicyUrl, "50", checkStatusUrl)
		time.Sleep(2 * time.Second)
	}

	log.Println("执行热门概率20%的测试，LRU和REPLICA分别执行3组，每组读取文件200次。")

	for i := 0; i < 4; i++ {
		allTest(startUrl, "LRU", "20", "200", checkStatusUrl)
		time.Sleep(2 * time.Second)
		allTest(startUrl, "REPLICA", "20", "200", checkStatusUrl)
		time.Sleep(2 * time.Second)
		dynamicTest(startUrl, readUrl, setPolicyUrl, "20", checkStatusUrl)
		time.Sleep(2 * time.Second)
	}

	log.Println("执行热门概率80%的测试，LRU和REPLICA分别执行3组，每组读取文件300次。")

	for i := 0; i < 4; i++ {
		allTest(startUrl, "LRU", "80", "300", checkStatusUrl)
		time.Sleep(2 * time.Second)
		allTest(startUrl, "REPLICA", "80", "300", checkStatusUrl)
		time.Sleep(2 * time.Second)
		//dynamicTest(startUrl, readUrl, setPolicyUrl, "80", checkStatusUrl)
		//time.Sleep(2 * time.Second)
	}
	log.Println("执行热门概率80%的测试，LRU和REPLICA分别执行3组，每组读取文件400次。")

	for i := 0; i < 4; i++ {
		allTest(startUrl, "LRU", "80", "400", checkStatusUrl)
		time.Sleep(2 * time.Second)
		allTest(startUrl, "REPLICA", "80", "400", checkStatusUrl)
		time.Sleep(2 * time.Second)
		//dynamicTest(startUrl, readUrl, setPolicyUrl, "80", checkStatusUrl)
		//time.Sleep(2 * time.Second)
	}
	////启动任务,根据设置的JSON，开始运行alluxio
	////RunStartFunc(jsonBytes, startUrl)
	//RunSetPolicyFunc(jsonBytes, setPolicyUrl)
	//time.Sleep(2 * time.Second)
	//运行的Alluxio，开始进行read数据
	//RunReadFunc(jsonBytes, readUrl)
	//
	//RunStartFunc(jsonBytes, startUrl)

	RunCheckStatusFunc(checkStatusUrl)

}

func allTest(startUrl string, policy string, ratio string,
	count string, checkStatusUrl string) {
	//启动全程policy测试
	jsonBytes := StartSetJSON(policy, ratio, "20", "100", count)
	RunStartFunc(jsonBytes, startUrl)

	//等待直到REPLICA训练完成
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		if RunCheckStatusFunc(checkStatusUrl).Train == "empty" {
			break
		}
	}

}

func dynamicTest(startUrl string, readUrl string,
	setPolicyUrl string, ratio string,
	checkStatusUrl string) {

	//前100次启动REPLICA策略
	jsonBytes := StartSetJSON("REPLICA", ratio, "20", "100", "100")
	RunStartFunc(jsonBytes, startUrl)

	//等待直到REPLICA训练完成
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		if RunCheckStatusFunc(checkStatusUrl).Train == "empty" {
			break
		}
	}

	//后100次切换成LRU策略
	jsonBytes = StartSetJSON("LRU", ratio, "20", "100", "100")
	RunSetPolicyFunc(jsonBytes, setPolicyUrl)
	//继续读取数据
	time.Sleep(2 * time.Second)
	RunReadFunc(jsonBytes, readUrl)

	//等待直到LRU策略读取完成
	for range ticker.C {
		if RunCheckStatusFunc(checkStatusUrl).Read == "empty" {
			break
		}
	}
}
