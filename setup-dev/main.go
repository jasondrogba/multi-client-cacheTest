package main

import (
	"fmt"
	"jasondrogba/multi-client-cacheTest/setup-dev/ec2test"
	"os/exec"
	"sync"
)

var masterDir = "/Users/sunbury/go/src/jasondrogba/multi-client-cacheTest/master/master-server"
var workerDir = "/Users/sunbury/go/src/jasondrogba/multi-client-cacheTest/worker/worker-server"
var ec2Username = "ec2-user"
var remoteFilePath = "~"

var wg sync.WaitGroup

func main() {
	//获得所有的ec2实例
	instanceMap := ec2test.Getec2Instance()
	//将master生成的程序，和worker生成的程序分别发送到ec2实例上
	for k, v := range instanceMap {
		wg.Add(1)
		if k == "Ec2Cluster-default-masters-0" {
			go sendCMD(masterDir, v)
			continue
		}
		go sendCMD(workerDir, v)
	}
	//time.Sleep(time.Second * 1)
	wg.Wait()

	//启动所有的ec2实例上的master和worker
	//for k, v := range instanceMap {
	//	wg.Add(1)
	//	if k == "Ec2Cluster-default-masters-0" {
	//		go runCMD("master-server", v)
	//		continue
	//	}
	//	go runCMD("worker-server", v)
	//}
	//wg.Wait()
}

func sendCMD(Dir string, hostname string) {
	defer wg.Done()
	//cmd := exec.Command("bash", "-c", "scp", "-r", "-i", "~/.ssh/id_rsa", Dir, "ec2",
	fmt.Sprintf("%s@%s:%s", ec2Username, hostname, remoteFilePath)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("scp -r -i ~/.ssh/id_rsa %s  %s@%s:%s", Dir, ec2Username, hostname, remoteFilePath))

	err := cmd.Run()
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println("send success", hostname, Dir)
}

func runCMD(Dir string, hostname string) {
	defer wg.Done()
	//Dir := workerDir

	cmd := exec.Command("bash", "-c", fmt.Sprintf("ssh -i ~/.ssh/id_rsa %s@%s ./%s", ec2Username, hostname, Dir))
	err := cmd.Run()
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println("run success", hostname)

}
