package startTest

import (
	"jasondrogba/multi-client-cacheTest/master/handleLock"
	"jasondrogba/multi-client-cacheTest/master/loadAlluxio"
	"jasondrogba/multi-client-cacheTest/master/readyForEc2"
	"jasondrogba/multi-client-cacheTest/master/userMasterInfo"
)

var totalReadUfs, totalRemote []float64

func StartTraining(workerListInfo userMasterInfo.WorkerInfoList) {
	instanceMap := readyForEc2.GetInstanceMap()
	//启动alluxio
	handleLock.Getrunning() <- struct{}{}
	StartTest(instanceMap["Ec2Cluster-default-masters-0"], workerListInfo.Policy)
	//向所有的alluxio worker发送预热数据的命令
	handleLock.GetLoadRunning() <- struct{}{}
	loadAlluxio.TotalLoad(workerListInfo)
	handleLock.GetLoadRunning() <- struct{}{}

	//向所有的alluxio worker发送读取数据的命令
	handleLock.GetReadRunning() <- struct{}{}
	loadAlluxio.TotalRead(workerListInfo)
	handleLock.GetReadRunning() <- struct{}{}

	//tmpReadUfs, tmpRemote := metrics.BackProcess()
	//totalReadUfs = append(totalReadUfs, tmpReadUfs)
	//totalRemote = append(totalRemote, tmpRemote)
	//tmpInfo := "StartTraining Policy:" + workerListInfo.Policy +
	//	"-Ratio:" + workerListInfo.WorkerInfoList[0].ReadRatio +
	//	"-LoadFile:" + workerListInfo.WorkerInfoList[0].LoadFile +
	//	"-HotFile:" + workerListInfo.WorkerInfoList[0].HotFile +
	//	"-TotalFile:" + workerListInfo.WorkerInfoList[0].TotalFile +
	//	"-Count:" + workerListInfo.WorkerInfoList[0].Count
	//
	//metrics.SetInfo(tmpInfo, tmpReadUfs, tmpRemote)
	//释放
	<-handleLock.GetReadRunning()
	<-handleLock.GetLoadRunning()
	<-handleLock.GetTrainRunning()

}

func GetResult() ([]float64, []float64) {
	return totalReadUfs, totalRemote
}
