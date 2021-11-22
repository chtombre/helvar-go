package main

import (
	"math/rand"
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

	RandomLoop(router)
}

func RandomLoop(router Router){
	loops := 0

	router.SendCommand(NewCommand(MT_COMMAND, DIRECT_LEVEL_GROUP, CommandParameter{GROUP, "701"}, CommandParameter{LEVEL, "0"}))
	router.SendCommand(NewCommand(MT_COMMAND, DIRECT_LEVEL_GROUP, CommandParameter{GROUP, "702"}, CommandParameter{LEVEL, "0"}))
	for loops < 2 {

		var (
			red = rand.Intn(100)
			blue = rand.Intn(100)
			green = rand.Intn(100)
		)

		SendCommand(router, "105.10.4.5", int(red),400)
		SendCommand(router, "105.10.4.6", int(green),400)
		SendCommand(router, "105.10.4.7", int(blue),400)

		SendCommand(router, "105.10.4.9", int(red),400)
		SendCommand(router, "105.10.4.10", int(green),400)
		SendCommand(router, "105.10.4.11", int(blue),400)

		time.Sleep(5 * time.Second)

		loops++
	}

	SendCommand(router, "105.10.4.8", 80,400)
	SendCommand(router, "105.10.4.12", 80,400)

	//router.SendCommand(NewCommand(MT_COMMAND, DIRECT_LEVEL_GROUP, CommandParameter{GROUP, "701"}, CommandParameter{LEVEL, "75"}))
	//router.SendCommand(NewCommand(MT_COMMAND, DIRECT_LEVEL_GROUP, CommandParameter{GROUP, "702"}, CommandParameter{LEVEL, "75"}))
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
