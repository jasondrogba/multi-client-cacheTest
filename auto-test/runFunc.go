package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func RunStartFunc(jsonBytes []byte, startUrl string) {
	client := &http.Client{}
	getReq, err := http.NewRequest("GET", startUrl, bytes.NewBuffer(jsonBytes))
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
		log.Println("Server started training.Status code:", getResp.StatusCode)

	} else {
		log.Println("Server is not ready. Status code:", getResp.StatusCode)
	}
}

func RunReadFunc(jsonBytes []byte, readUrl string) {
	client := &http.Client{}
	getReq, err := http.NewRequest("POST", readUrl, bytes.NewBuffer(jsonBytes))
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
		log.Println("Server started reading.")

	} else {
		log.Println("Server is not ready. Status code:", getResp.StatusCode)
	}
}

func RunSetPolicyFunc(jsonBytes []byte, setPolicyUrl string) {
	client := &http.Client{}
	getReq, err := http.NewRequest("GET", setPolicyUrl, bytes.NewBuffer(jsonBytes))
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
		log.Println("Server started reading.")

	} else {
		log.Println("Server is not ready. Status code:", getResp.StatusCode)
	}
}

func RunCheckStatusFunc(checkStatusUrl string) CheckResponse {
	client := &http.Client{}
	getReq, err := http.NewRequest("GET", checkStatusUrl, nil)
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
		log.Println("Server started reading.")

		//得到getResp的结果，解析JSON
		var checkResp CheckResponse
		err = json.NewDecoder(getResp.Body).Decode(&checkResp)
		if err != nil {
			log.Println("json.NewDecoder failed:", err)
		}
		//log.Println("checkResp:", checkResp)
		return checkResp
	} else {
		log.Println("Server is not ready. Status code:", getResp.StatusCode)
	}
	return CheckResponse{}
}
