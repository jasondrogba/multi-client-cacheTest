package userMasterInfo

type WorkerInfo struct {
	WorkerId  string `json:"worker_id"`
	ReadRatio string `json:"read_ratio"`
	LoadFile  string `json:"load_file"`
	HotFile   string `json:"hot_file"`
	TotalFile string `json:"total_file"`
	Count     string `json:"count"`
}

type WorkerInfoList struct {
	WorkerInfoList []WorkerInfo `json:"worker_info_list"`
}
