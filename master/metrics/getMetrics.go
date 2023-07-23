package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"jasondrogba/multi-client-cacheTest/master/ec2test"
	"net/http"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var totalReadUfs float64
var totalRemote float64
var mutex sync.Mutex

const ReadRemote = "BytesReadRemote"
const ReadUFS = "BytesReadPerUfs"

func BackProcess() (float64, float64) {
	//间隔5s执行一次getMetrics
	var count int
	var resultRemote, resultReadUfs float64
	instanceMap := ec2test.Getec2Instance()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	previousReadUfs, previousRemote := GetMetrics(instanceMap)
	totalRemote, totalReadUfs = 0, 0

	for range ticker.C {
		currentReadUfs, currentRemote := GetMetrics(instanceMap)
		if currentRemote == 0 && currentReadUfs == 0 {
			continue
		}
		if currentReadUfs == previousReadUfs && currentRemote == previousRemote {
			count++
			fmt.Println("count: ", count)
		} else {
			count = 0
		}

		if count == 7 {
			fmt.Println("task finish")
			count = 0
			resultRemote, resultReadUfs = currentRemote, currentReadUfs
			break
		}
		previousRemote, previousReadUfs = currentRemote, currentReadUfs
		totalRemote, totalReadUfs = 0, 0
	}
	fmt.Println(" go to next")
	return resultRemote, resultReadUfs

}

func GetMetrics(instanceMap map[string]string) (float64, float64) {
	//instanceMap := ec2test.Getec2Instance()
	//master := instanceMap["Ec2Cluster-default-masters-0"]
	//worker0 := instanceMap["Ec2Cluster-default-workers-0"]
	//worker1 := instanceMap["Ec2Cluster-default-workers-1"]
	//worker2 := instanceMap["Ec2Cluster-default-workers-2"]
	//worker3 := instanceMap["Ec2Cluster-default-workers-3"]

	// Get the metrics of the master node
	for key, value := range instanceMap {
		if strings.Contains(key, "workers") {
			wg.Add(1)
			go GetReadUfsFromWorker(value)
		}
	}
	wg.Wait()
	fmt.Println("total ReadFromUfs: ", totalReadUfs/1024/1024/1024, "GB")
	fmt.Println("total ReadFromRemote: ", totalRemote/1024/1024/1024, "GB")
	return totalReadUfs / 1024 / 1024 / 1024, totalRemote / 1024 / 1024 / 1024
	//GetReadUfsFromWorker(worker0)
	//GetReadUfsFromWorker(worker1)
	//GetReadUfsFromWorker(worker2)
	//GetReadUfsFromWorker(worker3)

}

func GetReadUfsFromWorker(hostname string) {
	// get prometheus metrics from master
	defer wg.Done()
	url := "http://" + hostname + ":30000/metrics/json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch data from Prometheus: %s\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read Prometheus response body: %s\n", err)
	}

	//fmt.Println(string(body))
	// 解析响应JSON
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		fmt.Printf("Failed to parse JSON response: %s\n", err)
		return
	}
	//fmt.Println(jsonData["gauges"][])
	// 如果jsonData["gauges"]不是一个数组，则尝试转换为map[string]interface{}
	countersMap, ok := jsonData["counters"].(map[string]interface{})
	if !ok {
		fmt.Println("Failed to parse 'gauges' field")
		return
	}
	// 现在可以通过字段名进行访问
	//value := gaugesMap["fieldName"]
	// ...
	var valueReadUfs interface{}
	//fmt.Println(countersMap["Worker.BytesReadPerUfs.UFS:s3:%2F%2Falluxio-tpch100%2FcacheTest.User:.Ec2Cluster-workers-0"])
	for key, value := range countersMap {
		if strings.Contains(key, ReadUFS) {
			//fmt.Println(key, value)
			valueReadUfs = value
			break
		}
	}

	if valueReadUfs == nil {
		fmt.Println("Failed to parse 'BytesReadPerUfs' field, valueReadUfs is nil")
		return
	}

	value, ok := valueReadUfs.(map[string]interface{})["count"]
	if !ok {
		fmt.Println("Failed to parse 'count' field")
		return
	}
	valueUfsFloat, ok := value.(float64)
	if !ok {
		fmt.Println("Failed to parse 'count' field")
		return
	}
	//fmt.Println(valueUfsFloat/1024/1024/1024, "GB")

	var valueRemote interface{}
	for key, value := range countersMap {
		if strings.Contains(key, ReadRemote) {
			//fmt.Println(key, value)
			valueRemote = value
			break
		}
	}
	if valueRemote == nil {
		fmt.Println("Failed to parse 'BytesReadRemote' field, valueRemote is nil")
		return
	}
	value, ok = valueRemote.(map[string]interface{})["count"]
	if !ok {
		fmt.Println("Failed to parse 'count' field")
		return
	}
	valueRemoteFloat, ok := value.(float64)
	if !ok {
		fmt.Println("Failed to parse 'count' field")
		return
	}
	//fmt.Println(valueRemoteFloat/1024/1024/1024, "GB")
	mutex.Lock()
	totalReadUfs += valueUfsFloat
	totalRemote += valueRemoteFloat
	mutex.Unlock()

}
