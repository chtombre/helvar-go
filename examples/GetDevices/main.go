package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	. "github.com/chtombre/helvar-go"
)

func main() {

	routers := []Router {
		{Host: "10.40.168.50", Port: 50000},
		{Host: "10.40.168.51", Port: 50000},
		{Host: "10.40.168.52", Port: 50000},
		{Host: "10.40.168.53", Port: 50000},
		{Host: "10.40.168.54", Port: 50000},
		{Host: "10.40.168.55", Port: 50000},
	}

	sumDevices := 0

	var wg sync.WaitGroup

	defer func() {
		for idx := range routers {
			router := routers[idx]
			router.Disconnect()
		}
	}()

	for idx := range routers {
		router := routers[idx]
		wg.Add(1)
		router.Initialize()
		go func() {
			defer wg.Done()
			groups, err := router.GetGroups()
			if err != nil {

			}

			for idx := range groups {
				group := groups[idx]
				sumDevices += len(group.Devices)
			}
		}()
	}

	wg.Wait()

	fmt.Printf("Found %d devices", sumDevices)
	fmt.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func GetGroups (router *Router) []Group {
	groups, err := router.GetGroups()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return groups
}