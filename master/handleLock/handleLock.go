package handleLock

var running = make(chan struct{}, 1)     // 创建一个channel作为互斥锁
var loadRunning = make(chan struct{}, 1) // 创建一个channel作为互斥锁
var readRunning = make(chan struct{}, 1) // 创建一个channel作为互斥锁
func Getrunning() chan struct{} {
	return running
}

func GetLoadRunning() chan struct{} {
	return loadRunning
}

func GetReadRunning() chan struct{} {
	return readRunning
}
