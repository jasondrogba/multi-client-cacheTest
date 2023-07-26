package loadAlluxio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jasondrogba/multi-client-cacheTest/master/handleLock"
	"jasondrogba/multi-client-cacheTest/master/metrics"
	"jasondrogba/multi-client-cacheTest/master/readyForEc2"
	"jasondrogba/multi-client-cacheTest/master/userMasterInfo"
	"log"
	"net/http"
	"runtime"
	"sync"
)

var readwg sync.WaitGroup
var infoReadUfs, infoRemote map[string]float64

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
	tmpReadUfs, tmpRemote := metrics.BackProcess()
	//totalReadUfs = append(totalReadUfs, tmpReadUfs)
	//totalRemote = append(totalRemote, tmpRemote)
	tmpInfo := "Read Policy:" + workerListInfo.Policy +
		"-Ratio:" + workerListInfo.WorkerInfoList[0].ReadRatio +
		"-LoadFile:" + workerListInfo.WorkerInfoList[0].LoadFile +
		"-HotFile:" + workerListInfo.WorkerInfoList[0].HotFile +
		"-TotalFile:" + workerListInfo.WorkerInfoList[0].TotalFile +
		"-Count:" + workerListInfo.WorkerInfoList[0].Count

	metrics.SetInfoUfs(tmpInfo, tmpReadUfs)
	metrics.SetInfoRemote(tmpInfo, tmpRemote)

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
	retryCount := 3
	for retryCount > 0 {
		readResp, err := client.Do(readReq)
		if err != nil {
			log.Println("http.Get err: ", err)
			retryCount--
			continue
		}
		if readResp.StatusCode != http.StatusOK {
			log.Print("read失败，状态码：", readResp.StatusCode)
			retryCount--
			continue
		}
		//err = readResp.Body.Close()
		//if err != nil {
		//	log.Println("readResp.Body.Close() err: ", err)
		//	return err
		//}
		break
	}
	//readResp, err := client.Do(readReq)
	//
	//if err != nil {
	//	log.Println("http.Get err: ", err)
	//	return err
	//}

	//if readResp.StatusCode != http.StatusOK {
	//	log.Print("read失败，状态码：", readResp.StatusCode)
	//	return err
	//}
	return nil
}
