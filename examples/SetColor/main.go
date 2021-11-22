package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"

	. "helvargo"
)

func main() {
	router := Router{
		Host: "10.40.105.10",
		Port: 50000,
	}

	defer router.Disconnect()

	router.Initialize()

	SetColor(router,"#EE9B00")
}

func SetColor(router Router, hexColor string){
	color, err := ParseHexColor(hexColor)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var (
		red = float32(color.R) / 255 * 100
		blue = float32(color.B) / 255 * 100
		green = float32(color.G) / 255 * 100
	)

	SendCommand(router, "105.10.4.9", int(red),20)
	SendCommand(router, "105.10.4.10", int(green),20)
	SendCommand(router, "105.10.4.11", int(blue),20)
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

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}
