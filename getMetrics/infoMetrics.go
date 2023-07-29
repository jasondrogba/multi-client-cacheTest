package main

type infoStruct struct {
	TestId int     `json:"testId"`
	Info   string  `json:"info"`
	Ufs    float64 `json:"ufs"`
	Remote float64 `json:"remote"`
}

type metricsInfo struct {
	Message string       `json:"message"`
	Result  []infoStruct `json:"result"`
}
