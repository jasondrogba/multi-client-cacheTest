package loadAlluxio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jasondrogba/multi-client-cacheTest/master/handleLock"
	"jasondrogba/multi-client-cacheTest/master/readyForEc2"
	"jasondrogba/multi-client-cacheTest/master/userMasterInfo"
	"log"
	"net/http"
	"runtime"
	"sync"
)

var readwg sync.WaitGroup

func TotalRead(workerListInfo userMasterInfo.WorkerInfoList) {
	instanceMap := readyForEc2.GetInstanceMap()
	MasterIp := instanceMap["Ec2Cluster-default-masters-0"]
	for _, w := range workerListInfo.WorkerInfoList {
		readwg.Add(1)
		fmt.Printf("WorkerID: %s, ReadRatio: %s, LoadFile: %s, "+
			"HotFile: %s, TotalFile: %s, Count: %s,MasterIP: %s\n",
			w.WorkerId, w.ReadRatio, w.LoadFile, w.HotFile, w.TotalFile, w.Count, MasterIp)
		go multiRead(w.WorkerId, MasterIp, w.Count, w.HotFile, w.TotalFile, w.ReadRatio)
	}
	readwg.Wait()
	<-handleLock.GetReadRunning()
}

func multiRead(workerId string, masterIp string, count string,
	hotFile string, totalFile string, readRatio string) error {
	instanceMap := readyForEc2.GetInstanceMap()
	if runtime.GOOS == "linux" {
		fmt.Println("Detected Linux system")
		err := sendReadToWorker(instanceMap["Ec2Cluster-default-workers-"+workerId]+":8888/readAlluxio",
			masterIp, count, hotFile, totalFile, readRatio)
		if err != nil {
			return err
		}
	} else if runtime.GOOS == "darwin" {
		fmt.Println("Detected macOS system")
		err := sendReadToWorker("localhost:8888/readAlluxio",
			masterIp, count, hotFile, totalFile, readRatio)
		if err != nil {
			return err
		}

	} else {
		fmt.Println("Unknown system")
	}
	return nil
}

func sendReadToWorker(workerIP string, masterIp string, count string,
	hotFile string, totalFile string, readRatio string) error {
	//fmt.Println("loadFile", loadFile)
	//将参数打包进JSON
	readwg.Done()

	// 创建一个包含参数的结构体实例
	requestData := userMasterInfo.RequestData{
		MasterIP:  masterIp,
		Count:     count,
		HotFile:   hotFile,
		TotalFile: totalFile,
		ReadRatio: readRatio,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Println("json.Marshal failed:", err)
		return err
	}
	url := "http://" + workerIP
	fmt.Println(url)

	readReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	readReq.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println("http.Get err: ", err)
	}
	client := &http.Client{}

	readResp, err := client.Do(readReq)
	
	if err != nil {
		log.Println("http.Get err: ", err)
		return err
	}
	defer readResp.Body.Close()
	if readResp.StatusCode != http.StatusOK {
		log.Print("read失败，状态码：", readResp.StatusCode)
		return err
	}
	return nil
}
