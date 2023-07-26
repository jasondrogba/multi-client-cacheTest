package handleLock

var running = make(chan struct{}, 1)       // 创建一个channel作为互斥锁
var loadRunning = make(chan struct{}, 1)   // 创建一个channel作为互斥锁
var readRunning = make(chan struct{}, 1)   // 创建一个channel作为互斥锁
var trainRunning = make(chan struct{}, 1)  // 创建一个channel作为互斥锁
var policyRunning = make(chan struct{}, 1) // 创建一个channel作为互斥锁
func Getrunning() chan struct{} {
	return running
}

func GetLoadRunning() chan struct{} {
	return loadRunning
}

func GetReadRunning() chan struct{} {
	return readRunning
}

func GetTrainRunning() chan struct{} {
	return trainRunning
}

func GetPolicyRunning() chan struct{} {
	return policyRunning
}
