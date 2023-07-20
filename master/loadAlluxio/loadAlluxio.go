package loadAlluxio

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"runtime"
	"sync"
)

var wg sync.WaitGroup
var mutex sync.Mutex

func SshTest(instanceMap map[string]string) {
	config := SetupSSH()

	// Establish an SSH connection to the EC2 instance
	//master := instanceMap["Ec2Cluster-default-masters-0"]
	worker0 := instanceMap["Ec2Cluster-default-workers-0"]
	worker1 := instanceMap["Ec2Cluster-default-workers-1"]
	worker2 := instanceMap["Ec2Cluster-default-workers-2"]
	worker3 := instanceMap["Ec2Cluster-default-workers-3"]
	port := "22"

	fmt.Println("start worker0")
	//establishSSH(worker0, port, config, 20)
	for i := 1; i <= 39; i++ {
		wg.Add(1)
		go multiSSH(worker0, port, config, i)
	}
	fmt.Println("start worker1")
	//establishSSH(worker1, port, config, 10)
	for i := 1; i <= 39; i++ {
		wg.Add(1)

		go multiSSH(worker1, port, config, i)
	}
	fmt.Println("start worker2")
	//establishSSH(worker2, port, config, 5)
	for i := 1; i <= 39; i++ {
		wg.Add(1)

		go multiSSH(worker2, port, config, i)
	}
	fmt.Println("start worker3")
	//establishSSH(worker2, port, config, 5)
	for i := 1; i <= 39; i++ {
		wg.Add(1)

		go multiSSH(worker3, port, config, i)
	}
	wg.Wait()

}
func SetupSSH() *ssh.ClientConfig {
	// Read the private key file for the SSH connection
	PrivateKeyPath := ""
	if runtime.GOOS == "linux" {
		fmt.Println("Detected Linux system")
		PrivateKeyPath = "/home/ec2-user/.ssh/id_rsa"

	} else if runtime.GOOS == "darwin" {
		fmt.Println("Detected macOS system")
		PrivateKeyPath = "/Users/sunbury/.ssh/id_rsa"

	} else {
		fmt.Println("Unknown system")
	}
	privateKeyBytes, err := os.ReadFile(PrivateKeyPath)
	if err != nil {
		fmt.Println("Failed to read private key file:", err)
		os.Exit(1)
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		fmt.Println("Failed to parse private key:", err)
		os.Exit(1)
	}

	// Set up the SSH configuration
	config := &ssh.ClientConfig{
		User: "ec2-user",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config
}
func multiSSH(instance string, port string, config *ssh.ClientConfig, i int) {
	defer wg.Done()
	mutex.Lock()
	conn, err := ssh.Dial("tcp", instance+":"+port, config)
	if err != nil {
		fmt.Println("Failed to establish SSH connection:", err)
		os.Exit(1)
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Failed to create session:", err)
		os.Exit(1)
	}
	defer session.Close()
	mutex.Unlock()

	cmd := fmt.Sprintf("sudo su alluxio -c \"cd /opt/alluxio && ./bin/alluxio fs load /%d.txt --local flag\"", i)
	output, err := session.Output(cmd)
	if err != nil {
		fmt.Println("Failed to run command:", err)
		os.Exit(1)
	}
	fmt.Print(string(output))

}
