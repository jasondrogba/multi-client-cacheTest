package userMasterInfo

type RequestData struct {
	MasterIP  string `json:"master_ip"`
	Count     string `json:"count"`
	HotFile   string `json:"hot_file"`
	TotalFile string `json:"total_file"`
	ReadRatio string `json:"read_ratio"`
}
