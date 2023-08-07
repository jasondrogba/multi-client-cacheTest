package main

import (
	"fmt"
	"github.com/Alluxio/alluxio-go"
	"github.com/Alluxio/alluxio-go/option"
	"io"
	"log"
	"math/rand"
	"os"
)

func main() {
	localReadTest(400, 20, 20, 100)

}

func localReadTest(count int, readRatio int, hotFile int, totalFile int) {
	for i := 1; i <= count; i++ {
		localRead(readRatio, hotFile, totalFile)
	}
}

func localRead(readRatio int, hotFile int, totalFile int) {
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
