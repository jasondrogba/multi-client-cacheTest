package startTest

import (
	"jasondrogba/multi-client-cacheTest/master/handleLock"
	"jasondrogba/multi-client-cacheTest/master/loadAlluxio"
	"jasondrogba/multi-client-cacheTest/master/readyForEc2"
	"jasondrogba/multi-client-cacheTest/master/userMasterInfo"
)

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

	//释放
	<-handleLock.GetReadRunning()
	<-handleLock.GetLoadRunning()

}
