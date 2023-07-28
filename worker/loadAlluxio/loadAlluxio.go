package loadAlluxio

import (
	"fmt"
	"jasondrogba/multi-client-cacheTest/worker/workerHandleLock"
	"os/exec"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func LoadAlluxio(fileCount int) {
	for i := 1; i <= fileCount; i++ {
		wg.Add(1)
		go multiLoad(i)
	}
	wg.Wait()
	//time.Sleep(10 * time.Second)
	<-workerHandleLock.GetLoadRunning()
}

func multiLoad(file int) {
	var cmd string
	if runtime.GOOS == "linux" {
		fmt.Println("Detected Linux system")
		cmd = fmt.Sprintf("sudo su alluxio -c \"cd /opt/alluxio && ./bin/alluxio fs load /%d.txt --local flag\"", file)

	} else if runtime.GOOS == "darwin" {
		fmt.Println("Detected macOS system")
		cmd = fmt.Sprintf("echo this is a test %d", file)

	} else {
		fmt.Println("Unknown system")
	}
	runcmd := exec.Command("bash", "-c", cmd)
	output, err := runcmd.Output()
	if err != nil {
		fmt.Println("Failed to run command:", err)
	}
	fmt.Print(string(output))
	wg.Done()

}
