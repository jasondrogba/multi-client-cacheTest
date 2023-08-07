package main

import (
	"jasondrogba/multi-client-cacheTest/auto-test/ec2test"
	"log"
	"strconv"
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

	log.Println("读取20000次，前1w次AI，后1w次DA。全程REPLICA和前REPLICA+后LRU各执行3组。")
	for i := 0; i < 2; i++ {
		//读取20000次，前1w次AI，后1w次DA。全程REPLICA
		switchTest(startUrl, readUrl, "LRU", "80", "20000", checkStatusUrl, false)
		time.Sleep(2 * time.Second)
		//switchTest(startUrl, readUrl, "REPLICA", "80", "20000", checkStatusUrl, false)
		//time.Sleep(2 * time.Second)
		//switchTest(startUrl, readUrl, "LRU", "20", "10", checkStatusUrl, false)
		//time.Sleep(2 * time.Second)
		////读取20000次，前1w次AI，后1w次DA。前1w次REPLICA，后1w次LRU
		//dynamicTest(startUrl, readUrl, setPolicyUrl, "REPLICA", "20", checkStatusUrl, "20000", false)
		//time.Sleep(2 * time.Second)
		////读取20000次，前1w次DA，后1w次AI。全程LRU
		//switchTest(startUrl, readUrl, "LRU", "80", "20000", checkStatusUrl, false)
		//time.Sleep(2 * time.Second)
		////读取20000次，前1w次DA，后1w次AI。前1w次LRU，后1w次REPLICA
		//dynamicTest(startUrl, readUrl, setPolicyUrl, "LRU", "80", checkStatusUrl, "20000", false)
		//time.Sleep(2 * time.Second)
		//读取20000次，前1w次DA，后1w次AI。前1w次LRU，后1w次REPLICA
		dynamicTest(startUrl, readUrl, setPolicyUrl, "LRU", "80", checkStatusUrl, "20000", false)
		time.Sleep(2 * time.Second)
		//读取20000次，前1w次AI，后1w次DA。前1w次REPLICA，后1w次LRU
		//dynamicTest(startUrl, readUrl, setPolicyUrl, "REPLICA", "20", checkStatusUrl, "20000", false)
		//time.Sleep(2 * time.Second)
	}

	log.Println("读取20000次，前1w次DA，后1w次AI。全程LRU和前LRU+后REPLICA各执行3组。")

	//log.Println("执行热门概率90%的测试，LRU和REPLICA分别执行3组，每组读取文件200次。")
	//for i := 0; i < 4; i++ {
	//	switchTest(startUrl, "LRU", "90", "200", checkStatusUrl)
	//	//time.Sleep(2 * time.Second)
	//	switchTest(startUrl, "REPLICA", "90", "200", checkStatusUrl)
	//	//time.Sleep(2 * time.Second)
	//	dynamicTest(startUrl, readUrl, setPolicyUrl, "90", checkStatusUrl, "200")
	//	time.Sleep(2 * time.Second)
	//}

	//log.Println("执行热门概率50%的测试，LRU和REPLICA分别执行3组，每组读取文件200次。")
	//
	//for i := 0; i < 4; i++ {
	//	//allTest(startUrl, "LRU", "50", "200", checkStatusUrl)
	//	//time.Sleep(2 * time.Second)
	//	//allTest(startUrl, "REPLICA", "50", "200", checkStatusUrl)
	//	//time.Sleep(2 * time.Second)
	//	dynamicTest(startUrl, readUrl, setPolicyUrl, "50", checkStatusUrl, "200")
	//	time.Sleep(2 * time.Second)
	//}
	//
	//log.Println("执行热门概率20%的测试，LRU和REPLICA分别执行3组，每组读取文件200次。")
	//
	//for i := 0; i < 4; i++ {
	//	switchTest(startUrl, "LRU", "20", "200", checkStatusUrl)
	//	//time.Sleep(2 * time.Second)
	//	switchTest(startUrl, "REPLICA", "20", "200", checkStatusUrl)
	//	//time.Sleep(2 * time.Second)
	//	dynamicTest(startUrl, readUrl, setPolicyUrl, "20", checkStatusUrl, "200")
	//	time.Sleep(2 * time.Second)
	//}
	//

}

func allTest(startUrl string, policy string, ratio string,
	hot string, total string,
	count string, checkStatusUrl string, warmUp bool) {
	//启动全程policy测试
	jsonBytes := StartSetJSON(policy, ratio, hot, total, count, warmUp)
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

func switchTest(startUrl string, readUrl string, policy string, ratio string,
	count string, checkStatusUrl string, warmUp bool) {

	//得到一半的count
	intCount, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		log.Println("err: ", err)
	}
	intHalfCount := intCount / 2
	strHalfCount := strconv.FormatInt(intHalfCount, 10)

	//得到冷门的ratio
	intRatio, err := strconv.ParseInt(ratio, 10, 64)
	if err != nil {
		log.Println("err: ", err)
	}
	intColdRatio := 100 - intRatio
	strColdRatio := strconv.FormatInt(intColdRatio, 10)

	//前一半的count的测试执行热门ratio的policy
	jsonBytes := StartSetJSON(policy, ratio, "20", "100", strHalfCount, warmUp)
	RunStartFunc(jsonBytes, startUrl)

	//等待直到REPLICA训练完成
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		if RunCheckStatusFunc(checkStatusUrl).Train == "empty" {
			break
		}
	}

	time.Sleep(2 * time.Second)
	//后一半的count的测试执行冷门ratio的policy
	jsonBytes = StartSetJSON(policy, strColdRatio, "20", "100", strHalfCount, warmUp)
	RunReadFunc(jsonBytes, readUrl)

	for range ticker.C {
		if RunCheckStatusFunc(checkStatusUrl).Read == "empty" {
			break
		}
	}

}

func dynamicTest(startUrl string, readUrl string,
	setPolicyUrl string, policy string, ratio string,
	checkStatusUrl string, count string, warmUp bool) {

	//得到一半的count
	intCount, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		log.Println("err: ", err)
	}
	intHalfCount := intCount / 2
	strHalfCount := strconv.FormatInt(intHalfCount, 10)

	//得到冷门的ratio
	intRatio, err := strconv.ParseInt(ratio, 10, 64)
	if err != nil {
		log.Println("err: ", err)
	}
	intColdRatio := 100 - intRatio
	strColdRatio := strconv.FormatInt(intColdRatio, 10)

	//前100次启动REPLICA策略
	jsonBytes := StartSetJSON(policy, ratio, "20", "100", strHalfCount, warmUp)
	RunStartFunc(jsonBytes, startUrl)

	//等待直到上一次策略的训练完成
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		if RunCheckStatusFunc(checkStatusUrl).Train == "empty" {
			break
		}
	}

	//后100次切换成LRU策略，同时反转热门概率
	if policy == "LRU" {
		policy = "REPLICA"
	} else {
		policy = "LRU"
	}
	jsonBytes = StartSetJSON(policy, strColdRatio, "20", "100", strHalfCount, warmUp)
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
