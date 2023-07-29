package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		getResult()
	}
	// Check the response status code

}

func getResult() {
	getMapResultUrl := "http://ec2-34-238-83-148.compute-1.amazonaws.com:8080/getMapResult"
	client := &http.Client{}
	getReq, err := http.NewRequest("GET", getMapResultUrl, nil)
	if err != nil {
		log.Println("http.NewRequest err: ", err)
	}
	getResp, err := client.Do(getReq)
	if err != nil {
		log.Println("http.Get err: ", err)
	}
	defer getResp.Body.Close()
	var resultList metricsInfo
	err = json.NewDecoder(getResp.Body).Decode(&resultList)
	if err != nil {
		log.Println("json.NewDecoder err: ", err)
	}
	log.Println("resultList: ", resultList)

	//创建文件
	path := "./result.json"
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("os.OpenFile err: ", err)
	}
	defer file.Close()

	if len(resultList.Result) == 0 {
		log.Println("resultList.Result is empty")
		return
	} else {
		//创建文件
		// 写入数据到文件
		dataBytes, err := json.Marshal(resultList.Result)
		if err != nil {
			fmt.Println("Error marshalling data:", err)
			return
		}
		_, err = file.Write(dataBytes)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	fmt.Println("Data written to result.json")

	//得到getResp的结果，解析JSON
}
