package main

import (
	"encoding/json"
	"log"
)

func SetJSON(policy string, worker []WorkerInfo) ([]byte, error) {
	jsonData := WorkerInfoList{
		Policy: policy,
		WorkerInfoList: []WorkerInfo{
			worker[0],
			worker[1],
			worker[2],
			worker[3],
		},
	}
	//fmt.Print(jsonData)
	jsonBytes, err := json.Marshal(jsonData)
	return jsonBytes, err
}

func SetWorkerJSON(workerId string, ratio string, loadFile string,
	hotFile string, totalFile string, count string) WorkerInfo {
	return WorkerInfo{
		WorkerId:  workerId,
		ReadRatio: ratio,
		LoadFile:  loadFile,
		HotFile:   hotFile,
		TotalFile: totalFile,
		Count:     count,
	}
}

func StartSetJSON(policy string, ratio string, hot string,
	total string, count string, warmUp bool) []byte {
	if warmUp {
		workerList := []WorkerInfo{}
		worker := SetWorkerJSON("0", ratio,
			"39", hot,
			total, count)
		workerList = append(workerList, worker)

		worker = SetWorkerJSON("1", ratio,
			"30", hot,
			total, count)
		workerList = append(workerList, worker)

		worker = SetWorkerJSON("2", ratio,
			"15", hot,
			total, count)
		workerList = append(workerList, worker)

		worker = SetWorkerJSON("3", ratio,
			"0", hot,
			total, count)
		workerList = append(workerList, worker)

		jsonBytes, err := SetJSON(policy, workerList)

		if err != nil {
			log.Println("json.Marshal failed:", err)
		}
		return jsonBytes
	} else {
		return RandSetJSON(policy, ratio, hot, total, count)
	}

}

func RandSetJSON(policy string, ratio string, hot string, total string, count string) []byte {
	workerList := []WorkerInfo{}
	worker := SetWorkerJSON("0", ratio,
		"0", hot,
		total, count)
	workerList = append(workerList, worker)

	worker = SetWorkerJSON("1", ratio,
		"0", hot,
		total, count)
	workerList = append(workerList, worker)

	worker = SetWorkerJSON("2", ratio,
		"0", hot,
		total, count)
	workerList = append(workerList, worker)

	worker = SetWorkerJSON("3", ratio,
		"0", hot,
		total, count)
	workerList = append(workerList, worker)

	jsonBytes, err := SetJSON(policy, workerList)

	if err != nil {
		log.Println("json.Marshal failed:", err)
	}
	return jsonBytes
}
