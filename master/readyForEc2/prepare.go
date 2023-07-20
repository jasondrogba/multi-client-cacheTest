package readyForEc2

import "jasondrogba/multi-client-cacheTest/master/ec2test"

var instanceMap map[string]string

func Prepare() {
	instanceMap = ec2test.Getec2Instance()
}
func GetInstanceMap() map[string]string {
	return instanceMap
}
