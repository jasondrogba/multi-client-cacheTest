package ec2test

//get ec2 instance ip
import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"os"
	"sync"
)

var wg sync.WaitGroup
var mutex sync.Mutex

func Getec2Instance() map[string]string {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		os.Exit(1)
	}

	client := ec2.NewFromConfig(cfg)
	resp, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	})
	if err != nil {
		fmt.Println("Error describing instances:", err)
		return nil
	}
	// Print the instance ID, instance name, and IP address for each instance.
	instanceMap := make(map[string]string)
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			wg.Add(1)
			go func() {
				defer wg.Done()
				instanceID := aws.ToString(instance.InstanceId)
				instanceName := getInstanceName(instance)
				ipAddress := getIPAddress(instance)
				mutex.Lock()
				instanceMap[instanceName] = ipAddress
				mutex.Unlock()
				fmt.Printf("Instance ID: %s, Instance Name: %s, IP Address: %s\n", instanceID, instanceName, ipAddress)
			}()

		}
	}
	wg.Wait()
	return instanceMap
}

// getInstanceName returns the value of the "Name" tag for the specified instance,
// or an empty string if the tag is not present.
func getInstanceName(instance types.Instance) string {
	for _, tag := range instance.Tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}
	return ""
}

// getIPAddress returns the IP address for the specified instance,
// or an empty string if the instance does not have a public IP address.
func getIPAddress(instance types.Instance) string {
	if instance.PublicIpAddress != nil {
		return aws.ToString(instance.PublicDnsName)
	}
	return ""
}
