package loadAlluxio

import (
	"fmt"
	alluxio "github.com/Alluxio/alluxio-go"
	"github.com/Alluxio/alluxio-go/option"
	"io"
	"jasondrogba/multi-client-cacheTest/worker/workerHandleLock"
	"log"
	"math/rand"
	"os"
)

func ReadAlluxio(masterIP string, count int, readRatio int, hotFile int, totalFile int) {
	//TODO():将循环改为并发
	//hostname := instanceMap["Ec2Cluster-default-masters-0"]
	rand.Seed(int64(12345))
	for i := 1; i <= count; i++ {

		multiReadRand(masterIP, readRatio, hotFile, totalFile)
		//if i == (count / 2) {
		//	resultRemote, resultUfs := metricsTest.BackProcess(instanceMap)
		//	fmt.Println("前一半的远程读取时间halfresultRemote：", resultRemote)
		//	fmt.Println("前一半的UFS读取时间halfresultUfs：", resultUfs)
		//	if dynamic {
		//		fmt.Println("现在动态缓存策略，需要调整REPLICA到LRU")
		//		startTest.SwitchLRU()
		//		dynamic = false
		//	}
		//
		//}
	}
	<-workerHandleLock.GetLoadRunning()

}

func multiReadRand(hostname string, readRatio int, hotFile int, totalFile int) {
	hostName, err := os.Hostname()
	if err != nil {
		fmt.Println("Failed to get hostname:", err)
		return
	}
	//fmt.Println("hostname:", hostName)
	fs := alluxio.NewClient(hostName, 39999, 0)
	index := rand.Int()
	//if i <= count/2 {
	//	index = index % 2600
	//	if index > 300 {
	//		index = (index-300)/60 + 1
	//	}
	//} else {
	//	index = index%300 + 1
	//}
	rangeNumber := (totalFile - hotFile) * 100 / (100 - readRatio)
	index = index % rangeNumber
	if index > totalFile {
		index = (index-totalFile)*hotFile/rangeNumber + 1
	}

	pathfile := fmt.Sprintf("/%d.txt", index)
	v := 0
	//log.Println("种子：", index)
	//v := 0
	//mutex.Lock()
	exists, err := fs.Exists(pathfile, &option.Exists{})
	if err != nil {
		log.Println(err)
	}
	//log.Println(index, "文件是否存在：", exists)
	if exists {
		f, err := fs.OpenFile(pathfile, &option.OpenFile{})
		defer fs.Close(f)
		if err != nil {
			log.Println(err)
		}
		//log.Println(index, "文件打开成功")
		data, err := fs.Read(f)
		if err != nil {
			log.Println(err)
		}
		//log.Println(index, "文件读取成功")
		//mutex.Unlock()
		defer data.Close()
		content, err := io.ReadAll(data)
		if err != nil {
			log.Println(err)
		}
		//log.Println(index, "文件本地IO成功")
		v = len(content)
		log.Print(index, "文件内容长度:", v)
		//TODO():加上运行的时间，可以对比时间消耗是否有减少
	} else {
		log.Println(index, "文件不存在")
		//mutex.Unlock()
	}
}
