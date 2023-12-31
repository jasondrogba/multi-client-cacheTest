package loadAlluxio

import (
	"fmt"
	"jasondrogba/multi-client-cacheTest/master/handleLock"
	"jasondrogba/multi-client-cacheTest/master/readyForEc2"
	"jasondrogba/multi-client-cacheTest/master/userMasterInfo"
	"log"
	"net/http"
	"runtime"
	"sync"
)

var wgLoad sync.WaitGroup

func TotalLoad(workerListInfo userMasterInfo.WorkerInfoList) {
	for _, w := range workerListInfo.WorkerInfoList {
		wgLoad.Add(1)
		fmt.Printf("WorkerID: %s, ReadRatio: %d, LoadFile: %s\n", w.WorkerId, w.ReadRatio, w.LoadFile)
		go multiLoad(w.WorkerId, w.LoadFile)
	}
	wgLoad.Wait()
	<-handleLock.GetLoadRunning()
}

func multiLoad(workerId string, loadFile string) error {
	instanceMap := readyForEc2.GetInstanceMap()
	if runtime.GOOS == "linux" {
		fmt.Println("Detected Linux system")
		err := sendLoadToWorker(instanceMap["Ec2Cluster-default-workers-"+workerId]+":8888/loadAlluxio", loadFile)
		if err != nil {
			return err
		}
	} else if runtime.GOOS == "darwin" {
		fmt.Println("Detected macOS system")
		err := sendLoadToWorker("localhost:8888/loadAlluxio", loadFile)
		if err != nil {
			return err
		}

	} else {
		fmt.Println("Unknown system")
	}
	return nil
}

func sendLoadToWorker(workerIP string, loadFile string) error {
	//fmt.Println("loadFile", loadFile)
	url := "http://" + workerIP + "?loadFile=" + loadFile
	fmt.Println(url)
	loadResp, err := http.Get(url)
	if err != nil {
		log.Println("http.Get err: ", err)
	}
	defer loadResp.Body.Close()
	if loadResp.StatusCode != http.StatusOK {
		log.Print("预热失败，状态码：", loadResp.StatusCode)
		wgLoad.Done()
		return err
	}
	wgLoad.Done()
	return nil

}
