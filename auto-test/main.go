package main

import (
	"bytes"
	"encoding/json"
	"jasondrogba/multi-client-cacheTest/auto-test/ec2test"
	"log"
	"net/http"
)

type ServerResponse struct {
	ReadUfs    []float64 `json:"readUfs"`
	ReadRemote []float64 `json:"remote"`
}

func main() {
	instanceMap := ec2test.Getec2Instance()
	MasterIp := instanceMap["Ec2Cluster-default-masters-0"]
	startUrl := "http://" + MasterIp + ":8080/startTraining"
	readUrl := "http://" + MasterIp + ":8080/readAlluxio"
	resultUrl := "http://" + MasterIp + ":8080/getResult"
	log.Println("startUrl:", startUrl)
	log.Println("readUrl:", readUrl)
	policy := "LRU"
	jsonData := WorkerInfoList{
		Policy: policy,
		WorkerInfoList: []WorkerInfo{
			{
				WorkerId:  "0",
				ReadRatio: "0.5",
				LoadFile:  "0",
				HotFile:   "0",
				TotalFile: "0",
				Count:     "0",
			},
			{
				WorkerId:  "1",
				ReadRatio: "0.5",
				LoadFile:  "0",
				HotFile:   "0",
				TotalFile: "0",
				Count:     "0",
			},
			{
				WorkerId:  "2",
				ReadRatio: "0.5",
				LoadFile:  "0",
				HotFile:   "0",
				TotalFile: "0",
				Count:     "0",
			},
			{
				WorkerId:  "3",
				ReadRatio: "0.5",
				LoadFile:  "0",
				HotFile:   "0",
				TotalFile: "0",
				Count:     "0",
			},
		},
	}

	jsonBytes, err := json.Marshal(jsonData)

	if err != nil {
		log.Println("json.Marshal failed:", err)
	}
	client := &http.Client{}
	getReq, err := http.NewRequest("GET", resultUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("http.NewRequest err: ", err)
	}
	getResp, err := client.Do(getReq)
	if err != nil {
		log.Println("http.Get err: ", err)
	}
	// Check the response status code
	if getResp.StatusCode == http.StatusOK {
		// The server has started the training, exit the loop
		log.Println("Server started training.")
		var responseData ServerResponse
		err := json.NewDecoder(getResp.Body).Decode(&responseData)
		if err != nil {
			log.Println("json.NewDecoder failed:", err)
		}
		//log.Println("readUfs:", responseData)
		//log.Println("readUfs:", responseData.ReadUfs)
		//log.Println("readRemote:", responseData.ReadRemote)

	} else {
		log.Println("Server is not ready. Status code:", getResp.StatusCode)
	}

}
