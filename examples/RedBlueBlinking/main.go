package main

import (
	"strconv"
	"time"

	. "github.com/chtombre/helvar-go"
)

func main() {
	router := Router{
		Host: "10.40.105.10",
		Port: 50000,
	}

	defer router.Disconnect()

	router.Initialize()

	DoLoop(router)
}

func DoLoop(router Router){
	loops := 0

	for loops < 10 {

		SendCommand(router, "105.10.4.9", 75,20)

		time.Sleep(200 * time.Millisecond)
		SendCommand(router, "105.10.4.11", 75,20)

		time.Sleep(1 * time.Millisecond)
		SendCommand(router, "105.10.4.9", 0, 20)

		time.Sleep(200 * time.Millisecond)
		SendCommand(router, "105.10.4.11", 0, 20)

		loops++
	}

	router.SendCommand(NewCommand(MT_COMMAND, DIRECT_LEVEL_GROUP, CommandParameter{GROUP, "702"}, CommandParameter{LEVEL, "0"}))
}

func SendCommand(router Router, address string, level int, fade int){
	router.SendCommand(NewCommand(
		MT_COMMAND,
		DIRECT_LEVEL_DEVICE,
		CommandParameter{LEVEL, strconv.Itoa(level)},
		CommandParameter{ADDRESS, address},
		CommandParameter{FADE_TIME, strconv.Itoa(fade)},
	))
}