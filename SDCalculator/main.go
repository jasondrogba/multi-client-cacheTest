package main

import (
	"fmt"
	"math"
	"math/rand"
)

func randomGenerate(totalFile int, hotFile int, readRatio int) int {
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
	//index = index % totalFile
	return index
}

func calculateMean(data map[int]int) float64 {
	total := 0
	for _, value := range data {
		total += value
	}
	return float64(total) / float64(len(data))
}

func calculateStandardDeviation(data map[int]int) float64 {
	mean := calculateMean(data)
	varianceSum := 0.0
	max := 0.0
	for _, value := range data {

		deviation := (float64(value) - mean)

		varianceSum += deviation * deviation
		if deviation*deviation > max {
			max = deviation * deviation
		}
	}

	variance := varianceSum / float64(len(data))
	standardDeviation := math.Sqrt(variance)
	return standardDeviation
}

func fairnessIndex(data map[int]int) float64 {
	//mean := calculateMean(data)
	sumPower := 0
	powerSum := 0
	for _, value := range data {
		sumPower += value * value
		powerSum += value
	}
	powerSum = powerSum * powerSum
	fairnessIndex := float64(powerSum) / float64(len(data)) / float64(sumPower)
	return fairnessIndex
}

func main() {
	// 示例的map[string]int
	//data := map[string]int{
	//	"file1.txt": 1,
	//	"file2.txt": 1,
	//	"file3.txt": 1,
	//	"file4.txt": 1,
	//	"file5.txt": 1,
	//}
	//
	//// 计算标准差
	//standardDeviation := calculateStandardDeviation(data)
	//fmt.Printf("标准差：%.2f\n", standardDeviation)
	//var counter int
	//data := make(map[int]int)
	//for i := 0; i < 100; i++ {
	//	generator := randomGenerate(100, 20, 20)
	//	data[generator]++
	//	if generator < 20 {
	//		counter++
	//	}
	//}
	//
	//standardDeviation := calculateStandardDeviation(data)
	//fairnessIndex := fairnessIndex(data)
	//fmt.Printf("标准差：%v\n", standardDeviation)
	//fmt.Printf("公平性指数：%v\n", fairnessIndex)
	//fmt.Println(counter)
	//fmt.Println(data)

	for i := 0; i < 3; i++ {
		fairnessIndexTest(400, 5000, 1000, 20)
		//fairnessIndexTest(400, 200, 40, 20)
		//fairnessIndexTest(400, 300, 60, 20)

	}
	for i := 0; i < 3; i++ {

		fairnessIndexTest(400, 5000, 1000, 80)

		//fairnessIndexTest(400, 200, 40, 80)
		//fairnessIndexTest(400, 300, 60, 80)

	}

}

func fairnessIndexTest(count int, total int, hotFile int, ratio int) {
	data := make(map[int]int)

	for i := 0; i < count; i++ {
		generator := randomGenerate(total, hotFile, ratio)
		data[generator]++
	}

	standardDeviation := calculateStandardDeviation(data)
	fairnessIndex := fairnessIndex(data)
	fmt.Printf("读取次数：%v,文件总数：%v，热门比例：%v，\n标准差：\n%v\n，公平性指数：\n%v\n",
		count, total, ratio, standardDeviation, fairnessIndex)

}
