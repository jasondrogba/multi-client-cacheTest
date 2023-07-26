package main

import (
	"flag"
	"fmt"
)

func ParseArgs() (int, string, int) {
	count := flag.Int("count", 0, "The count parameter")
	policy := flag.String("policy", "", "The policy parameter")
	circle := flag.Int("circle", 0, "The circle parameter")

	flag.Parse()

	fmt.Println("Count:", *count)
	fmt.Println("Policy:", *policy)
	fmt.Println("Circle:", *circle)
	return *count, *policy, *circle
}
