package main

import (
	"bufio"
	"fmt"
	"os"

	. "github.com/chtombre/helvar-go"
)

func main() {
	router := Router{
		Host: "10.40.105.10",
		Port: 50000,
	}

	defer router.Disconnect()

	router.Initialize()

	groups, err := router.GetGroups()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sumDevices := 0
	for _, g := range groups {
		sumDevices += len(g.Devices)
	}

	fmt.Printf("Found %d groups and %d devices", len(groups), sumDevices)
	fmt.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}