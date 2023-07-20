package workerHandleLock

var loadRunning = make(chan struct{}, 1) // 创建一个channel作为互斥锁

func GetLoadRunning() chan struct{} {
	return loadRunning
}
