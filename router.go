package helvargo

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var CommandResponseTimeout = 60 * time.Second
var KeepAliveDuration = 120 * time.Second

type Router struct {
	RouterId int
	Config *RouterConfig
	Host string
	Port int

	connection net.Conn

	connected bool

	commandReceived chan Command

	keepAlive chan bool
}

type RouterConfig struct {
	RouterId int
	//HasValue bool //Possible workaround nullable structs
}

func (r *Router) Initialize(){
	if !r.connected {
		r.Connect()
	}
}


func (r *Router) Id() int {
	if r.Config != nil {
		return r.Config.RouterId
	}
	return r.RouterId
}

func (r *Router) Connect() {
	conn, err := net.Dial("tcp",  r.Host + ":" + strconv.Itoa(r.Port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//setup channel for responses
	//r.commandReceived = make(chan Command)
	r.keepAlive = make(chan bool, 1)

	r.connection = conn
	r.connected = true

	//go r.Listener()

	go r.Keepalive()
}

func (r *Router) Disconnect() {
	r.connection.Close()
	r.connected = false

	close(r.keepAlive)
	//Log
}

func (r *Router) Reconnect() {
	r.Disconnect()
	r.Connect()
}

//func (r *Router) Listener(){
//	c := bufio.NewReader(r.connection)
//	for r.connected {
//		response, err := c.ReadString('#')
//		if err != nil {
//			//skip for now
//		}
//
//		command, err := ParseCommand(response)
//		if err != nil {
//			//skip for now
//		}
//		r.commandReceived <- *command
//	}
//}

func (r *Router) SendCommand(command Command){
	startTime := time.Now()
	_ = startTime

	strCommand := command.ToString()

	println("Sendt command: %s", strCommand)
	_, err := r.connection.Write([]byte(strCommand))
	if err != nil {
		println("Write to server failed:", err.Error())
		//	os.Exit(1)
	}
	time.Sleep(1 * time.Millisecond)

	//Tell keepalive that we just sent command to router, no need for pinging
	//If it's full, continue
	select {
		case r.keepAlive <- true:
		default:
	}


	//Check if command is not expecting answers
	for _, b := range COMMAND_TYPES_DONT_LISTEN_FOR_RESPONSE {
		if b == command.CommandType {
			return
		}
	}

	rc := make(chan string)

	go func() {
		c := bufio.NewReader(r.connection)
		response, err := c.ReadString('#')
		if err != nil {
			//skip for now
		}

		rc <- response
		close(rc)
	}()

	select {
		case response := <- rc:
			responseCommand, err := ParseCommand(response)
			if err != nil {
				//skip for now
			}
			println("Received response: ", responseCommand.ToString())

		case <- time.After(CommandResponseTimeout):
			//close(rc) //close the channel to indicate we're no longer interested in the response
			println("Response took too long")
			return
	}
}

func (r *Router) Keepalive() {
	for {
		select {
			case <- r.keepAlive:
				continue
			case <- time.After(KeepAliveDuration): {
				if !r.connected { return }

				//send command in a goroutine, so we will go back to listen to keepalive channel to keep from blocking
				r.SendCommand(NewCommand(MT_COMMAND, QUERY_ROUTER_TIME))
				r.SendCommand(NewCommand(MT_COMMAND, QUERY_ROUTER_TIME))
			}

		}
	}
}


