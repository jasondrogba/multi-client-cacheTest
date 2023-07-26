package metrics

type infoStruct struct {
	TestId int     `json:"testId"`
	Info   string  `json:"info"`
	Ufs    float64 `json:"ufs"`
	Remote float64 `json:"remote"`
}

type metricsInfo struct {
	InfoList []infoStruct `json:"infoList"`
}
